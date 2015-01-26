package middleware

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/casualjim/go-swagger/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func terminator(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestRouterMiddleware(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	mw := context.RouterMiddleware(http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/api/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)

	methods := strings.Split(recorder.Header().Get("Allow"), ",")
	sort.Sort(sort.StringSlice(methods))
	assert.Equal(t, "GET,POST", strings.Join(methods, ","))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/nopets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	spec, api = petstore.NewRootAPI(t)
	context = NewContext(spec, api, nil)
	mw = context.RouterMiddleware(http.HandlerFunc(terminator))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)

	methods = strings.Split(recorder.Header().Get("Allow"), ",")
	sort.Sort(sort.StringSlice(methods))
	assert.Equal(t, "GET,POST", strings.Join(methods, ","))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/nopets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}
