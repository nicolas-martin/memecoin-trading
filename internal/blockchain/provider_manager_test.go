package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProviderManagerRegistration(t *testing.T) {
	pm := NewProviderManager()
	assert.NotNil(t, pm)

	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	config := ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	}

	err := pm.RegisterProvider(mockProvider, config)
	assert.NoError(t, err)
}

func TestProviderManagerFallback(t *testing.T) {
	pm := NewProviderManager()

	// Create two providers
	mockProvider1 := new(MockProvider)
	mockProvider1.On("Network").Return(NetworkSolana)
	mockProvider2 := new(MockProvider)
	mockProvider2.On("Network").Return(NetworkSolana)

	// Register providers with different priorities
	err := pm.RegisterProvider(mockProvider1, ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	err = pm.RegisterProvider(mockProvider2, ProviderConfig{
		Priority:           2,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	// Set up test operation
	testOp := func(p Provider) error {
		return nil
	}

	// Test successful operation
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)
}

func TestProviderManagerRateLimiting(t *testing.T) {
	pm := NewProviderManager()

	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	// Register provider with low request limit
	err := pm.RegisterProvider(mockProvider, ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  2, // Only allow 2 requests per window
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	// Set up test operation
	testOp := func(p Provider) error {
		return nil
	}

	// First request should succeed
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)

	// Second request should succeed
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)

	// Third request should fail due to rate limiting
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no available providers")
}

func TestProviderManagerHealthCheck(t *testing.T) {
	pm := NewProviderManager()

	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	// Register provider with low error threshold
	err := pm.RegisterProvider(mockProvider, ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 2, // Mark as unhealthy after 2 consecutive errors
	})
	assert.NoError(t, err)

	// Set up failing operation
	failingOp := func(p Provider) error {
		return assert.AnError
	}

	// First error
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)

	// Second error should mark the provider as unhealthy
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)

	// Third request should fail immediately as no healthy providers are available
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no available providers for network solana")
}

func TestProviderManagerPriority(t *testing.T) {
	pm := NewProviderManager()

	// Create three providers with different priorities
	mockProvider1 := new(MockProvider)
	mockProvider1.On("Network").Return(NetworkSolana)
	mockProvider2 := new(MockProvider)
	mockProvider2.On("Network").Return(NetworkSolana)
	mockProvider3 := new(MockProvider)
	mockProvider3.On("Network").Return(NetworkSolana)

	// Register providers with different priorities
	configs := []struct {
		provider Provider
		priority int
	}{
		{mockProvider1, 3}, // Lowest priority
		{mockProvider2, 1}, // Highest priority
		{mockProvider3, 2}, // Medium priority
	}

	for _, c := range configs {
		err := pm.RegisterProvider(c.provider, ProviderConfig{
			Priority:           c.priority,
			RequestsPerWindow:  100,
			WindowDuration:     time.Minute,
			HealthCheckPeriod:  time.Minute,
			MaxConsecutiveErrs: 3,
		})
		assert.NoError(t, err)
	}

	// Get a provider and verify it's the highest priority one
	provider, err := pm.getHealthyProvider(context.Background(), NetworkSolana)
	assert.NoError(t, err)
	assert.Equal(t, mockProvider2, provider)
}
