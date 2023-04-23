package main

import (
	"net/http"
	"testing"

	"github.com/onahvictor/Snippet/internal/assert"
)

// end to end testing
func TestPing(t *testing.T) {
	app := newTestApplication(t)

	//This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the
	//duration of the test.
	ts := newTestServer(t, app.routes())
	//defer close shuts down this server when test is over
	defer ts.Close()

	//the network address thet the server is listening
	// is contained in the ts.URL field.
	sc, mh, body := ts.get(t, "/ping")

	expectedValue := "origin-when-cross-origin"
	assert.Equal(t, mh.Get("Referrer-Policy"), expectedValue)

	expectedValue = "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, mh.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, mh.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, mh.Get("X-Frame-Options"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, mh.Get("X-Content-Type-Options"), expectedValue)

	assert.Equal(t, body, "OK")
	assert.Equal(t, sc, http.StatusOK)

}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},

		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
