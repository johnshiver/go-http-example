package utils

import (
	"net/url"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	// from validator docs: use a single instance of Validate, it caches struct info
	v                 *validator.Validate
	initValidatorOnce sync.Once
)

func GetValidator() *validator.Validate {
	initValidatorOnce.Do(func() {
		v = validator.New()
	})
	return v
}

func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
