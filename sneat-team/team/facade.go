package team

import "context"

// Follower is the facade interface for the follow operation.
// The implementation lives in the private sneat-team repo; this definition-only
// interface is the public contract consumed by callers.
type Follower interface {
	// Follow subscribes the given user to the specified team's updates.
	// It is idempotent: calling it when the user is already following returns
	// FollowTeamResponse.Followed == false.
	Follow(ctx context.Context, userID, teamID string) (FollowTeamResponse, error)
}
