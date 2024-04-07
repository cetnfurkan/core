package config

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func DurationHook() mapstructure.DecodeHookFuncType {
	// Wrapped in a function call to add optional input parameters (eg. separator)
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data interface{}, // raw data
	) (interface{}, error) {

		if t != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			duration, err := time.ParseDuration(data.(string))
			if err != nil {
				return data, nil
			}

			return duration * time.Millisecond, nil

		case reflect.Float64:
			return time.Duration(data.(float64) * float64(time.Millisecond)), nil

		case reflect.Int:
			return time.Duration(data.(int) * int(time.Millisecond)), nil

		default:
			return data, nil

		}
	}
}
