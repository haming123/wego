package wego

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TesWriterJson(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	_, c := CreateTestContext(w, req)

	var user User
	user.Name = "lisi"
	user.Age = 12
	c.WriteJSON(200, user)
}
