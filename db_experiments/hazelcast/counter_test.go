package hazelcast

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestCounterWithoutBlocking(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	counter := CreateDao(ctx)

	mapName := "counter_without_blocking"
	keyName := mapName + "_key"

	distMap := counter.GetMap(ctx, mapName)
	err := distMap.Set(ctx, keyName, 0)
	require.NoError(t, err)

	// WHEN
	counter.ExecuteCounterWithoutBlocking(ctx, mapName, keyName)

	// THEN
	value, err := distMap.Get(ctx, keyName)

	log.Printf("Final value: %v", value)

	require.NoError(t, err)
	require.True(t, value.(int64) > int64(0))
}

func TestCounterWithPessimisticBlocking(t *testing.T) {
	t.Skip("taking too long for now... TODO: fix")

	// GIVEN
	ctx := context.TODO()
	counter := CreateDao(ctx)

	mapName := "counter_with_pessimistic_blocking"
	keyName := mapName + "_key"

	distMap := counter.GetMap(ctx, mapName)
	err := distMap.Set(ctx, keyName, 0)
	require.NoError(t, err)

	// WHEN
	counter.ExecuteCounterWithPessimisticBlocking(ctx, mapName, keyName)

	// THEN
	value, err := distMap.Get(ctx, keyName)
	require.NoError(t, err)

	require.Equal(t, int64(100_000), value.(int64))
}

func TestCounterWithOptimisticBlocking(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	counter := CreateDao(ctx)

	mapName := "counter_with_optimistic_blocking"
	keyName := mapName + "_key"

	distMap := counter.GetMap(ctx, mapName)
	err := distMap.Set(ctx, keyName, 0)
	require.NoError(t, err)

	// WHEN
	counter.ExecuteCounterWithOptimisticBlocking(ctx, mapName, keyName, true)

	// THEN
	value, err := distMap.Get(ctx, keyName)
	require.NoError(t, err)

	require.Equal(t, int64(100_000), value.(int64))
}

func TestCounterWithAtomicLong(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	counter := CreateDao(ctx)

	name := "counter_with_atomic_long"

	atomic := counter.GetAtomicLong(ctx, name)
	err := atomic.Set(ctx, 0)
	require.NoError(t, err)

	// WHEN
	counter.ExecuteCounterWithAtomicLong(ctx, name)

	// THEN
	value, err := atomic.Get(ctx)
	require.NoError(t, err)

	require.Equal(t, int64(100_000), value)
}

func BenchmarkCounterWithoutBlocking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// GIVEN
		ctx := context.TODO()
		counter := CreateDao(ctx)

		rand.Seed(time.Now().Unix())
		mapName := fmt.Sprintf("without_blocking_%d", rand.Int())
		keyName := mapName + "_key"

		distMap := counter.GetMap(ctx, mapName)
		_ = distMap.Set(ctx, keyName, 0)

		b.StartTimer() // Important!

		// MEASURE
		counter.ExecuteCounterWithoutBlocking(ctx, mapName, keyName)
	}
}

func BenchmarkCounterWithOptimisticBlocking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// GIVEN
		ctx := context.TODO()
		counter := CreateDao(ctx)

		rand.Seed(time.Now().Unix())
		mapName := fmt.Sprintf("optimistic_blocking_%d", rand.Int())
		keyName := mapName + "_key"

		distMap := counter.GetMap(ctx, mapName)
		_ = distMap.Set(ctx, keyName, 0)

		b.StartTimer() // Important!

		// MEASURE
		counter.ExecuteCounterWithOptimisticBlocking(ctx, mapName, keyName, false)
	}
}

func BenchmarkCounterWithAtomicLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// GIVEN
		ctx := context.TODO()
		counter := CreateDao(ctx)

		rand.Seed(time.Now().Unix())
		name := fmt.Sprintf("atomic_counter_%d", rand.Int())

		atomic := counter.GetAtomicLong(ctx, name)
		_ = atomic.Set(ctx, 0)

		b.StartTimer() // Important!

		// MEASURE
		counter.ExecuteCounterWithAtomicLong(ctx, name)
	}
}
