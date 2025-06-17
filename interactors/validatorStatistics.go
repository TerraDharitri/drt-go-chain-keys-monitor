package interactors

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-keys-monitor/core"
)

const validatorStatisticsURL = "validator/statistics"
const dataFieldName = "statistics"

var log = logger.GetOrCreate("interactors")

type validatorStatisticsInteractor struct {
	httpClientWrapper HTTPClientWrapper
}

// NewValidatorStatisticsInteractor can interact with the validator statistics API endpoint route
func NewValidatorStatisticsInteractor(httpClientWrapper HTTPClientWrapper) (*validatorStatisticsInteractor, error) {
	if check.IfNil(httpClientWrapper) {
		return nil, errNilHTTPClientWrapper
	}

	return &validatorStatisticsInteractor{
		httpClientWrapper: httpClientWrapper,
	}, nil
}

// Query will query the validator statistics API endpoint route
func (interactor *validatorStatisticsInteractor) Query(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
	log.Debug("validatorStatisticsInteractor.Query")
	bytes, statusCode, err := interactor.httpClientWrapper.GetHTTP(ctx, validatorStatisticsURL)
	if err != nil {
		return nil, err
	}
	if !core.IsHttpStatusCodeSuccess(statusCode) {
		return nil, fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}

	validatorStatisticsResponse := &core.ValidatorStatisticsResponse{}
	err = json.Unmarshal(bytes, validatorStatisticsResponse)
	if err != nil {
		return nil, err
	}

	statisticsMap := validatorStatisticsResponse.Data[dataFieldName]
	if statisticsMap == nil {
		return nil, errInvalidResponse
	}

	return statisticsMap, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (interactor *validatorStatisticsInteractor) IsInterfaceNil() bool {
	return interactor == nil
}
