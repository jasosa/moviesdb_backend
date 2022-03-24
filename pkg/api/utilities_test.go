package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorJSON_NoStatusCode(t *testing.T) {
	expectedStatus := http.StatusBadRequest
	expectedBody := "{\"error\":{\"message\":\"a test error\"}}"

	rr := httptest.NewRecorder()

	ErrorJSON(rr, errors.New("a test error"))
	if rr.Code != expectedStatus {
		t.Errorf("Error json assigned wrong status code: got %v want %v",
			rr.Code, expectedStatus)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("Error json assigned wrong body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestErrorJSON_StatusCodeAssigned(t *testing.T) {
	expectedErrorMessage := "status forbidden test error message"
	expectedStatus := http.StatusForbidden
	expectedBody := "{\"error\":{\"message\":\"" + expectedErrorMessage + "\"}}"

	rr := httptest.NewRecorder()

	ErrorJSON(rr, errors.New(expectedErrorMessage), http.StatusForbidden)
	if rr.Code != expectedStatus {
		t.Errorf("Error json assigned wrong status code: got %v want %v",
			rr.Code, expectedStatus)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("Error json assigned wrong body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

type testStruct struct {
	FieldA string `json:"fieldA"`
	FieldB int    `json:"fieldB"`
}

func TestWriteJSON(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedBody := "{\"mywrapper\":{\"fieldA\":\"field 1\",\"fieldB\":2}}"

	rr := httptest.NewRecorder()
	testStruct := &testStruct{
		FieldA: "field 1",
		FieldB: 2,
	}

	WriteJSON(rr, http.StatusOK, testStruct, "mywrapper")
	if rr.Code != expectedStatus {
		t.Errorf("Error json assigned wrong status code: got %v want %v",
			rr.Code, expectedStatus)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("Error json assigned wrong body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}
