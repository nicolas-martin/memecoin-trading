package blockchain

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProviderStatus represents the current status of a provider
type ProviderStatus int

const (
	ProviderStatusHealthy ProviderStatus = iota
	ProviderStatusDegraded
	ProviderStatusUnhealthy
)

// ProviderHealth tracks the health and rate limit status of a provider
type ProviderHealth struct {
	status           ProviderStatus
	lastChecked      time.Time
	consecutiveErrs  int
	rateLimitExpiry  time.Time
	requestsInWindow int
	mu               sync.RWMutex
}

// ProviderConfig holds configuration for a provider
type ProviderConfig struct {
	Priority           int           // Lower number means higher priority
	RequestsPerWindow  int           // Number of requests allowed in the time window
	WindowDuration     time.Duration // Duration of the rate limiting window
	HealthCheckPeriod  time.Duration // How often to check provider health
	MaxConsecutiveErrs int           // Number of consecutive errors before marking as unhealthy
}

// providerEntry represents a provider and its associated metadata
type providerEntry struct {
	provider Provider
	health   *ProviderHealth
	config   ProviderConfig
}

// ProviderManager manages multiple providers per network with fallback and load balancing
type ProviderManager struct {
	providers map[Network][]providerEntry
	mu        sync.RWMutex
}

// NewProviderManager creates a new provider manager
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[Network][]providerEntry),
	}
}

// RegisterProvider registers a new provider with its configuration
func (pm *ProviderManager) RegisterProvider(provider Provider, config ProviderConfig) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	network := provider.Network()
	health := &ProviderHealth{
		status:      ProviderStatusHealthy,
		lastChecked: time.Now(),
	}

	entry := providerEntry{
		provider: provider,
		health:   health,
		config:   config,
	}

	providers := pm.providers[network]
	pm.providers[network] = append(providers, entry)

	// Sort providers by priority
	pm.sortProviders(network)

	return nil
}

// getHealthyProvider returns the highest priority healthy provider for a network
func (pm *ProviderManager) getHealthyProvider(_ context.Context, network Network) (Provider, error) {
	pm.mu.RLock()
	providers, exists := pm.providers[network]
	pm.mu.RUnlock()

	if !exists || len(providers) == 0 {
		return nil, fmt.Errorf("no providers registered for network %s", network)
	}

	for _, entry := range providers {
		if pm.isProviderAvailable(entry) {
			return entry.provider, nil
		}
	}

	return nil, fmt.Errorf("no healthy providers available for network %s", network)
}

// isProviderAvailable checks if a provider is healthy and not rate limited
func (pm *ProviderManager) isProviderAvailable(entry providerEntry) bool {
	entry.health.mu.RLock()
	defer entry.health.mu.RUnlock()

	now := time.Now()

	// Check if rate limited
	if now.Before(entry.health.rateLimitExpiry) {
		return false
	}

	// Check if within rate limit window
	if entry.health.requestsInWindow >= entry.config.RequestsPerWindow {
		entry.health.rateLimitExpiry = now.Add(entry.config.WindowDuration)
		return false
	}

	return entry.health.status == ProviderStatusHealthy
}

// recordSuccess records a successful request for a provider
func (pm *ProviderManager) recordSuccess(entry providerEntry) {
	entry.health.mu.Lock()
	defer entry.health.mu.Unlock()

	entry.health.consecutiveErrs = 0
	entry.health.requestsInWindow++
}

// recordError records a failed request for a provider
func (pm *ProviderManager) recordError(entry providerEntry) {
	entry.health.mu.Lock()
	defer entry.health.mu.Unlock()

	entry.health.consecutiveErrs++
	if entry.health.consecutiveErrs >= entry.config.MaxConsecutiveErrs {
		entry.health.status = ProviderStatusUnhealthy
	}
}

// sortProviders sorts providers by priority
func (pm *ProviderManager) sortProviders(network Network) {
	providers := pm.providers[network]
	// Sort providers by priority (lower number means higher priority)
	for i := 0; i < len(providers)-1; i++ {
		for j := i + 1; j < len(providers); j++ {
			if providers[i].config.Priority > providers[j].config.Priority {
				providers[i], providers[j] = providers[j], providers[i]
			}
		}
	}
}

// executeWithFallback executes an operation with automatic fallback
func (pm *ProviderManager) executeWithFallback(_ context.Context, network Network, op func(Provider) error) error {
	pm.mu.RLock()
	providers := pm.providers[network]
	pm.mu.RUnlock()

	if len(providers) == 0 {
		return fmt.Errorf("no providers registered for network %s", network)
	}

	var lastErr error
	for _, entry := range providers {
		if !pm.isProviderAvailable(entry) {
			continue
		}

		err := op(entry.provider)
		if err == nil {
			pm.recordSuccess(entry)
			return nil
		}

		pm.recordError(entry)
		lastErr = err
	}

	if lastErr != nil {
		return fmt.Errorf("all providers failed: %w", lastErr)
	}
	return fmt.Errorf("no available providers for network %s", network)
}
