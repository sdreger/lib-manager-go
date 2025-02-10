package response

import (
	"encoding/json"
	"net/http"
)

const (
	dataWrapper  = "data"
	errorWrapper = "errors"
)

type APIError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func RenderErrorJSON(w http.ResponseWriter, statusCode int, err []APIError) error {
	return RenderJSONWithHeaders(w, statusCode, map[string]interface{}{errorWrapper: err}, nil)
}

func RenderDataJSON(w http.ResponseWriter, statusCode int, data any) error {
	return RenderJSONWithHeaders(w, statusCode, map[string]interface{}{dataWrapper: data}, nil)
}

func RenderJSONWithHeaders(w http.ResponseWriter, statusCode int, data any, headers http.Header) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
