package factory

import (
	"time"

	"github.com/TerraDharitri/drt-go-sdk/core/http"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/checkers"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/config"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/executors"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/interactors"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/monitor"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/parsers"
)

const timeBetweenBLSKeysFetch = time.Second

// NewBLSKeysMonitor will create a BLS keys monitor based on the configs & other internal components
func NewBLSKeysMonitor(
	cfg config.BLSKeysMonitorConfig,
	snoozeConfig config.AlarmSnoozeConfig,
	notifiersHandler OutputNotifiersHandler,
	statusHandler executors.StatusHandler,
) (Monitor, error) {
	parser := parsers.NewListParser()
	intentitiesHolder, err := parser.ParseFile(cfg.ListFile)
	if err != nil {
		return nil, err
	}

	ratingsChecker, err := checkers.NewBLSRatingsChecker(intentitiesHolder.BlsHexKeys, cfg.Name, cfg.AlarmDeltaRatingDrop)
	if err != nil {
		return nil, err
	}

	httpWrapper := http.NewHttpClientWrapper(nil, cfg.ApiURL)
	interactor, err := interactors.NewValidatorStatisticsInteractor(httpWrapper)
	if err != nil {
		return nil, err
	}

	fetcher, err := interactors.NewBLSKeysFetcher(httpWrapper, intentitiesHolder.Addresses, timeBetweenBLSKeysFetch)
	if err != nil {
		return nil, err
	}

	blsKeysFilter, err := NewBLSKeysFilter(snoozeConfig)
	if err != nil {
		return nil, err
	}

	argsExecutor := executors.ArgsBLSKeysExecutor{
		OutputNotifiersHandler:     notifiersHandler,
		RatingsChecker:             ratingsChecker,
		ValidatorStatisticsQuerier: interactor,
		BlsKeysFetcher:             fetcher,
		StatusHandler:              statusHandler,
		BLSKeysFilter:              blsKeysFilter,
		Name:                       cfg.Name,
		ExplorerURL:                cfg.ExplorerURL,
	}

	executor, err := executors.NewBLSKeysExecutor(argsExecutor)
	if err != nil {
		return nil, err
	}

	return monitor.NewBLSKeysMonitor(
		executor,
		time.Duration(cfg.PollingIntervalInSeconds)*time.Second,
		cfg.Name)
}
