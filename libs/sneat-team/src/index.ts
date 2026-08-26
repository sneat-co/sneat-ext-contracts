// @sneat/extension-team-contract — public TS contract for the team extension.
// Types are hand-implemented to match ../../typespec/api4team.tsp (no emitters).
// See also: github.com/sneat-co/ext-sneat-team/backend/team for the Go counterpart.

/** Request body for POST /v0/api4team/follow. */
export interface FollowTeamRequest {
	/** The ID of the team to follow. */
	teamID: string;
}

/** Response body returned after a follow operation. */
export interface FollowTeamResponse {
	/** The ID of the team that was followed. */
	teamID: string;

	/** True if the user is now following the team; false if they were already following. */
	followed: boolean;
}
