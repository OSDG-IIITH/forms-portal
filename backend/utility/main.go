package utils

import (
	"github.com/go-playground/validator/v10"
)

// global variable, use after calling utils.LoadConfig()
var Config config

// global variable, use after calling utils.LoadValidator()
var Validate *validator.Validate

func EmptyArrayIfNull[T any](data []T) []T {
	if data == nil {
		return []T{}
	} else {
		return data
	}
}
