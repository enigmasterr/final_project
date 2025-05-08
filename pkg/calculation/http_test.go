package calculation_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"net/http/httptest"

	"github.com/enigmasterr/final_project/internal/application"
)

func TestRequestHandlerSuccessCase(t *testing.T) {
	// expected := "Hello John"
	// req := httptest.NewRequest(http.MethodGet, "/greet?name=John", nil)
	// w := httptest.NewRecorder()
	validData := []byte(`{"expression": "2+2*7"}`)
	expectedData := 16.0
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
	// var response map[string]string
	// json.Unmarshal(w.Body.Bytes(), &response)
	// if response["status"] != "success" || response["message"] != "Data processed successfully" {
	// 	t.Errorf("Unexpected response body: %v", response)
	// }

	defer res.Body.Close()
	// data, err := io.ReadAll(res.Body)
	// type ResStr struct {
	// 	Result string `json:"result"`
	// }
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	// if err != nil {
	// 	t.Errorf("Error: %v", err)
	// }
	rez, _ := strconv.ParseFloat(response["result"], 64)
	if rez != expectedData {
		t.Errorf("Expected %f but got %f", expectedData, rez)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("wrong status code")
	}
}

func TestRequestHandlerBadRequestCase(t *testing.T) {
	notvalidData := []byte(`{"expression": "(2+2*7"}`)
	expectedAns := "Expression is not valid"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(notvalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, w.Code)
	}
	// var response map[string]string
	// json.Unmarshal(w.Body.Bytes(), &response)
	// if response["status"] != "success" || response["message"] != "Data processed successfully" {
	// 	t.Errorf("Unexpected response body: %v", response)
	// }

	defer res.Body.Close()
	// data, err := io.ReadAll(res.Body)
	// type ResStr struct {
	// 	Result string `json:"result"`
	// }
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	// if err != nil {
	// 	t.Errorf("Error: %v", err)
	// }

	if response["error"] != expectedAns {
		t.Errorf("Expected %s but got %s", expectedAns, response["error"])
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("wrong status code")
	}
}

func TestRequestHandlerUnprocessableEntityCase(t *testing.T) {
	notvalidData := []byte(`{"expression": "yuyi2+2*7"}`)
	expectedAns := "Expression is not valid"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(notvalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %v, got %v", http.StatusUnprocessableEntity, w.Code)
	}
	// var response map[string]string
	// json.Unmarshal(w.Body.Bytes(), &response)
	// if response["status"] != "success" || response["message"] != "Data processed successfully" {
	// 	t.Errorf("Unexpected response body: %v", response)
	// }

	defer res.Body.Close()
	// data, err := io.ReadAll(res.Body)
	// type ResStr struct {
	// 	Result string `json:"result"`
	// }
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	// if err != nil {
	// 	t.Errorf("Error: %v", err)
	// }

	if response["error"] != expectedAns {
		t.Errorf("Expected %s but got %s", expectedAns, response["error"])
	}

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("wrong status code")
	}
}
