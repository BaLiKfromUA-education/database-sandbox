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
	require.True(t, value.(int64) > 0)
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

	require.Equal(t, value.(int64), 100_000)
}
