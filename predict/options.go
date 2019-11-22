package predict

import (
	"fmt"
	"github.com/posener/complete/v2"
	"strings"
)

// Option provides prediction through options pattern.
//
// Usage:
//
//  func(o ...predict.Option) {
//  	cfg := predict.Options(o)
//  	// use cfg.Predict...
//  }
type Option func(*Config)

// OptValues allows to set a desired set of valid values for the flag.
func OptValues(values ...string) Option {
	return OptPredictor(Set(values))
}

// OptPredictor allows to set a custom predictor.
func OptPredictor(p complete.Predictor) Option {
	return func(o *Config) { o.Predictor = p }
}

// OptCheck enforces the valid values on the predicted flag.
func OptCheck() Option {
	return func(o *Config) { o.check = true }
}

type Config struct {
	complete.Predictor
	check bool
}

func Options(os ...Option) Config {
	var op Config
	for _, f := range os {
		f(&op)
	}
	return op
}

func (c Config) Predict(prefix string) []string {
	if c.Predictor != nil {
		return c.Predictor.Predict(prefix)
	}
	return nil
}

func (c Config) Check(value string) error {
	predictions := c.Predictor.Predict(value)
	if !c.check || len(predictions) == 0 {
		return nil
	}
	for _, vv := range predictions {
		if value == vv {
			return nil
		}
	}
	return fmt.Errorf("not in allowed values: %s", strings.Join(predictions, ","))
}
