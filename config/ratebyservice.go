package config

import (
	"sync"
)

// RateByService stores the sampling rate per service. It is thread-safe, so
// one can read/write on it concurrently, using getters and setters.
type RateByService struct {
	rates map[string]float64
	mutex sync.RWMutex
}

// SetAll the sampling rate for all services. If a service/env is not
// in the map, then the entry is removed.
func (rbs *RateByService) SetAll(rates map[string]float64) {
	rbs.mutex.Lock()
	defer rbs.mutex.Unlock()

	if rbs.rates == nil {
		rbs.rates = make(map[string]float64, len(rates))
	}
	for k := range rbs.rates {
		if _, ok := rates[k]; !ok {
			delete(rbs.rates, k)
		}
	}
	for k, v := range rates {
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		rbs.rates[k] = v
	}
}

// GetAll returns all sampling rates for all services.
func (rbs *RateByService) GetAll() map[string]float64 {
	rbs.mutex.RLock()
	defer rbs.mutex.RUnlock()

	ret := make(map[string]float64, len(rbs.rates))
	for k, v := range rbs.rates {
		ret[k] = v
	}

	return ret
}
