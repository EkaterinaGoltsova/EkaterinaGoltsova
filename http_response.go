package main

import (
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

type responseSuccess struct {
	Data   string      `json:"data"`
	Meta   interface{} `json:"meta"`
	Errors []string    `json:"errors"`
}

type responseError struct {
	Errors []string `json:"errors"`
}

func getErrorResponse(err error) responseError {
	return responseError{Errors: getErrors(err)}
}

func getErrors(err error) []string {
	var errors []string
	if err != nil {
		for _, err := range strings.Split(err.Error(), hierr.BranchDelimiter) {
			errors = append(errors, strings.TrimSpace(err))
		}
	}

	return errors
}

func getSuccessResponse(data string, meta interface{}, err error) responseSuccess {
	return responseSuccess{
		Data:   data,
		Meta:   meta,
		Errors: getErrors(err),
	}
}
