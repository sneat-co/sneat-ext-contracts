package briefs4contactus

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/validation"
)

// ContactBase is used in dbo4contactus.ContactDbo and in requests to create a contactBrief
type ContactBase struct {
	ContactBrief

	Status dbmodels.Status `json:"status" firestore:"status"` // active, archived

	dbmodels.WithUpdatedAndVersion

	WithGroupIDs

	Address   *dbmodels.Address `json:"address,omitempty" firestore:"address,omitempty"`
	VATNumber string            `json:"vatNumber,omitempty" firestore:"vatNumber,omitempty"`

	//EmailsObsolete []dbmodels.PersonEmail `json:"emails,omitempty" firestore:"emails,omitempty"`
	//PhonesObsolete []dbmodels.PersonPhone `json:"phones,omitempty" firestore:"phones,omitempty"`
	Avatars []dbprofile.Avatar `json:"avatars,omitempty" firestore:"avatars,omitempty"`

	dbmodels.WithTimezone
}

func (v *ContactBase) Equal(v2 *ContactBase) bool {
	return v.ContactBrief.Equal(&v2.ContactBrief) && v.VATNumber == v2.VATNumber
}

// Validate returns error if not valid
func (v *ContactBase) Validate() error {
	var errs []error
	if err := v.ContactBrief.Validate(); err != nil {
		errs = append(errs, err)
	}
	switch v.Status {
	case ContactStatusActive, ContactStatusArchived, ContactStatusDeleted: // OK
	case "":
		errs = append(errs, validation.NewErrRequestIsMissingRequiredField("status"))
	default:
		errs = append(errs, validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status))
	}
	if v.Type == ContactTypeCompany {
		if v.CountryID == "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("countryID", "missing required field for a contactBrief of type=company"))
		}
	}
	if strings.TrimSpace(v.Title) == "" && v.Names == nil {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("name|title"))
	}
	if v.Names != nil {
		if err := v.Names.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if v.DoB != "" {
		if t, err := time.Parse(time.DateOnly, v.DoB); err != nil {
			return validation.NewErrBadRecordFieldValue("dob", "invalid date of birth: "+v.DoB)
		} else if t.After(time.Now()) {
			return validation.NewErrBadRecordFieldValue("dob", "date of birth cannot be in the future: "+v.DoB)
		}
	}
	switch v.Type {
	case ContactTypePerson, ContactTypeAnimal:
		if v.VATNumber != "" {
			return validation.NewErrBadRecordFieldValue("vatNumber", "should be empty for a contactBrief of type=person")
		}
		if err := dbmodels.ValidateGender(v.Gender, true); err != nil {
			errs = append(errs, err)
		}
		if v.AgeGroup != "" || v.Type != ContactTypeAnimal {
			if err := dbmodels.ValidateAgeGroup(v.AgeGroup, true); err != nil {
				errs = append(errs, err)
			}
		}
	case ContactTypeCompany:
		if v.Gender != "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("gender", "expected to be empty for contactBrief type=company, got: "+v.Gender))
		}
		if err := dbmodels.ValidateGender(v.Gender, false); err != nil {
			return err
		}
		if v.AgeGroup != "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("ageGroup", "expected to be empty for contactBrief type=company, got: "+v.Gender))
		}
		if err := dbmodels.ValidateAgeGroup(v.AgeGroup, false); err != nil {
			errs = append(errs, err)
		}
	}
	for i, avatar := range v.Avatars {
		if err := avatar.Validate(); err != nil {
			errs = append(errs, validation.NewErrBadRecordFieldValue(fmt.Sprintf("avatars[%d]", i), err.Error()))
		}
	}
	if err := v.WithGroupIDs.Validate(); err != nil {
		errs = append(errs, validation.NewErrBadRecordFieldValue("withGroupIDs", err.Error()))
	}

	if err := v.Timezone.Validate(); err != nil {
		errs = append(errs, validation.NewErrBadRecordFieldValue("timezone", err.Error()))
	}

	if l := len(errs); l == 1 {
		return validation.NewErrBadRecordFieldValue("ContactBase", errs[0].Error())
	} else if l > 0 {
		return validation.NewErrBadRecordFieldValue("ContactBase", fmt.Errorf("%d errors:\n%w", l, errors.Join(errs...)).Error())
	}

	return nil
}
