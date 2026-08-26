package dto4kidsclub

import (
	"encoding/json"
	"testing"
)

// TestRegistrationRequestJSONRoundTrip freezes the wire field names.
func TestRegistrationRequestJSONRoundTrip(t *testing.T) {
	in := `{"children":[{"name":"Ann","birthYear":2019,"gender":"f","allergies":"nuts"}],` +
		`"emergencyContactName":"Bob","emergencyContactPhone":"+353551234","pickups":["Bob"],` +
		`"acceptedConsentVersion":"v1"}`
	var r RegistrationRequest
	if err := json.Unmarshal([]byte(in), &r); err != nil {
		t.Fatal(err)
	}
	if len(r.Children) != 1 || r.Children[0].BirthYear != 2019 || r.AcceptedConsentVersion != "v1" {
		t.Fatalf("unexpected decode: %+v", r)
	}
}
