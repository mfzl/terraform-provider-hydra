package main

import (
	"net/http"

	"github.com/pkg/errors"
)

type validStrOptions map[string]bool

func (v validStrOptions) keys() []string {
	keys := []string{}

	for k := range v {
		keys = append(keys, k)
	}

	return keys
}

func toStringSlice(vals []interface{}) []string {
	strSlice := []string{}

	for _, v := range vals {
		if sv, ok := v.(string); ok {
			strSlice = append(strSlice, sv)
		}
	}

	return strSlice
}

func httpOk(code int) bool {
	if code >= http.StatusOK && code < 299 {
		return true
	}

	return false
}

func httpStatusErr(code int) error {
	if !httpOk(code) {
		return errors.Errorf("unxpected HTTP status code %d", code)
	}

	return nil
}
