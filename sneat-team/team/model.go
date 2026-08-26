// Package team holds the public contract types for the team extension follow operation.
// Types here are hand-implemented to match ../../../typespec/api4team.tsp (no emitters).
package team

// FollowTeamRequest is the request body for POST /v0/api4team/follow.
// JSON field names match the TypeSpec model (camelCase).
type FollowTeamRequest struct {
	TeamID string `json:"teamID"`
}

// FollowTeamResponse is the response body returned after a follow operation.
// JSON field names match the TypeSpec model (camelCase).
type FollowTeamResponse struct {
	TeamID   string `json:"teamID"`
	Followed bool   `json:"followed"`
}
