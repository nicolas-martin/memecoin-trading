package blockchain

import (
	"context"
	"fmt"
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

	// Verify provider was registered
	provider, err := pm.getHealthyProvider(context.Background(), NetworkSolana)
	assert.NoError(t, err)
	assert.Same(t, mockProvider, provider)
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

	// Set up test operation that fails for first provider but succeeds for second
	var calledProvider Provider
	testOp := func(p Provider) error {
		calledProvider = p
		if mp, ok := p.(*MockProvider); ok && mp == mockProvider1 {
			return fmt.Errorf("simulated error")
		}
		return nil
	}

	// Test fallback behavior
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)
	assert.Same(t, mockProvider2, calledProvider, "Should have fallen back to second provider")
}

func TestProviderManagerRateLimiting(t *testing.T) {
	pm := NewProviderManager()

	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	// Register provider with low request limit
	err := pm.RegisterProvider(mockProvider, ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  2,           // Only allow 2 requests per window
		WindowDuration:     time.Second, // Short window for testing
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	// Set up test operation
	requestCount := 0
	testOp := func(p Provider) error {
		requestCount++
		return nil
	}

	// First request should succeed
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)

	// Second request should succeed
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)
	assert.Equal(t, 2, requestCount)

	// Third request should fail due to rate limiting
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no available providers")
	assert.Equal(t, 2, requestCount, "No additional requests should have been made")

	// Wait for rate limit window to expire
	time.Sleep(time.Second)

	// Should be able to make requests again
	err = pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.NoError(t, err)
	assert.Equal(t, 3, requestCount)
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
	requestCount := 0
	failingOp := func(p Provider) error {
		requestCount++
		return fmt.Errorf("simulated error")
	}

	// First error
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated error")
	assert.Equal(t, 1, requestCount)

	// Second error should mark the provider as unhealthy
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated error")
	assert.Equal(t, 2, requestCount)

	// Third request should fail immediately as no healthy providers are available
	err = pm.executeWithFallback(context.Background(), NetworkSolana, failingOp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no available providers")
	assert.Equal(t, 2, requestCount, "No additional requests should have been made")

	// Verify provider is marked as unhealthy
	provider, err := pm.getHealthyProvider(context.Background(), NetworkSolana)
	assert.Error(t, err)
	assert.Nil(t, provider)
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

	// Track order of provider usage
	var usedProviders []Provider
	testOp := func(p Provider) error {
		usedProviders = append(usedProviders, p)
		return fmt.Errorf("simulated error") // Force fallback to next provider
	}

	// Execute operation that will try all providers
	err := pm.executeWithFallback(context.Background(), NetworkSolana, testOp)
	assert.Error(t, err)

	// Verify providers were tried in priority order
	assert.Equal(t, 3, len(usedProviders))
	assert.Same(t, mockProvider2, usedProviders[0], "Highest priority provider should be tried first")
	assert.Same(t, mockProvider3, usedProviders[1], "Medium priority provider should be tried second")
	assert.Same(t, mockProvider1, usedProviders[2], "Lowest priority provider should be tried last")
}

func TestProviderManagerNoProviders(t *testing.T) {
	pm := NewProviderManager()

	// Try to get provider for network with no providers
	provider, err := pm.getHealthyProvider(context.Background(), NetworkSolana)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no providers registered")
	assert.Nil(t, provider)

	// Try to execute operation
	err = pm.executeWithFallback(context.Background(), NetworkSolana, func(p Provider) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no providers registered")
}
