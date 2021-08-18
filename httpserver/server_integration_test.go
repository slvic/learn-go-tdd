package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/vkenrik117/httpserver"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := httpserver.NewInMemoryPlayerStore()
	server := httpserver.NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
