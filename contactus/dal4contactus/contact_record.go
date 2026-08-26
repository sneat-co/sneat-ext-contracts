// Package dal4contactus provides Contactus storage keys and record envelopes
// that other extensions may use as cross-boundary persistence glue.
package dal4contactus

import (
	"fmt"

	"github.com/dal-go/record"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/const4contactus"
	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// NewContactKey returns a Contactus contact key. The schema remains owned by
// Contactus; consumers receive no Contactus DTO through this helper.
func NewContactKey(spaceID coretypes.SpaceID, contactID string) *record.Key {
	if !core.IsAlphanumericOrUnderscore(contactID) {
		panic(fmt.Errorf("contactID should be alphanumeric, got: [%s]", contactID))
	}
	spaceModuleKey := coretypes.NewSpaceModuleKey(spaceID, const4contactus.ExtensionID)
	return record.NewKeyWithParentAndID(spaceModuleKey, const4contactus.ContactsCollection, contactID)
}

// NewContactRecord returns an envelope used only to check a contact's record
// existence. Its map payload intentionally avoids exposing a Contactus DTO.
func NewContactRecord(spaceID coretypes.SpaceID, contactID string) record.Record {
	return record.NewRecordWithData(NewContactKey(spaceID, contactID), new(map[string]any))
}
