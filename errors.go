/*
apitinyfile errors
*/

package main

import "net/http"

type jsonError struct {
	Error string `json:"error"`
}

func unauthError(errorStr string) (int, jsonError) {
	if errorStr == "" {
		errorStr = "Unauthorized"
	}
	return http.StatusUnauthorized, jsonError{Error: errorStr}
}

func badRequestError(errorStr string) (int, jsonError) {
	if errorStr == "" {
		errorStr = "Invalid request"
	}
	return http.StatusBadRequest, jsonError{Error: errorStr}
}

func notFoundError(errorStr string) (int, jsonError) {
	if errorStr == "" {
		errorStr = "Not found"
	}
	return http.StatusNotFound, jsonError{Error: errorStr}
}

func internalError(errorStr string) (int, jsonError) {
	if errorStr == "" {
		errorStr = "Internal error"
	}
	return http.StatusInternalServerError, jsonError{Error: errorStr}
}
