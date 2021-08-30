package poker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vkenrik117/poker"
)

var (
	dummyGame = &poker.GameSpy{}
)

func TestGETPlayers(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	},
		nil, nil)

	server, _ := poker.NewPlayerServer(store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.NewStubPlayerStore(
		map[string]int{},
		nil,
		nil,
	)
	server, _ := poker.NewPlayerServer(store, dummyGame)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusAccepted)

		winCalls := store.GetWinCalls()

		if len(winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(winCalls), 1)
		}

		if winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []poker.Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := poker.NewStubPlayerStore(nil, nil, wantedLeague)
		server, _ := poker.NewPlayerServer(store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertContentType(t, response, poker.JsonContentType)

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /playGame returns 200", func(t *testing.T) {
		server, err := poker.NewPlayerServer(&poker.StubPlayerStore{}, dummyGame)
		if err != nil {
			fmt.Printf("problem creating new player server: %v:", err)
		}

		request, err :=  newGameRequest()
		if err != nil {
			fmt.Printf("problem requesting a new game: %v:", err)
		}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
	})

	t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
		const tenMS = 10 * time.Millisecond
		var err error
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &poker.GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer func() {
			err = ws.Close()
		}()
		if err != nil {
			fmt.Printf("problem closing websocket: %v", err)
		}

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(tenMS)

		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func newGameRequest() (*http.Request, error) {
	return http.NewRequest(http.MethodGet, "/playGame", nil)
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}