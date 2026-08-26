package sportius

import "context"

// TwoPlayerRosterSchemaVersion is the version identifier for the stable
// two-player competition roster snapshot schema.
//
// A competition persists both this value and TwoPlayerRosterSnapshot.Version.
// The schema version identifies the snapshot shape; Version identifies the
// exact accepted team and player identities at registration time.
const TwoPlayerRosterSchemaVersion = "sportius.team-roster.v1"

// TeamRosterAuthority is the privileged server-to-server capability used by a
// competition host to resolve a registered team. It is intentionally separate
// from the user-facing Facade: callers receive an authoritative identity
// snapshot, not a general team-membership query.
//
// Implementations MUST reject a request when the team does not currently have
// exactly two eligible authenticated players. When ExpectedVersion is present,
// they MUST reject a roster whose current snapshot version differs, rather
// than silently replacing the accepted players.
type TeamRosterAuthority interface {
	ResolveTwoPlayerRoster(ctx context.Context, request TwoPlayerRosterRequest) (TwoPlayerRosterSnapshot, error)
}

// TwoPlayerRosterRequest identifies a team Space. ExpectedVersion is empty
// when the competition first accepts the team; thereafter it is the Version
// from the previously accepted snapshot and makes roster changes fail closed.
type TwoPlayerRosterRequest struct {
	TeamSpaceID     string `json:"teamSpaceID"`
	ExpectedVersion string `json:"expectedVersion,omitempty"`
}

// TwoPlayerRosterSnapshot is the deterministic, immutable-at-return-time
// identity record a competition persists for a two-player entry. It deliberately
// contains no mutable display or profile data.
type TwoPlayerRosterSnapshot struct {
	SchemaVersion string                  `json:"schemaVersion"`
	TeamSpaceID   string                  `json:"teamSpaceID"`
	Players       []TwoPlayerRosterMember `json:"players"`
	Version       string                  `json:"version"`
}

// TwoPlayerRosterMember is one authenticated player accepted in a competition
// roster. UserID is the entrant identity; ContactID carries the stable Sportius
// contact association without exposing profile data.
type TwoPlayerRosterMember struct {
	UserID    string `json:"userID"`
	ContactID string `json:"contactID"`
}
