package briefs4contactus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// GetFullContactID returns full member ContactID
func GetFullContactID(spaceID coretypes.SpaceID, memberID string) string {
	if spaceID == "" {
		panic("spaceID is required parameter")
	}
	if memberID == "" {
		panic("memberID is required parameter")
	}
	return string(spaceID) + ":" + memberID
}

// IsUniqueShortTitle checks if a given value is an unique member title
func IsUniqueShortTitle(v string, contacts map[string]*ContactBrief, role string) bool {
	for _, c := range contacts {
		if c.ShortTitle == v && (role == "" || c.HasRole(role)) {
			return false
		}
	}
	return true
}
