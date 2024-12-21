package leaderboard

import (
	"context"
	"testing"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/mocks"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTopTraders(t *testing.T) {
	mockDB := &mocks.MockLeaderboardRepository{}
	mockCache := &mocks.MockCache{}
	service := NewService(mockDB, mockCache)

	ctx := context.Background()
	duration := 24 * time.Hour
	timeframe := duration.String()

	testEntries := []models.LeaderboardEntry{
		{
			Username:         "trader1",
			Profit:           1000.0,
			ProfitPercentage: 10.0,
			Rank:             1,
		},
	}

	t.Run("returns cached entries when available", func(t *testing.T) {
		mockCache.On("GetLeaderboard", ctx, timeframe).Return(testEntries, nil)

		entries, err := service.GetTopTraders(ctx, duration)

		assert.NoError(t, err)
		assert.Equal(t, testEntries, entries)
		mockDB.AssertNotCalled(t, "GetTopTraders")
	})

	t.Run("fetches from DB when cache misses", func(t *testing.T) {
		mockCache.On("GetLeaderboard", ctx, timeframe).Return(nil, nil)
		mockDB.On("GetTopTraders", ctx, duration, 100).Return(testEntries, nil)
		mockCache.On("SetLeaderboard", ctx, timeframe, testEntries).Return(nil)

		entries, err := service.GetTopTraders(ctx, duration)

		assert.NoError(t, err)
		assert.Equal(t, testEntries, entries)
		mockDB.AssertCalled(t, "GetTopTraders", ctx, duration, 100)
	})
}
