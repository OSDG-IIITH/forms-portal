package utils

// global variable, use after calling utils.LoadConfig()
var Config config

func EmptyArrayIfNull[T any](data []T) []T {
	if data == nil {
		return []T{}
	} else {
		return data
	}
}
