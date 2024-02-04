package main

import (
	"context"
	"github.com/hazelcast/hazelcast-go-client"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

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

func buildConfig() hazelcast.Config {
	config := hazelcast.NewConfig()
	config.Cluster.Name = "test_lock"
	addresses, ok := os.LookupEnv("HAZELCAST_ADDRESSES")

	if !ok {
		log.Fatal("Failed to get addresses of hazecast nodes")
	}

	config.Cluster.Network.Addresses = strings.Split(addresses, ",")

	return config
}

func main() {
	time.Sleep(10 * time.Second) // just wait for node to be ready

	ctx := context.Background()
	client, err := hazelcast.StartNewClientWithConfig(ctx, buildConfig())

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	distMap, err := client.GetMap(ctx, "my-test")

	if err != nil {
		log.Fatalf("Failed to get dist. map: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(10)

	key := "test-key"

	start := time.Now()
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 10_000; j++ {
				if j%1000 == 0 {
					log.Printf("[Pessimistic] At %d\n", j)
				}
				// taken from https://github.com/hazelcast/hazelcast-go-client/blob/e7a962174982d98a3e3840ab3bec917bf67596a0/examples/map/lock/main.go
				lockAndIncrement(ctx, distMap, key)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	v, err := distMap.Get(ctx, key)
	if err != nil {
		log.Fatalf("Error on Get: %v", err)
	}

	log.Printf("Pessimistic locking update took %s, value=%d", elapsed, v.(int64))
}
