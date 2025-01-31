package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	errorLogger *log.Logger
	once        sync.Once
)

func getErrorLogger() *log.Logger {
	once.Do(func() {
		file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open error.log file: %v", err)
		}
		errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return errorLogger
}
func LogError(err error) {
	logger := getErrorLogger()
	logger.Println(err.Error())
}

func WriteJSONError(w http.ResponseWriter, statusCode int, errMsg string) {
	LogError(errors.New(errMsg))
	w.WriteHeader(statusCode)
	resp := ErrorResponse{Error: errMsg}
	err := jsonNewEncoder(w).Encode(resp)
	if err != nil {
		LogError(err)
	}
}

func jsonNewEncoder(w http.ResponseWriter) *jsonEncoderWrapper {
	return &jsonEncoderWrapper{w: w}
}

type jsonEncoderWrapper struct {
	w http.ResponseWriter
}

func (j *jsonEncoderWrapper) Encode(v interface{}) error {
	j.w.Header().Set("Content-Type", "application/json")
	return EncodeAsJSON(v, j.w)
}

func EncodeAsJSON(v interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
