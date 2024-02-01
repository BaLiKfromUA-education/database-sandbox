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

func (dao *CounterDao) GetAtomicLong(ctx context.Context, name string) *hazelcast.AtomicLong {
	cp := dao.client.CPSubsystem()
	counter, err := cp.GetAtomicLong(ctx, name)
	if err != nil {
		log.Fatalf("Failed to get atomic long: %v", err)
	}
	return counter
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

func lockAndIncrement(ctx context.Context, distMap *hazelcast.Map, key string) {
	intValue := int64(0)
	// Create a new unique lock context.
	// https://pkg.go.dev/github.com/hazelcast/hazelcast-go-client#hdr-Using_Locks
	lockCtx := distMap.NewLockContext(ctx)
	// Lock the key.
	// The key cannot be unlocked without the same lock context.
	if err := distMap.Lock(lockCtx, key); err != nil {
		log.Fatalf("Error on Lock: %v", err)
	}
	// Remember to unlock the key, otherwise it won't be accessible elsewhere.
	defer distMap.Unlock(lockCtx, key)
	// The same lock context, or a derived one from that lock context must be used,
	// otherwise the Get operation below will block.
	v, err := distMap.Get(lockCtx, key)
	if err != nil {
		log.Fatalf("Error on Get: %v", err)
	}
	// If v is not nil, then there's already a value for the key.
	if v != nil {
		intValue = v.(int64)
	}
	// Increment and set the value back.
	intValue++
	// The same lock context, or a derived one from that lock context must be used,
	// otherwise the Set operation below will block.
	if err = distMap.Set(lockCtx, key, intValue); err != nil {
		log.Fatalf("Error on Set: %v", err)
	}
}

func (dao *CounterDao) counterWithPessimisticBlockingImpl(ctx context.Context, key string, distMap *hazelcast.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10_000; i++ {
		if i%100 == 0 {
			log.Printf("[Pessimistic] At %d\n", i)
		}
		// taken from https://github.com/hazelcast/hazelcast-go-client/blob/e7a962174982d98a3e3840ab3bec917bf67596a0/examples/map/lock/main.go
		lockAndIncrement(ctx, distMap, key)
	}
}

func (dao *CounterDao) counterWithOptimisticBlockingImpl(ctx context.Context, key string, distMap *hazelcast.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	debugMode := dao.debugMode.Load()

	for i := 0; i < 10_000; i++ {
		/*if i%1000 == 0 {
			log.Printf("[Optimistic] At %d\n", i)
		}*/

		for {
			counter, err := distMap.Get(ctx, key)
			if err != nil {
				log.Fatalf("Error on Get: %v", err)
			}

			cnt := counter.(int64)

			ok, err := distMap.ReplaceIfSame(ctx, key, cnt, cnt+1)
			if err != nil {
				log.Fatalf("Error on Set: %v", err)
			}

			if ok {
				break
			} else if debugMode {
				dao.misses.Add(1)
			}
		}

	}
}

func (dao *CounterDao) ExecuteCounterWithoutBlocking(ctx context.Context, name string, key string) {
	dao.execute(ctx, name, key, dao.counterWithoutBlockingImpl)
}

func (dao *CounterDao) ExecuteCounterWithPessimisticBlocking(ctx context.Context, name string, key string) {
	dao.execute(ctx, name, key, dao.counterWithPessimisticBlockingImpl)
}

func (dao *CounterDao) ExecuteCounterWithOptimisticBlocking(ctx context.Context, name string, key string, debugMode bool) {
	dao.debugMode.Store(debugMode)
	dao.misses.Store(0)

	dao.execute(ctx, name, key, dao.counterWithOptimisticBlockingImpl)

	if debugMode {
		log.Printf("Number of misses: %d\n", dao.misses.Load())
	}
	dao.debugMode.Store(false)
}

func (dao *CounterDao) counterWithAtomicLongImpl(ctx context.Context, counter *hazelcast.AtomicLong, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := 0; j < 10_000; j++ {
		_, err := counter.IncrementAndGet(ctx)
		if err != nil {
			log.Fatalf("Failed to increment counter: %v", err)
		}
	}
}

func (dao *CounterDao) ExecuteCounterWithAtomicLong(ctx context.Context, name string) {
	counter := dao.GetAtomicLong(ctx, name)
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go dao.counterWithAtomicLongImpl(ctx, counter, &wg)
	}

	wg.Wait()
}
