package dal4contactus

import (
	"testing"

	dalrecord "github.com/dal-go/record"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestNewContactRecordUsesContactusSchemaKey(t *testing.T) {
	t.Parallel()
	spaceID := coretypes.SpaceID("space_1")
	contactID := "contact_1"
	parent := coretypes.NewSpaceModuleKey(spaceID, const4contactus.ExtensionID)
	want := dalrecord.NewKeyWithParentAndID(parent, const4contactus.ContactsCollection, contactID)

	key := NewContactKey(spaceID, contactID)
	if key.String() != want.String() {
		t.Fatalf("NewContactKey() = %v, want %v", key, want)
	}
	if record := NewContactRecord(spaceID, contactID); record.Key().String() != key.String() {
		t.Fatalf("NewContactRecord() key = %v, want %v", record.Key(), key)
	}
}

func TestNewContactKeyRejectsNonSchemaID(t *testing.T) {
	t.Parallel()
	defer func() {
		if recover() == nil {
			t.Fatal("NewContactKey() did not reject an invalid ID")
		}
	}()
	NewContactKey("space_1", "not a valid id")
}
