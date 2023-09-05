package service

import "fmt"

type PwServiceError struct {
	HTTPStatusCode int
	ErrorCode      int
	Log            bool
	Public         string
	Private        string
}

type ErrorOutput struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (d PwServiceError) Error() string {
	return fmt.Sprintf("Public Message : '%s' | Private message '%s'", d.Public, d.Private)
}
