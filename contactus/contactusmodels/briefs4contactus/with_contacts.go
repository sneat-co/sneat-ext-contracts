package briefs4contactus

import (
	"fmt"
	"slices"

	"github.com/dal-go/record/update"
	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

type contactBrief interface {
	dbmodels.UserIDGetter
	dbmodels.RelatedAs
}

type WithSingleSpaceContactsWithoutContactIDs[
	T interface {
		contactBrief
		HasRole(role string) bool
		Equal(v T) bool
	},
] struct {
	WithContactsBase[T]
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) Validate() error {
	for id, brief := range v.Contacts {
		if err := validate.RecordID(id); err != nil {
			return validation.NewErrBadRecordFieldValue(const4contactus.ContactsField,
				fmt.Sprintf("invalid contact ContactID=%s: %v", id, err))
		}
		if err := brief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("contacts."+id, err.Error())
		}
	}
	return nil
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) ContactIDs() (contactIDs []string) {
	contactIDs = make([]string, 0, len(v.Contacts))
	for id := range v.Contacts {
		contactIDs = append(contactIDs, id)
	}
	return
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) HasContact(contactID string) bool {
	_, ok := v.Contacts[contactID]
	return ok
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) AddContact(contactID string, contact T) update.Update {
	if v.Contacts == nil {
		v.Contacts = make(map[string]T)
	}
	v.Contacts[contactID] = contact
	return update.ByFieldPath([]string{"contacts", contactID}, contact)
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) RemoveContact(contactID string) update.Update {
	delete(v.Contacts, contactID)
	return update.ByFieldPath([]string{"contacts", contactID}, update.DeleteField)
}

// WithMultiSpaceContacts mixin that adds WithMultiSpaceContactIDs.ContactIDs & Contacts fields
type WithMultiSpaceContacts[
	T interface {
		contactBrief
		HasRole(role string) bool
		Equal(v T) bool
	},
] struct {
	WithMultiSpaceContactIDs
	WithContactsBase[T]
}

// Validate returns error if not valid.
//
// NB the contactIDs check below must RETURN its error. It used to `return nil`
// on failure, silently discarding the verdict of a real format validator --
// and since this mixin is embedded in contactus's ContactDbo, that meant
// contactIDs were never actually validated on any contact record.
func (v *WithMultiSpaceContacts[T]) Validate() error {
	if err := v.WithMultiSpaceContactIDs.Validate(); err != nil {
		return err
	}
	return dbmodels.ValidateWithIdsAndBriefs("contactIDs", const4contactus.ContactsField, v.ContactIDs, v.Contacts)
}

func (v *WithMultiSpaceContacts[T]) Updates(contactIDs ...dbmodels.SpaceItemID) (updates []update.Update) {
	updates = append(updates, update.ByFieldName("contactIDs", v.ContactIDs))
	if len(contactIDs) == 0 {
		updates = append(updates, update.ByFieldName(const4contactus.ContactsField, v.Contacts))
	} else {
		for _, id := range contactIDs {
			updates = append(updates, update.ByFieldName(const4contactus.ContactsField+"."+string(id), v.Contacts[string(id)]))
		}
	}
	return
}

// SetContactBrief sets contactBrief brief by ContactID
func (v *WithMultiSpaceContacts[T]) SetContactBrief(spaceID coretypes.SpaceID, contactID string, contactBrief T) (updates []update.Update) {
	id := string(dbmodels.NewSpaceItemID(spaceID, contactID))
	if !slices.Contains(v.ContactIDs, id) {
		v.ContactIDs = append(v.ContactIDs, id)
		updates = append(updates, update.ByFieldName("contactIDs", v.ContactIDs))
	}
	if v.Contacts == nil {
		v.Contacts = make(map[string]T)
	}
	if currentBrief, ok := v.Contacts[id]; !ok || !currentBrief.Equal(contactBrief) {
		v.Contacts[id] = contactBrief
		updates = append(updates, update.ByFieldName(const4contactus.ContactsField+"."+id, contactBrief))
	}
	return
}

// ParentContactBrief returns parent contactBrief brief
func (v *WithMultiSpaceContacts[T]) ParentContactBrief() (i int, id dbmodels.SpaceItemID, brief T) {
	for i, id := range v.ContactIDs {
		brief := v.Contacts[id]
		if brief.GetRelatedAs() == "parent" {
			return i, dbmodels.SpaceItemID(id), brief
		}
	}
	return -1, "", brief
}

// GetContactBriefByID returns contactBrief brief by ContactID
func (v *WithMultiSpaceContacts[T]) GetContactBriefByID(spaceID coretypes.SpaceID, contactID string) (i int, brief T) {
	id := dbmodels.NewSpaceItemID(spaceID, contactID)
	if brief, ok := v.Contacts[string(id)]; !ok {
		return -1, brief
	}
	return slice.Index(v.ContactIDs, string(id)), brief
}

// GetContactBriefByUserID returns contactBrief brief by user ContactID
func (v *WithMultiSpaceContacts[T]) GetContactBriefByUserID(userID string) (id dbmodels.SpaceItemID, t T) {
	for cID, c := range v.Contacts {
		if c.GetUserID() == userID {
			return dbmodels.SpaceItemID(cID), c
		}
	}
	return
}

func (v *WithMultiSpaceContacts[T]) AddContact(spaceID coretypes.SpaceID, contactID string, c T) (updates []update.Update) {
	id := dbmodels.NewSpaceItemID(spaceID, contactID)
	if !slices.Contains(v.ContactIDs, string(id)) {
		if len(v.ContactIDs) == 0 {
			v.ContactIDs = make([]string, 1, 2)
			v.ContactIDs[0] = "*"
		}
		v.ContactIDs = append(v.ContactIDs, string(id))
		updates = append(updates, update.ByFieldName("contactIDs", v.ContactIDs))
	}
	if _, ok := v.Contacts[string(id)]; !ok {
		updates = append(updates, update.ByFieldName(const4contactus.ContactsField+"."+string(id), c))
	}
	if v.Contacts == nil {
		v.Contacts = make(map[string]T)
	}
	v.Contacts[string(id)] = c
	return
}

func (v *WithMultiSpaceContacts[T]) RemoveContact(spaceID coretypes.SpaceID, contactID string) (updates []update.Update) {
	id := dbmodels.NewSpaceItemID(spaceID, contactID)
	contactIDs := slice.RemoveInPlaceByValue(v.ContactIDs, string(id))
	if len(contactIDs) != len(v.ContactIDs) {
		v.ContactIDs = contactIDs
		updates = append(updates, update.ByFieldName("contactIDs", v.ContactIDs))
	}
	if _, ok := v.Contacts[string(id)]; ok {
		delete(v.Contacts, string(id))
		updates = append(updates, update.ByFieldName(const4contactus.ContactsField+"."+string(id), update.DeleteField))
	}
	return
}
