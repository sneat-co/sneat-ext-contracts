package facade4contactus

import (
	"time"

	"github.com/dal-go/record"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// ContactusSpaceContributor builds the contactus records that must be persisted
// when a new space is created. It is implemented and registered by the contactus
// module so that spaceus does not depend on contactus DAL types directly.
type ContactusSpaceContributor interface {
	// BuildSpaceCreationRecords returns the contactus records (contactus space + creator member)
	// to insert as part of creating a new space.
	BuildSpaceCreationRecords(
		spaceID coretypes.SpaceID,
		userContactID string,
		creatorBrief briefs4contactus.ContactBrief,
		createdAt time.Time,
		byUserID string,
	) (records []record.Record, err error)
}
