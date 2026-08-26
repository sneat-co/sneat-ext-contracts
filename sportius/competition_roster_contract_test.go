package sportius

import (
	"context"
	"encoding/json"
	"testing"
)

type rosterAuthorityCompileGuard struct{}

func (rosterAuthorityCompileGuard) ResolveTwoPlayerRoster(_ context.Context, _ TwoPlayerRosterRequest) (TwoPlayerRosterSnapshot, error) {
	return TwoPlayerRosterSnapshot{}, nil
}

var _ TeamRosterAuthority = rosterAuthorityCompileGuard{}

func TestTwoPlayerRosterContractJSON(t *testing.T) {
	request, err := json.Marshal(TwoPlayerRosterRequest{
		TeamSpaceID:     "team-space-1",
		ExpectedVersion: "accepted-version",
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	if got, want := string(request), `{"teamSpaceID":"team-space-1","expectedVersion":"accepted-version"}`; got != want {
		t.Fatalf("request JSON = %s, want %s", got, want)
	}

	snapshot, err := json.Marshal(TwoPlayerRosterSnapshot{
		SchemaVersion: TwoPlayerRosterSchemaVersion,
		TeamSpaceID:   "team-space-1",
		Players: []TwoPlayerRosterMember{
			{UserID: "user-1", ContactID: "contact-1"},
			{UserID: "user-2", ContactID: "contact-2"},
		},
		Version: "accepted-version",
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if got, want := string(snapshot), `{"schemaVersion":"sportius.team-roster.v1","teamSpaceID":"team-space-1","players":[{"userID":"user-1","contactID":"contact-1"},{"userID":"user-2","contactID":"contact-2"}],"version":"accepted-version"}`; got != want {
		t.Fatalf("snapshot JSON = %s, want %s", got, want)
	}
}

func TestTwoPlayerRosterRequestOmitsEmptyExpectedVersion(t *testing.T) {
	data, err := json.Marshal(TwoPlayerRosterRequest{TeamSpaceID: "team-space-1"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	if got, want := string(data), `{"teamSpaceID":"team-space-1"}`; got != want {
		t.Fatalf("request JSON = %s, want %s", got, want)
	}
}
