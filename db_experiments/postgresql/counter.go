package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type CounterDao struct {
	pool      *sql.DB
	debugMode atomic.Bool
	misses    atomic.Int64
}

func buildDsn() string {
	username, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		log.Fatalf("USER is not set")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		log.Fatalf("PASSWORD is not set")
	}

	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		log.Fatalf("HOST is not set")
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		log.Fatalf("PORT is not set")
	}

	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		log.Fatalf("DB name is not set")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}

func CreateDao(ctx context.Context) *CounterDao {
	// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
	dbPool, err := sql.Open("pgx", buildDsn())
	dbPool.SetMaxIdleConns(10)

	if err != nil {
		log.Fatalf("error during pool creation: %v", err)
	}

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS user_counter(
                user_id integer,
        		counter integer,
        		version integer
        )
  	`

	_, e := dbPool.ExecContext(ctx, createTableQuery)

	if e != nil {
		log.Fatalf("error during creation of table: %v", e)
	}

	return &CounterDao{pool: dbPool}
}

func (db *CounterDao) InsertBaseRecord(ctx context.Context, id int) {
	statement := "INSERT INTO user_counter (user_id, counter, version) VALUES ($1, 0, 0);"
	_, err := db.pool.ExecContext(ctx, statement, id)
	if err != nil {
		log.Fatalf("error during inserting base record: %v", err)
	}
}

func (db *CounterDao) CleanUp(ctx context.Context, id int) {
	statement := "DELETE FROM user_counter WHERE user_id=$1;"
	_, err := db.pool.ExecContext(ctx, statement, id)
	if err != nil {
		log.Fatalf("error during clean up: %v", err)
	}
}

func (db *CounterDao) GetResult(ctx context.Context, id int) int {
	statement := "SELECT counter FROM user_counter WHERE user_id=$1;"
	var counter int
	err := db.pool.QueryRowContext(ctx, statement, id).Scan(&counter)
	if err != nil {
		log.Fatalf("error during select: %v", err)
	}

	return counter
}

func (db *CounterDao) execute(ctx context.Context, id int, task func(context.Context, int, *sync.WaitGroup)) {
	var wg sync.WaitGroup
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go task(ctx, id, &wg)
	}

	wg.Wait()
}

func (db *CounterDao) lostUpdateImpl(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	selectStatement := "SELECT counter FROM user_counter WHERE user_id = $1;"
	updateStatement := "UPDATE user_counter SET counter = $1 WHERE user_id = $2;"

	var counter int
	for i := 0; i < 10_000; i++ {
		if err := db.pool.QueryRowContext(ctx, selectStatement, id).Scan(&counter); err != nil {
			log.Fatalf("error during select: %v", err)
		}

		counter += 1

		tx, err := db.pool.BeginTx(ctx, nil)
		if err != nil {
			log.Fatalf("error during creation of transaction: %v", err)
		}

		if _, err := tx.ExecContext(ctx, updateStatement, counter, id); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during update: %v", err)
		}

		if err = tx.Commit(); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during commit")
		}
	}

}

func (db *CounterDao) inplaceUpdateImpl(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	updateStatement := "UPDATE user_counter SET counter = counter + 1 WHERE user_id = $1;"
	for i := 0; i < 10_000; i++ {
		tx, err := db.pool.BeginTx(ctx, nil)
		if err != nil {
			log.Fatalf("error during creation of transaction: %v", err)
		}

		if _, err := tx.ExecContext(ctx, updateStatement, id); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during update: %v", err)
		}

		if err = tx.Commit(); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during commit")
		}
	}
}

func (db *CounterDao) rowLevelLockingImpl(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	selectStatement := "SELECT counter FROM user_counter WHERE user_id = $1 FOR UPDATE;"
	updateStatement := "UPDATE user_counter SET counter = $1 WHERE user_id = $2;"

	var counter int
	for i := 0; i < 10_000; i++ {
		tx, err := db.pool.BeginTx(ctx, nil)
		if err != nil {
			log.Fatalf("error during creation of transaction: %v", err)
		}

		if err := tx.QueryRowContext(ctx, selectStatement, id).Scan(&counter); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during select: %v", err)
		}

		counter += 1

		if _, err := tx.ExecContext(ctx, updateStatement, counter, id); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during update: %v", err)
		}

		if err = tx.Commit(); err != nil {
			_ = tx.Rollback()
			log.Fatalf("error during commit")
		}
	}
}

func (db *CounterDao) optimisticConcurrencyControlImpl(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	isDebug := db.debugMode.Load()

	selectStatement := "SELECT counter, version FROM user_counter WHERE user_id = $1;"
	updateStatement := "UPDATE user_counter SET counter = $1, version = $2 + 1 WHERE user_id = $3 AND version = $2;"

	var counter, version int
	for i := 0; i < 10_000; i++ {
		for {
			if err := db.pool.QueryRowContext(ctx, selectStatement, id).Scan(&counter, &version); err != nil {
				log.Fatalf("error during select: %v", err)
			}

			counter += 1

			tx, err := db.pool.BeginTx(ctx, nil)
			if err != nil {
				log.Fatalf("error during creation of transaction: %v", err)
			}

			res, err := tx.ExecContext(ctx, updateStatement, counter, version, id)

			if err != nil {
				_ = tx.Rollback()
				log.Fatalf("error during update: %v", err)
			}

			if err = tx.Commit(); err != nil {
				_ = tx.Rollback()
				log.Fatalf("error during commit")
			}

			rowCount, err := res.RowsAffected()
			if err != nil {
				log.Fatalf("error during row affected: %v", err)
			}

			if rowCount > 0 {
				break
			} else if isDebug {
				db.misses.Add(1)
			}
		}
	}
}

func (db *CounterDao) ExecuteLostUpdate(ctx context.Context, id int) {
	db.execute(ctx, id, db.lostUpdateImpl)
}

func (db *CounterDao) ExecuteInPlaceUpdate(ctx context.Context, id int) {
	db.execute(ctx, id, db.inplaceUpdateImpl)
}

func (db *CounterDao) ExecuteRowLevelLockingUpdate(ctx context.Context, id int) {
	db.execute(ctx, id, db.rowLevelLockingImpl)
}

func (db *CounterDao) ExecuteOptimisticConcurrencyControl(ctx context.Context, id int, isDebug bool) {
	if isDebug {
		db.debugMode.Store(isDebug)
		db.misses.Store(0)

		defer func() {
			db.debugMode.Store(false)
			db.misses.Store(0)
		}()
	}

	db.execute(ctx, id, db.optimisticConcurrencyControlImpl)

	if isDebug {
		log.Printf("Number of misses: %d\n", db.misses.Load())
	}
}
