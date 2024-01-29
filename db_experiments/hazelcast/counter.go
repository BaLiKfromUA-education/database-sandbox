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
	config := hazelcast.Config{}
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

func (dao *CounterDao) execute(ctx context.Context, name string, key string, task func(context.Context, string, string, *sync.WaitGroup)) {
	var wg sync.WaitGroup
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go task(ctx, name, key, &wg)
	}

	wg.Wait()
}

func (dao *CounterDao) counterWithoutBlockingImpl(ctx context.Context, name string, key string, wg *sync.WaitGroup) {
	defer wg.Done()

	distMap := dao.GetMap(ctx, name)

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

func (dao *CounterDao) ExecuteCounterWithoutBlocking(ctx context.Context, name string, key string) {
	dao.execute(ctx, name, key, dao.counterWithoutBlockingImpl)
}
