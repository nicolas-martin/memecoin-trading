package mocks

import (
	"context"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetLeaderboard(ctx context.Context, timeframe string) ([]models.LeaderboardEntry, error) {
	args := m.Called(ctx, timeframe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LeaderboardEntry), args.Error(1)
}

func (m *MockCache) SetLeaderboard(ctx context.Context, timeframe string, entries []models.LeaderboardEntry) error {
	args := m.Called(ctx, timeframe, entries)
	return args.Error(0)
}

// Add more mock methods as needed...
