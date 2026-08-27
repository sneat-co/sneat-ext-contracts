package contract4chessraiders_test

import (
	"testing"

	"github.com/sneat-co/sneat-ext-contracts/chessraiders/contract4chessraiders"
)

func TestExtensionID(t *testing.T) {
	if contract4chessraiders.ExtensionID != "chess-raiders" {
		t.Fatalf("ExtensionID = %q, want %q", contract4chessraiders.ExtensionID, "chess-raiders")
	}
}

func TestSideValues(t *testing.T) {
	cases := map[contract4chessraiders.Side]string{
		contract4chessraiders.SideNone:  "",
		contract4chessraiders.SideWhite: "w",
		contract4chessraiders.SideBlack: "b",
	}
	for side, want := range cases {
		if string(side) != want {
			t.Errorf("Side %v = %q, want %q", side, string(side), want)
		}
	}
}

func TestLobbyViewRoundTrips(t *testing.T) {
	view := contract4chessraiders.LobbyView{
		MatchID:     "match-1",
		Lifecycle:   contract4chessraiders.LifecycleLobby,
		CaptureMode: contract4chessraiders.CaptureModeKillOnly,
		CreatedBy:   "user-1",
		White: []contract4chessraiders.LobbySeat{
			{UserID: "user-1", DisplayName: "Creator", Ready: true},
		},
	}
	if len(view.White) != 1 || view.White[0].UserID != "user-1" {
		t.Fatalf("LobbyView.White = %+v, want one seat for user-1", view.White)
	}
	if view.Lifecycle != contract4chessraiders.LifecycleLobby {
		t.Fatalf("LobbyView.Lifecycle = %v, want LifecycleLobby", view.Lifecycle)
	}
}
