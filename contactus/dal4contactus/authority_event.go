package dal4contactus

import (
	"github.com/dal-go/record"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contract4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// NewAuthorityEventRecord uses the native Contactus extension root. Records
// must be inserted immutably in the contact mutation's original transaction.
func NewAuthorityEventRecord(event contract4contactus.ContactAuthorityEvent) (record.Record, error) {
	if err := event.Validate(); err != nil {
		return nil, err
	}
	parent := coretypes.NewSpaceModuleKey(coretypes.SpaceID(event.SpaceID), const4contactus.ExtensionID)
	key := record.NewKeyWithParentAndID(parent, contract4contactus.AuthorityEventsCollection, event.ID())
	return record.NewRecordWithData(key, &event), nil
}
