package postgresql

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type PostgresCounterTestSuite struct {
	suite.Suite
	dao    *CounterDao
	userId int
	ctx    context.Context
}

func (suite *PostgresCounterTestSuite) SetupSuite() {
	log.Println("Start postgresql tests...")
	suite.ctx = context.TODO()
	suite.dao = CreateDao(suite.ctx)
	suite.userId = 42
}

func (suite *PostgresCounterTestSuite) SetupTest() {
	log.Println("Populate table...")
	suite.dao.InsertBaseRecord(suite.ctx, suite.userId)
	require.Equal(suite.T(), 0, suite.dao.GetResult(suite.ctx, suite.userId))
}

func (suite *PostgresCounterTestSuite) TearDownTest() {
	fmt.Println("Clean up table...")
	suite.dao.CleanUp(suite.ctx, suite.userId)
}

func (suite *PostgresCounterTestSuite) TestLostUpdate() {
	// GIVEN: lost update strategy

	// WHEN
	suite.dao.ExecuteLostUpdate(suite.ctx, suite.userId)

	// THEN
	result := suite.dao.GetResult(suite.ctx, suite.userId)
	log.Printf("Lost update result: %d", result)
	require.True(suite.T(), result > 0)
}

func (suite *PostgresCounterTestSuite) TestInPlaceUpdate() {
	// GIVEN: in-place update strategy

	// WHEN
	suite.dao.ExecuteInPlaceUpdate(suite.ctx, suite.userId)

	// THEN
	result := suite.dao.GetResult(suite.ctx, suite.userId)
	require.Equal(suite.T(), 100_000, result)
}

func TestPostgresCounterTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresCounterTestSuite))
}

func BenchmarkLostUpdate(b *testing.B) {
	// GIVEN
	ctx := context.TODO()
	id := 42

	dao := CreateDao(ctx)
	dao.CleanUp(ctx, id)
	dao.InsertBaseRecord(ctx, id)

	b.ResetTimer() // Important!

	// MEASURE
	dao.ExecuteLostUpdate(ctx, id)
}

func BenchmarkInPlaceUpdate(b *testing.B) {
	// GIVEN
	ctx := context.TODO()
	id := 42

	dao := CreateDao(ctx)
	dao.CleanUp(ctx, id)
	dao.InsertBaseRecord(ctx, id)

	b.ResetTimer() // Important!

	// MEASURE
	dao.ExecuteInPlaceUpdate(ctx, id)
}
