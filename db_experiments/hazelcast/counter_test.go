package hazelcast

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
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
