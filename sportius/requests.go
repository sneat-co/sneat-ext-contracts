package sportius

type PutPersonalSportRequest struct {
	RoleIDs    []RoleID          `json:"roleIDs"`
	Visibility ProfileVisibility `json:"visibility"`
}

type SearchRequest struct {
	Name     string  `json:"name"`
	SportID  SportID `json:"sportID,omitempty"`
	Locality string  `json:"locality,omitempty"`
}

type CreateTeamRequest struct {
	RequestID      string         `json:"requestID"`
	Name           string         `json:"name"`
	SportID        SportID        `json:"sportID"`
	CreatorRoleIDs []RoleID       `json:"creatorRoleIDs"`
	Gender         GenderCategory `json:"gender,omitempty"`
	Age            *AgeRange      `json:"age,omitempty"`
	Location       *LocationHint  `json:"location,omitempty"`
	Media          *MediaRef      `json:"media,omitempty"`
	JoinPolicy     JoinPolicy     `json:"joinPolicy,omitempty"`
}

type UpdateTeamRequest struct {
	RequestID     string          `json:"requestID"`
	Name          *string         `json:"name,omitempty"`
	SportID       *SportID        `json:"sportID,omitempty"`
	Gender        *GenderCategory `json:"gender,omitempty"`
	Age           *AgeRange       `json:"age,omitempty"`
	Location      *LocationHint   `json:"location,omitempty"`
	Media         *MediaRef       `json:"media,omitempty"`
	JoinPolicy    *JoinPolicy     `json:"joinPolicy,omitempty"`
	ClearAge      bool            `json:"clearAge,omitempty"`
	ClearLocation bool            `json:"clearLocation,omitempty"`
	ClearMedia    bool            `json:"clearMedia,omitempty"`
}

type JoinTeamRequest struct {
	RequestID    string   `json:"requestID"`
	RoleIDs      []RoleID `json:"roleIDs"`
	InvitationID string   `json:"invitationID,omitempty"`
	// ClaimToken is required whenever InvitationID is supplied.
	ClaimToken string `json:"claimToken,omitempty"`
}

type JoinTeamResponse struct {
	Team                TeamBrief  `json:"team"`
	Status              JoinStatus `json:"status"`
	RoleIDs             []RoleID   `json:"roleIDs"`
	MembershipRequestID string     `json:"membershipRequestID,omitempty"`
}

type AddPlayerRequest struct {
	RequestID   string   `json:"requestID"`
	DisplayName string   `json:"displayName"`
	RoleIDs     []RoleID `json:"roleIDs"`
	UserID      string   `json:"userID,omitempty"`
}

type AddStaffRequest struct {
	RequestID   string   `json:"requestID"`
	DisplayName string   `json:"displayName"`
	RoleIDs     []RoleID `json:"roleIDs"`
	UserID      string   `json:"userID,omitempty"`
}

type SetParticipantRolesRequest struct {
	RequestID string   `json:"requestID"`
	RoleIDs   []RoleID `json:"roleIDs"`
}

type LinkGuardianRequest struct {
	RequestID           string `json:"requestID"`
	GuardianContactID   string `json:"guardianContactID,omitempty"`
	GuardianDisplayName string `json:"guardianDisplayName,omitempty"`
	RelationshipRoleID  string `json:"relationshipRoleID"`
}

type CreateClubRequest struct {
	RequestID      string        `json:"requestID"`
	Name           string        `json:"name"`
	PrimarySportID SportID       `json:"primarySportID,omitempty"`
	SportIDs       []SportID     `json:"sportIDs"`
	CreatorRoleIDs []RoleID      `json:"creatorRoleIDs"`
	Location       *LocationHint `json:"location,omitempty"`
	Media          *MediaRef     `json:"media,omitempty"`
}

type UpdateClubRequest struct {
	RequestID         string        `json:"requestID"`
	Name              *string       `json:"name,omitempty"`
	PrimarySportID    *SportID      `json:"primarySportID,omitempty"`
	SportIDs          []SportID     `json:"sportIDs,omitempty"`
	Location          *LocationHint `json:"location,omitempty"`
	Media             *MediaRef     `json:"media,omitempty"`
	ClearPrimarySport bool          `json:"clearPrimarySport,omitempty"`
	ReplaceSportIDs   bool          `json:"replaceSportIDs,omitempty"`
	ClearLocation     bool          `json:"clearLocation,omitempty"`
	ClearMedia        bool          `json:"clearMedia,omitempty"`
}

type LinkTeamToClubRequest struct {
	RequestID   string `json:"requestID"`
	TeamSpaceID string `json:"teamSpaceID"`
	ClubSpaceID string `json:"clubSpaceID"`
}

type CreateInvitationRequest struct {
	RequestID string    `json:"requestID"`
	SpaceID   string    `json:"spaceID"`
	Kind      SpaceKind `json:"kind"`
	// ContactID targets an existing contact in the space. Exactly one of
	// ContactID and InviteeDisplayName is required.
	ContactID string `json:"contactID,omitempty"`
	// InviteeDisplayName asks the implementation to create a non-member
	// contact before issuing the invitation.
	InviteeDisplayName string   `json:"inviteeDisplayName,omitempty"`
	SuggestedRoleIDs   []RoleID `json:"suggestedRoleIDs"`
}

type AcceptInvitationRequest struct {
	RequestID  string   `json:"requestID"`
	ClaimToken string   `json:"claimToken"`
	RoleIDs    []RoleID `json:"roleIDs"`
}
