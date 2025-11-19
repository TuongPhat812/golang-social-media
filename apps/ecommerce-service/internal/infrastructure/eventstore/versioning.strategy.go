package eventstore

import (
	"encoding/json"
	"fmt"
)

// EventVersioningStrategy handles event versioning and migration
type EventVersioningStrategy struct {
	// Map of event type to current version
	currentVersions map[string]int
	// Map of event type to migration functions
	migrations map[string]map[int]func(map[string]interface{}) (map[string]interface{}, error)
}

// NewEventVersioningStrategy creates a new EventVersioningStrategy
func NewEventVersioningStrategy() *EventVersioningStrategy {
	return &EventVersioningStrategy{
		currentVersions: make(map[string]int),
		migrations:      make(map[string]map[int]func(map[string]interface{}) (map[string]interface{}, error)),
	}
}

// RegisterVersion registers the current version for an event type
func (s *EventVersioningStrategy) RegisterVersion(eventType string, version int) {
	s.currentVersions[eventType] = version
}

// RegisterMigration registers a migration function for an event type and version
func (s *EventVersioningStrategy) RegisterMigration(eventType string, fromVersion, toVersion int, migration func(map[string]interface{}) (map[string]interface{}, error)) {
	if s.migrations[eventType] == nil {
		s.migrations[eventType] = make(map[int]func(map[string]interface{}) (map[string]interface{}, error))
	}
	// Store migration as a key that represents the target version
	key := toVersion*1000 + fromVersion // Simple encoding: toVersion*1000 + fromVersion
	s.migrations[eventType][key] = migration
}

// MigrateEvent migrates an event from one version to another
func (s *EventVersioningStrategy) MigrateEvent(eventType string, fromVersion, toVersion int, payload map[string]interface{}) (map[string]interface{}, error) {
	if fromVersion == toVersion {
		return payload, nil
	}

	if fromVersion > toVersion {
		return nil, fmt.Errorf("cannot downgrade event %s from version %d to %d", eventType, fromVersion, toVersion)
	}

	// Apply migrations step by step
	currentPayload := payload
	currentVersion := fromVersion

	for currentVersion < toVersion {
		targetVersion := currentVersion + 1
		key := targetVersion*1000 + currentVersion

		migration, exists := s.migrations[eventType][key]
		if !exists {
			// No migration found, try to find a generic migration
			// For now, return error if no migration found
			return nil, fmt.Errorf("no migration found for event %s from version %d to %d", eventType, currentVersion, targetVersion)
		}

		migratedPayload, err := migration(currentPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate event %s from version %d to %d: %w", eventType, currentVersion, targetVersion, err)
		}

		currentPayload = migratedPayload
		currentVersion = targetVersion
	}

	return currentPayload, nil
}

// GetCurrentVersion returns the current version for an event type
func (s *EventVersioningStrategy) GetCurrentVersion(eventType string) int {
	if version, ok := s.currentVersions[eventType]; ok {
		return version
	}
	return 1 // Default version
}

// UpgradeEvent upgrades an event to the current version
func (s *EventVersioningStrategy) UpgradeEvent(eventType string, version int, payloadJSON string) (map[string]interface{}, error) {
	// Parse payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, err
	}

	currentVersion := s.GetCurrentVersion(eventType)
	if version >= currentVersion {
		return payload, nil
	}

	// Migrate to current version
	return s.MigrateEvent(eventType, version, currentVersion, payload)
}

