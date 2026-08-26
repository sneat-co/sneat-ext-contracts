package sportius

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestErrorKeepsCausePrivate(t *testing.T) {
	cause := errors.New("private backend detail")
	err := &Error{
		Code:       ErrorCodeRetryable,
		MessageKey: "sportius.error.try_again",
		Retryable:  true,
		Cause:      cause,
	}

	if !errors.Is(err, cause) {
		t.Fatal("error does not unwrap its internal cause")
	}
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}
	if string(data) != `{"code":"retryable","messageKey":"sportius.error.try_again","retryable":true}` {
		t.Fatalf("unexpected public error JSON: %s", data)
	}
	if err.Error() == cause.Error() {
		t.Fatal("public error string exposed its private cause")
	}
}

func TestUpdateRequestsCanExplicitlyClearOptionalFields(t *testing.T) {
	team := UpdateTeamRequest{RequestID: "team-update", ClearAge: true, ClearLocation: true, ClearMedia: true}
	club := UpdateClubRequest{
		RequestID:         "club-update",
		ClearPrimarySport: true,
		ReplaceSportIDs:   true,
		SportIDs:          []SportID{},
		ClearLocation:     true,
		ClearMedia:        true,
	}

	teamJSON, err := json.Marshal(team)
	if err != nil {
		t.Fatalf("marshal team update: %v", err)
	}
	clubJSON, err := json.Marshal(club)
	if err != nil {
		t.Fatalf("marshal club update: %v", err)
	}
	if string(teamJSON) != `{"requestID":"team-update","clearAge":true,"clearLocation":true,"clearMedia":true}` {
		t.Fatalf("unexpected team clear patch: %s", teamJSON)
	}
	if string(clubJSON) != `{"requestID":"club-update","clearPrimarySport":true,"replaceSportIDs":true,"clearLocation":true,"clearMedia":true}` {
		t.Fatalf("unexpected club clear patch: %s", clubJSON)
	}
}
