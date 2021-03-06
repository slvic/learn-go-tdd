package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type GameSpy struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	FinishedCalled   bool
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	_, err := out.Write(g.BlindAlert)
	if err != nil {
		fmt.Printf("problem writing out at the Start: %v", err)
	}
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalledWith = winner
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips At %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, _ io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{At: at, Amount: amount})
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func NewStubPlayerStore(s map[string]int, w []string, l League) *StubPlayerStore {
	return &StubPlayerStore{
		scores:   s,
		winCalls: w,
		league:   l,
	}
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) GetWinCalls() []string {
	return s.winCalls
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func GetLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func AssertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func AssertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertStatus(t testing.TB, got *httptest.ResponseRecorder, want int) {
	t.Helper()
	responseCode := got.Code
	if responseCode != want {
		t.Errorf("did not get correct status, got %d, want %d", responseCode, want)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func AssertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got ScheduledAlert %v want %v", got, want)
	}
}
