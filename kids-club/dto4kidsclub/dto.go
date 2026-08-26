// Package dto4kidsclub holds the frozen wire DTOs of the kidsclub extension
// (contract surface shared between the Go backend and other repos). The
// TypeSpec source of truth is ../../typespec/api4kidsclub.tsp.
package dto4kidsclub

import "time"

// CreateClubRequest creates a club in the operator's business space.
// POST /v0/api4kidsclub/spaces/{spaceID}/clubs
type CreateClubRequest struct {
	Title string `json:"title"`
	// Venue is a free-text venue type, e.g. "hotel", "jump zone".
	Venue string `json:"venue,omitempty"`
	// ConsentText is the club's v1 waiver text shown to parents at
	// registration. Required — registration records acceptance of it.
	ConsentText string `json:"consentText"`
}

// CreateClubResponse identifies the new club so the UI can navigate to it.
type CreateClubResponse struct {
	ID string `json:"id"`
	// PublicPath is the site-relative public club page path, e.g. "/club/abc".
	PublicPath string `json:"publicPath"`
}

// ClubPublic is the unauthenticated public projection of a club.
// GET /v0/api4kidsclub/clubs/{clubID}
type ClubPublic struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Venue          string `json:"venue,omitempty"`
	ConsentVersion string `json:"consentVersion"`
	ConsentText    string `json:"consentText"`
}

// Child is one child in a registration.
type Child struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birthYear"`
	Gender    string `json:"gender,omitempty"` // "f" | "m" | "" (unspecified)
	Allergies string `json:"allergies,omitempty"`
}

// ConsentAcceptance records acceptance of the club waiver.
type ConsentAcceptance struct {
	Version    string    `json:"version"`
	AcceptedAt time.Time `json:"acceptedAt"`
	ByUID      string    `json:"byUID"`
}

// RegistrationRequest submits children to a club.
// POST /v0/api4kidsclub/clubs/{clubID}/registrations
// AcceptedConsentVersion must equal the club's current ConsentVersion; the
// server stamps AcceptedAt and ByUID.
type RegistrationRequest struct {
	Children               []Child  `json:"children"`
	EmergencyContactName   string   `json:"emergencyContactName"`
	EmergencyContactPhone  string   `json:"emergencyContactPhone"`
	Pickups                []string `json:"pickups,omitempty"`
	AcceptedConsentVersion string   `json:"acceptedConsentVersion"`
}

// RegistrationResponse identifies the new registration.
type RegistrationResponse struct {
	ID string `json:"id"`
}

// RosterEntry is one registration row in the operator's roster.
// GET /v0/api4kidsclub/spaces/{spaceID}/clubs/{clubID}/registrations
type RosterEntry struct {
	ID                    string            `json:"id"`
	Children              []Child           `json:"children"`
	EmergencyContactName  string            `json:"emergencyContactName"`
	EmergencyContactPhone string            `json:"emergencyContactPhone"`
	Pickups               []string          `json:"pickups,omitempty"`
	Consent               ConsentAcceptance `json:"consent"`
	CreatedAt             time.Time         `json:"createdAt"`
}
