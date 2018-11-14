package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetByIdHandler(t *testing.T) {
	tt := []struct {
		variable   string
		shouldPass bool
	}{
		{"de", false},
		{"en", true},
		{"es", true},
	}
	for _, tc := range tt {
		path := fmt.Sprintf("/content/%s", tc.variable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on routeVariable %s: got %v want %v",
				tc.variable, rr.Code, http.StatusOK)
		}

	}
}

func TestDeleteByIdHandler(t *testing.T) {
	for i := 0; i < 2; i++ {
		path := fmt.Sprintf("/content/%s/%d", "en", 1)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler should have return bad status on second time delete")
		}
		// TODO: Fix data management then, this will work. right now every requests creates a new DB
		//if rr.Code == http.StatusNotFound && i == 0 {
		//	t.Errorf("handler should have return OK on first time delete")
		//}
	}

}

func TestGetAllHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/all", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if typ := rr.Header().Get("Content-Type"); typ != "application/json" {
		t.Errorf("Expected content-type application/json, found %s", typ)
	}

	var responseObject JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &responseObject)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body is what we expect.
	if len(responseObject.Data) != 4 {
		t.Errorf("handler returned unexpected length: got %d expected %d",
			rr.Body, 4)
	}
}

func TestAddInstanceHandler(t *testing.T) {
	tt := []struct {
		lang       string
		id         uint
		q          string
		a          string
		shouldPass bool
	}{
		{"de", 3, "test1", "test2", true},
		{"en", 2, "hi there", "you too", true},
	}
	for _, tc := range tt {
		path := fmt.Sprintf("/content/%s/%d/%s/%s", tc.lang, tc.id, tc.q, tc.a)
		req, err := http.NewRequest("PUT", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler did not pass on PUT %s", path)
		}

	}
}
