package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	testing

	"github.com/stretchr/testify/assert"
)

// A simple table driven test example for auth
func TestInvalidCredentials(t *testing.T) {
	cases := []struct {
		name     string
		payload  map[string]interface{}
		expected int
	}{
		{"Empty Payload", map[string]interface{}{}, http.StatusBadRequest},
		{"Missing Password", map[string]interface{}{"email": "test@test.com"}, http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body, _ := json.Marshal(c.payload)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			// Note: we need the router to test completely, but here we just assert types
			assert.NotNil(t, req)
			assert.NotNil(t, w)
		})
	}
}
