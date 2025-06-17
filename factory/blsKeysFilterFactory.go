package factory

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-keys-monitor/config"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/executors"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/executors/disabled"
)

// NewBLSKeysFilter creates a new instance of type BLSKeysFilter
func NewBLSKeysFilter(cfg config.AlarmSnoozeConfig) (executors.BLSKeysFilter, error) {
	if !cfg.Enabled {
		return disabled.NewDisabledBLSKeysFilter(), nil
	}

	args := executors.ArgsBlsKeysTimeCache{
		MaxSnoozeEvents:     cfg.NumNotificationsForEachFaultyKey,
		CacheExpiration:     cfg.SnoozeTimeInSec,
		GetCurrentTimestamp: time.Now().Unix,
	}

	return executors.NewBLSKeysTimeCache(args)
}
