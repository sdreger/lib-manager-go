package response

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderDataJSON(t *testing.T) {
	w := httptest.NewRecorder()
	statusCode := http.StatusOK

	err := RenderDataJSON(w, statusCode, map[string]string{"title": "Hello", "subtitle": "World"})
	result := w.Result()
	defer result.Body.Close()

	if assert.NoError(t, err, "should have no error during render data") {
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
		assert.Equal(t, statusCode, result.StatusCode)
		bytes, err := io.ReadAll(result.Body)
		assert.NoError(t, err, "should read body")
		assert.Equal(t, "{\"data\":{\"subtitle\":\"World\",\"title\":\"Hello\"}}", string(bytes))
	}
}

func TestRenderErrorJSON(t *testing.T) {
	w := httptest.NewRecorder()
	statusCode := http.StatusBadRequest

	err := RenderErrorJSON(w, http.StatusBadRequest, []APIError{{
		Field:   "title",
		Message: "can not be empty",
	}})
	result := w.Result()
	defer result.Body.Close()

	if assert.NoError(t, err, "should have no error during render data") {
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
		assert.Equal(t, statusCode, result.StatusCode)
		bytes, err := io.ReadAll(result.Body)
		assert.NoError(t, err, "should read body")
		assert.Equal(t, "{\"errors\":[{\"field\":\"title\",\"message\":\"can not be empty\"}]}", string(bytes))
	}
}

func TestRenderJSONWithHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	statusCode := http.StatusOK

	err := RenderJSONWithHeaders(w, statusCode, map[string]string{"title": "Hello", "subtitle": "World"}, http.Header{
		"Content-Disposition": []string{"attachment; filename=\"Hello\""},
	})
	result := w.Result()
	defer result.Body.Close()

	if assert.NoError(t, err, "should have no error during render data") {
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
		assert.Equal(t, "attachment; filename=\"Hello\"", result.Header.Get("Content-Disposition"))
		assert.Equal(t, statusCode, result.StatusCode)
		bytes, err := io.ReadAll(result.Body)
		assert.NoError(t, err, "should read body")
		assert.Equal(t, "{\"subtitle\":\"World\",\"title\":\"Hello\"}", string(bytes))
	}
}

func TestRenderJSONWithHeaders_MarshallingError(t *testing.T) {
	w := httptest.NewRecorder()
	statusCode := http.StatusOK

	unsupportedType := make(chan int)
	err := RenderDataJSON(w, statusCode, unsupportedType)
	result := w.Result()
	defer result.Body.Close()
	if assert.Error(t, err, "should have an error during render data") {
		assert.Equal(t, "json: unsupported type: chan int", err.Error())
	}
}

//func TestRenderJSONWithHeaders_WriteError(t *testing.T) {
//
//	recorder := httptest.NewRecorder()
//	recorder.
//	statusCode := http.StatusOK
//
//	err := RenderDataJSON(recorder, statusCode, "test data")
//	//result := w.Result()
//	//defer result.Body.Close()
//	if assert.Error(t, err, "should have an error during render data") {
//		assert.Equal(t, "json: unsupported type: chan int", err.Error())
//	}
//}
