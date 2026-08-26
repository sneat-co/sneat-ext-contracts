package sportius

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInvitationCarriesStableContactIdentity(t *testing.T) {
	value := Invitation{
		InvitationID:       "invite-1",
		SpaceID:            "team-1",
		Kind:               SpaceKindTeam,
		ContactID:          "contact-1",
		InviteeDisplayName: "Alex",
		SuggestedRoleIDs:   []RoleID{RolePlayer},
		DeepLink:           "https://t.me/sneat_bot?start=invite-1",
	}

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	jsonValue := string(data)
	for _, fragment := range []string{
		`"contactID":"contact-1"`,
		`"inviteeDisplayName":"Alex"`,
		`"suggestedRoleIDs":["player"]`,
	} {
		if !strings.Contains(jsonValue, fragment) {
			t.Fatalf("invitation JSON %s does not contain %s", jsonValue, fragment)
		}
	}
}

func TestInvitationAcceptanceReturnsClaimedContact(t *testing.T) {
	value := InvitationAcceptance{
		InvitationID: "invite-1",
		SpaceID:      "team-1",
		Kind:         SpaceKindTeam,
		ContactID:    "contact-1",
	}

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"contactID":"contact-1"`) {
		t.Fatalf("acceptance JSON does not identify the claimed contact: %s", data)
	}
}

func TestInvitationAcceptanceCarriesClaimProofOnlyInRequest(t *testing.T) {
	request := AcceptInvitationRequest{
		RequestID:  "accept-1",
		ClaimToken: "opaque-proof",
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"claimToken":"opaque-proof"`) {
		t.Fatalf("acceptance request JSON has no claim proof: %s", data)
	}

	invitationData, err := json.Marshal(Invitation{InvitationID: "invite-1"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(invitationData), "claimToken") {
		t.Fatalf("invitation response must not expose claim proof: %s", invitationData)
	}
	if strings.Contains(string(invitationData), "deepLink") {
		t.Fatalf("invitation inspection omits creation-only deep link: %s", invitationData)
	}
}

func TestInvitationJoinCarriesClaimProofWithInvitationID(t *testing.T) {
	request := JoinTeamRequest{
		RequestID:    "join-1",
		InvitationID: "invite-1",
		ClaimToken:   "opaque-proof",
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	jsonValue := string(data)
	for _, fragment := range []string{
		`"invitationID":"invite-1"`,
		`"claimToken":"opaque-proof"`,
	} {
		if !strings.Contains(jsonValue, fragment) {
			t.Fatalf("invitation join JSON %s does not contain %s", jsonValue, fragment)
		}
	}
}
