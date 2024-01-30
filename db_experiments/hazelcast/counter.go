package hazelcast

import (
	"context"
	"github.com/hazelcast/hazelcast-go-client"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type CounterDao struct {
	client    *hazelcast.Client
	debugMode atomic.Bool
	misses    atomic.Int64
}

func buildConfig() hazelcast.Config {
	config := hazelcast.NewConfig()
	config.Cluster.Name = "distributed_databases"
	addresses, ok := os.LookupEnv("HAZELCAST_ADDRESSES")

	if !ok {
		log.Fatal("Failed to get addresses of hazecast nodes")
	}

	config.Cluster.Network.Addresses = strings.Split(addresses, ",")

	return config
}

func CreateDao(ctx context.Context) *CounterDao {
	client, err := hazelcast.StartNewClientWithConfig(ctx, buildConfig())

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &CounterDao{
		client: client,
	}
}

func (dao *CounterDao) GetMap(ctx context.Context, name string) *hazelcast.Map {
	distMap, err := dao.client.GetMap(ctx, name)

	if err != nil {
		log.Fatalf("Failed to get dist. map: %v", err)
	}

	return distMap
}

func (dao *CounterDao) execute(ctx context.Context, name string, key string, task func(context.Context, string, *hazelcast.Map, *sync.WaitGroup)) {
	var wg sync.WaitGroup
	n := 10
	wg.Add(n)

	distMap := dao.GetMap(ctx, name)

	for i := 0; i < n; i++ {
		go task(ctx, key, distMap, &wg)
	}

	wg.Wait()
}

func (dao *CounterDao) counterWithoutBlockingImpl(ctx context.Context, key string, distMap *hazelcast.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10_000; i++ {
		counter, err := distMap.Get(ctx, key)
		if err != nil {
			log.Fatalf("Error on Get: %v", err)
		}

		cnt := counter.(int64)
		cnt += 1

		err = distMap.Set(ctx, key, cnt)
		if err != nil {
			log.Fatalf("Error on Set: %v", err)
		}
	}
}

func (dao *CounterDao) counterWithPessimisticBlockingImpl(ctx context.Context, key string, distMap *hazelcast.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10_000; i++ {
		if i%100 == 0 {
			log.Printf("At %d\n", i)
		}

		// https://pkg.go.dev/github.com/hazelcast/hazelcast-go-client#hdr-Using_Locks
		lockCtx := distMap.NewLockContext(nil)

		err := distMap.Lock(lockCtx, key)
		if err != nil {
			log.Fatalf("Error on Lock: %v", err)
		}

		counter, err := distMap.Get(lockCtx, key)
		if err != nil {
			_ = distMap.Unlock(lockCtx, key)
			log.Fatalf("Error on Get: %v", err)
		}

		cnt := counter.(int64)
		cnt += 1

		err = distMap.Set(lockCtx, key, cnt)
		if err != nil {
			_ = distMap.Unlock(lockCtx, key)
			log.Fatalf("Error on Set: %v", err)
		}

		err = distMap.Unlock(lockCtx, key)
		if err != nil {
			log.Fatalf("Error of Unlock: %v", err)
		}
	}
}

func (dao *CounterDao) ExecuteCounterWithoutBlocking(ctx context.Context, name string, key string) {
	dao.execute(ctx, name, key, dao.counterWithoutBlockingImpl)
}

func (dao *CounterDao) ExecuteCounterWithPessimisticBlocking(ctx context.Context, name string, key string) {
	dao.execute(ctx, name, key, dao.counterWithPessimisticBlockingImpl)
}
