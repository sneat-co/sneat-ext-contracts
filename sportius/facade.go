package sportius

import "context"

// Facade is the stable application boundary used by Telegram and other
// Sportius surfaces. actorUserID comes from authenticated host context and is
// never accepted from a wire request body.
type Facade interface {
	GetHome(ctx context.Context, actorUserID string) (SportsHome, error)
	GetPersonalProfile(ctx context.Context, actorUserID string) (PersonalSportsProfile, error)
	PutPersonalSport(ctx context.Context, actorUserID string, sportID SportID, request PutPersonalSportRequest) (PersonalSportsProfile, error)
	DeletePersonalSport(ctx context.Context, actorUserID string, sportID SportID) (PersonalSportsProfile, error)

	SearchTeams(ctx context.Context, actorUserID string, request SearchRequest) ([]TeamBrief, error)
	CreateTeam(ctx context.Context, actorUserID string, request CreateTeamRequest) (TeamView, error)
	GetTeam(ctx context.Context, actorUserID, spaceID string) (TeamView, error)
	UpdateTeam(ctx context.Context, actorUserID, spaceID string, request UpdateTeamRequest) (TeamView, error)
	JoinTeam(ctx context.Context, actorUserID, spaceID string, request JoinTeamRequest) (JoinTeamResponse, error)
	AddTeamPlayer(ctx context.Context, actorUserID, spaceID string, request AddPlayerRequest) (PlayerView, error)
	AddTeamStaff(ctx context.Context, actorUserID, spaceID string, request AddStaffRequest) (Participant, error)
	GetTeamPlayer(ctx context.Context, actorUserID, spaceID, playerContactID string) (PlayerView, error)
	// ListTeamGuardians returns reusable parent and guardian contacts to a team
	// manager. A returned contact is not necessarily a space member and this
	// read must not grant access.
	ListTeamGuardians(ctx context.Context, actorUserID, spaceID string) ([]ContactBrief, error)
	LinkGuardian(ctx context.Context, actorUserID, spaceID, playerContactID string, request LinkGuardianRequest) (PlayerView, error)

	SearchClubs(ctx context.Context, actorUserID string, request SearchRequest) ([]ClubBrief, error)
	CreateClub(ctx context.Context, actorUserID string, request CreateClubRequest) (ClubView, error)
	GetClub(ctx context.Context, actorUserID, spaceID string) (ClubView, error)
	UpdateClub(ctx context.Context, actorUserID, spaceID string, request UpdateClubRequest) (ClubView, error)
	AddClubStaff(ctx context.Context, actorUserID, spaceID string, request AddStaffRequest) (Participant, error)
	SetParticipantRoles(ctx context.Context, actorUserID string, kind SpaceKind, spaceID, contactID string, request SetParticipantRolesRequest) (Participant, error)
	LinkTeamToClub(ctx context.Context, actorUserID string, request LinkTeamToClubRequest) (ClubView, error)

	CreateInvitation(ctx context.Context, actorUserID string, request CreateInvitationRequest) (Invitation, error)
	GetInvitation(ctx context.Context, actorUserID, invitationID, claimToken string) (InvitationView, error)
	AcceptInvitation(ctx context.Context, actorUserID, invitationID string, request AcceptInvitationRequest) (InvitationAcceptance, error)
}
