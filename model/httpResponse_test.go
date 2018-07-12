package model

import (
	"testing"
)

// TestGetHTTPResponse : Test de la fonction GetHTTPResponse
func TestGetHTTPResponse(t *testing.T) {
	response := GetHTTPResponse(200, "Success", nil)
	responseValid := HTTPResponse{200, "Success", nil}

	if response != responseValid {
		t.Errorf("GetHTTPResponse - got %+v: , want: %+v.", response, responseValid)
	}

	response = GetHTTPResponse(404, "Not found", 15)
	responseValid = HTTPResponse{404, "Not found", 15}

	if response != responseValid {
		t.Errorf("GetHTTPResponse - got %+v: , want: %+v.", response, responseValid)
	}
}
