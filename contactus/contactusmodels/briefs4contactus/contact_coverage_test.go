package briefs4contactus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-ext-contracts/contactus/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
)

func TestContactGroupBrief_InvalidTitle(t *testing.T) {
	b := &ContactGroupBrief{Title: " leading space"}
	if err := b.Validate(); err == nil {
		t.Fatal("expected error on title with leading space")
	}
}

func TestWithContactIDs_EdgeCases(t *testing.T) {
	single := &WithSingleSpaceContactIDs{}
	single.AddContactID("c1")
	single.AddContactID("c2")
	if len(single.ContactIDs) != 3 {
		t.Fatalf("expected 3 contact IDs, got %d", len(single.ContactIDs))
	}

	multi := &WithMultiSpaceContactIDs{} // empty ContactIDs -> fails WithContactIDs.Validate()
	if err := multi.Validate(); err == nil {
		t.Fatal("expected error on empty multi.Validate()")
	}
}

func TestContact_ValidationHelpers(t *testing.T) {
	// WithGroupIDs
	g := WithGroupIDs{GroupIDs: []string{"g1"}}
	if err := g.Validate(); err != nil {
		t.Fatalf("expected valid group IDs, got %v", err)
	}
	gInvalid := WithGroupIDs{GroupIDs: []string{" g1"}}
	if err := gInvalid.Validate(); err == nil {
		t.Fatal("expected error on untrimmed group ID")
	}

	// ValidateContactIDRecordField
	if err := ValidateContactIDRecordField("cid", "", false); err != nil {
		t.Fatalf("expected nil when optional and empty, got %v", err)
	}
	if err := ValidateContactIDRecordField("cid", "", true); err == nil {
		t.Fatal("expected error when required and empty")
	}
	if err := ValidateContactIDRecordField("cid", "valid_id", true); err != nil {
		t.Fatalf("expected nil on valid id, got %v", err)
	}
	if err := ValidateContactIDRecordField("cid", "invalid id", true); err == nil {
		t.Fatal("expected error on id with spaces")
	}

	// ValidateContactType
	if err := ValidateContactType(""); err == nil {
		t.Fatal("expected error on empty ContactType")
	}
	if err := ValidateContactType("unknown_type"); err == nil {
		t.Fatal("expected error on unknown ContactType")
	}
	for _, ct := range ContactTypes {
		if err := ValidateContactType(ct); err != nil {
			t.Fatalf("expected valid for %s, got %v", ct, err)
		}
	}
}

func TestContactBrief_SetNameAndRoles(t *testing.T) {
	b := &ContactBrief{}

	b.SetName("first", "John")
	if b.Names.FirstName != "John" {
		t.Fatalf("expected FirstName 'John', got %q", b.Names.FirstName)
	}
	b.SetName("last", "Doe")
	if b.Names.LastName != "Doe" {
		t.Fatalf("expected LastName 'Doe', got %q", b.Names.LastName)
	}
	b.SetName("middle", "M")
	if b.Names.MiddleName != "M" {
		t.Fatalf("expected MiddleName 'M', got %q", b.Names.MiddleName)
	}
	b.SetName("full", "John M Doe")
	if b.Names.FullName != "John M Doe" {
		t.Fatalf("expected FullName 'John M Doe', got %q", b.Names.FullName)
	}
	b.SetName("nick", "Johnny")
	if b.Names.NickName != "Johnny" {
		t.Fatalf("expected NickName 'Johnny', got %q", b.Names.NickName)
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on unsupported field")
			}
		}()
		b.SetName("unsupported", "val")
	}()

	b.RolesField = with.RolesField{Roles: []string{const4contactus.SpaceMemberRoleMember}}
	if !b.IsSpaceMember() {
		t.Fatal("expected IsSpaceMember to be true")
	}
	b.UserID = "u1"
	if b.GetUserID() != "u1" {
		t.Fatalf("expected GetUserID 'u1', got %q", b.GetUserID())
	}

	b2 := &ContactBrief{}
	*b2 = *b
	b2.Names = &person.NameFields{}
	*b2.Names = *b.Names
	if !b.Equal(b2) {
		t.Fatal("expected b.Equal(b2) to be true")
	}

	b2.Type = ContactTypeCompany
	if b.Equal(b2) {
		t.Fatal("expected b.Equal(b2) to be false on different type")
	}
}

func TestContactBrief_Validate(t *testing.T) {
	// Valid person brief
	b := &ContactBrief{
		Type:  ContactTypePerson,
		Title: "John Doe",
	}
	if err := b.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Invalid type
	bInvalidType := &ContactBrief{Type: "unknown", Title: "Test"}
	if err := bInvalidType.Validate(); err == nil {
		t.Fatal("expected error on unknown type")
	}

	// Missing title & names
	bNoName := &ContactBrief{Type: ContactTypePerson}
	if err := bNoName.Validate(); err == nil {
		t.Fatal("expected error on missing name|title")
	}

	// Invalid names
	bBadNames := &ContactBrief{Type: ContactTypePerson, Names: &person.NameFields{FirstName: " a "}}
	if err := bBadNames.Validate(); err == nil {
		t.Fatal("expected error on untrimmed name")
	}

	// Invalid UserID (spaces/symbols)
	bBadUID := &ContactBrief{Type: ContactTypePerson, Title: "John", WithUserID: dbmodels.WithUserID{UserID: "bad user id!"}}
	if err := bBadUID.Validate(); err == nil {
		t.Fatal("expected error on invalid UserID")
	}

	// Location requires ParentID
	bLoc := &ContactBrief{Type: ContactTypeLocation, Title: "Loc 1"}
	if err := bLoc.Validate(); err == nil {
		t.Fatal("expected error on location without ParentID")
	}
	bLoc.ParentID = "parent1"
	if err := bLoc.Validate(); err != nil {
		t.Fatalf("expected valid location with ParentID, got %v", err)
	}

	// PetKind
	bPet := &ContactBrief{Type: ContactTypeAnimal, Title: "Doggy", PetKind: "unknown_species"}
	if err := bPet.Validate(); err == nil {
		t.Fatal("expected error on unknown PetKind")
	}
	bPet.PetKind = const4contactus.PetKindDog
	if err := bPet.Validate(); err != nil {
		t.Fatalf("expected valid pet, got %v", err)
	}
}

func TestContactBrief_TitlesAndShortTitles(t *testing.T) {
	b := &ContactBrief{Title: "Title One"}
	if b.GetTitle() != "Title One" {
		t.Fatalf("expected 'Title One', got %q", b.GetTitle())
	}
	b.Title = ""
	b.Names = &person.NameFields{FullName: "Full Name"}
	if b.GetTitle() != "Full Name" {
		t.Fatalf("expected 'Full Name', got %q", b.GetTitle())
	}

	// CleanTitle
	if got := CleanTitle("  Hello   World  "); got != "Hello World" {
		t.Fatalf("CleanTitle failed, got %q", got)
	}

	// GetShortNames
	names := GetShortNames("John John Doe  ")
	if len(names) != 2 || names[0].Name != "John" || names[1].Name != "Doe" {
		t.Fatalf("GetShortNames failed, got %#v", names)
	}

	// DetermineShortTitle
	contacts := map[string]*ContactBrief{
		"c1": {ShortTitle: "John"},
	}
	bUnique := &ContactBrief{Names: &person.NameFields{FirstName: "UniqueFirst"}}
	if bUnique.DetermineShortTitle("", contacts) != "" || bUnique.ShortTitle != "UniqueFirst" {
		t.Fatalf("expected ShortTitle to be set to 'UniqueFirst', got %q", bUnique.ShortTitle)
	}

	bNick := &ContactBrief{Names: &person.NameFields{NickName: "Johnny"}}
	if bNick.DetermineShortTitle("", contacts) != "Johnny" {
		t.Fatalf("expected 'Johnny', got %q", bNick.DetermineShortTitle("", contacts))
	}

	bFull := &ContactBrief{Names: &person.NameFields{FullName: "John Doe"}}
	if bFull.DetermineShortTitle("", contacts) != "Doe" {
		t.Fatalf("expected 'Doe' (since John is taken), got %q", bFull.DetermineShortTitle("", contacts))
	}

	bTitleOnly := &ContactBrief{Names: &person.NameFields{}}
	if bTitleOnly.DetermineShortTitle("John Smith", contacts) != "Smith" {
		t.Fatalf("expected 'Smith', got %q", bTitleOnly.DetermineShortTitle("John Smith", contacts))
	}

	bEmpty := &ContactBrief{Names: &person.NameFields{}}
	if bEmpty.DetermineShortTitle("", contacts) != "" {
		t.Fatalf("expected empty string on empty inputs, got %q", bEmpty.DetermineShortTitle("", contacts))
	}
}

func TestContactBase_ValidateAndEqual(t *testing.T) {
	c1 := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John"},
		Status:       ContactStatusActive,
		VATNumber:    "123",
	}
	c2 := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John"},
		Status:       ContactStatusActive,
		VATNumber:    "123",
	}
	if !c1.Equal(c2) {
		t.Fatal("expected c1.Equal(c2) to be true")
	}
	c2.VATNumber = "456"
	if c1.Equal(c2) {
		t.Fatal("expected c1.Equal(c2) to be false on different VAT")
	}

	// Valid ContactBase
	validBase := &ContactBase{
		ContactBrief: ContactBrief{
			Type:     ContactTypePerson,
			Title:    "John",
			Gender:   dbmodels.GenderMale,
			AgeGroup: dbmodels.AgeGroupAdult,
		},
		Status: ContactStatusActive,
	}
	if err := validBase.Validate(); err != nil {
		t.Fatalf("expected valid ContactBase, got %v", err)
	}

	// Missing status
	badStatus := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John"},
	}
	if err := badStatus.Validate(); err == nil {
		t.Fatal("expected error on missing status")
	}

	// Unknown status
	badStatus.Status = "unknown_status"
	if err := badStatus.Validate(); err == nil {
		t.Fatal("expected error on unknown status")
	}

	// Company type requires CountryID
	comp := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypeCompany, Title: "Acme Corp"},
		Status:       ContactStatusActive,
	}
	if err := comp.Validate(); err == nil {
		t.Fatal("expected error on company without CountryID")
	}
	comp.CountryID = "US"
	if err := comp.Validate(); err != nil {
		t.Fatalf("expected valid company, got %v", err)
	}
	// Company with gender or ageGroup is error
	comp.Gender = dbmodels.GenderMale
	if err := comp.Validate(); err == nil {
		t.Fatal("expected error on company with gender")
	}
	comp.Gender = ""
	comp.AgeGroup = dbmodels.AgeGroupAdult
	if err := comp.Validate(); err == nil {
		t.Fatal("expected error on company with ageGroup")
	}

	// Person with VATNumber
	personWithVAT := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John"},
		Status:       ContactStatusActive,
		VATNumber:    "VAT123",
	}
	if err := personWithVAT.Validate(); err == nil {
		t.Fatal("expected error on person with VATNumber")
	}

	// Date of birth
	badDobFormat := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John", DoB: "invalid-date"},
		Status:       ContactStatusActive,
	}
	if err := badDobFormat.Validate(); err == nil {
		t.Fatal("expected error on invalid DoB format")
	}
	futureDob := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John", DoB: time.Now().AddDate(10, 0, 0).Format(time.DateOnly)},
		Status:       ContactStatusActive,
	}
	if err := futureDob.Validate(); err == nil {
		t.Fatal("expected error on future DoB")
	}
	pastDob := &ContactBase{
		ContactBrief: ContactBrief{
			Type:     ContactTypePerson,
			Title:    "John",
			Gender:   dbmodels.GenderMale,
			AgeGroup: dbmodels.AgeGroupAdult,
			DoB:      "1990-01-01",
		},
		Status: ContactStatusActive,
	}
	if err := pastDob.Validate(); err != nil {
		t.Fatalf("expected valid past DoB, got %v", err)
	}

	// Avatars validation
	badAvatar := &ContactBase{
		ContactBrief: ContactBrief{
			Type:     ContactTypePerson,
			Title:    "John",
			Gender:   dbmodels.GenderMale,
			AgeGroup: dbmodels.AgeGroupAdult,
		},
		Status:  ContactStatusActive,
		Avatars: []dbprofile.Avatar{{URL: " http://example.com/avatar.png "}},
	}
	if err := badAvatar.Validate(); err == nil {
		t.Fatal("expected error on bad avatar")
	}
}

type testBrief struct {
	ContactBrief
}

func (b *testBrief) GetRelatedAs() string {
	if b == nil {
		return ""
	}
	return b.ContactBrief.GetRelatedAs()
}

func (b *testBrief) Equal(v *testBrief) bool {
	if b == nil || v == nil {
		return b == v
	}
	return b.ContactBrief.Equal(&v.ContactBrief)
}

func TestWithContactBriefs_Methods(t *testing.T) {
	wb := &WithContactsBase[*testBrief]{}

	// SetContactBrief with empty id panics
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on empty id")
			}
		}()
		wb.SetContactBrief("   ", &testBrief{})
	}()

	tb := &testBrief{ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John", WithUserID: dbmodels.WithUserID{UserID: "u1"}, RolesField: with.RolesField{Roles: []string{"admin"}}}}
	wb.SetContactBrief("c1", tb)
	if wb.GetContactBriefByContactID("c1").Title != "John" {
		t.Fatal("expected c1 title to be John")
	}

	id, brief := wb.GetContactBriefByUserID("u1")
	if id != "c1" || brief.Title != "John" {
		t.Fatalf("expected c1/John, got %s/%s", id, brief.Title)
	}
	id, _ = wb.GetContactBriefByUserID("non-existent")
	if id != "" {
		t.Fatal("expected empty id for non-existent user")
	}

	// Roles
	byRole := wb.GetContactBriefsByRoles("admin")
	if len(byRole) != 1 || byRole["c1"].Title != "John" {
		t.Fatalf("GetContactBriefsByRoles failed, got %#v", byRole)
	}
	sortedByRole := wb.GetSortedContactBriefsByRoles("admin")
	if len(sortedByRole) != 1 || sortedByRole[0].Title != "John" {
		t.Fatalf("GetSortedContactBriefsByRoles failed, got %#v", sortedByRole)
	}
	if count := wb.GetContactsCount("admin"); count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	// Validate WithContactsBase
	wb.WithUserIDs = dbmodels.WithUserIDs{UserIDs: []string{"u1"}}
	if err := wb.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	// Missing UserID in WithUserIDs
	tbEmptyUID := &testBrief{ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "John", Gender: dbmodels.GenderMale, AgeGroup: dbmodels.AgeGroupAdult}}
	wb.SetContactBrief("c2", tbEmptyUID)
	if err := wb.Validate(); err == nil {
		t.Fatal("expected validation error when contact has UserID not in WithUserIDs")
	}
}

func TestWithSingleSpaceContactsWithoutContactIDs_Methods(t *testing.T) {
	s := &WithSingleSpaceContactsWithoutContactIDs[*testBrief]{}

	tb := &testBrief{ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "Alice"}}
	s.AddContact("contact_1", tb)
	if !s.HasContact("contact_1") {
		t.Fatal("expected HasContact to be true")
	}
	if len(s.ContactIDs()) != 1 || s.ContactIDs()[0] != "contact_1" {
		t.Fatalf("expected ['contact_1'], got %v", s.ContactIDs())
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Invalid contact ID
	s.AddContact("bad id!", tb)
	if err := s.Validate(); err == nil {
		t.Fatal("expected error on invalid contact ID")
	}

	s.RemoveContact("bad id!")
	if s.HasContact("bad id!") {
		t.Fatal("expected bad id to be removed")
	}
}

func TestWithMultiSpaceContacts_Methods(t *testing.T) {
	m := &WithMultiSpaceContacts[*testBrief]{}

	tb := &testBrief{ContactBrief: ContactBrief{
		Type:        ContactTypePerson,
		Title:       "Bob",
		WithUserID:  dbmodels.WithUserID{UserID: "u_bob"},
		WithOptionalRelatedAs: dbmodels.WithOptionalRelatedAs{RelatedAs: "parent"},
	}}

	// AddContact
	m.AddContact(coretypes.SpaceID("s1"), "c1", tb)
	if len(m.ContactIDs) != 2 || m.ContactIDs[0] != "*" || m.ContactIDs[1] != "s1_c1" {
		t.Fatalf("expected ['*', 's1_c1'], got %v", m.ContactIDs)
	}

	// ParentContactBrief
	idx, parentID, parentBrief := m.ParentContactBrief()
	if idx != 1 || parentID != "s1_c1" || parentBrief.Title != "Bob" {
		t.Fatalf("ParentContactBrief failed, got %d, %s, %s", idx, parentID, parentBrief.Title)
	}

	// GetContactBriefByID
	idx, brief := m.GetContactBriefByID(coretypes.SpaceID("s1"), "c1")
	if idx != 1 || brief.Title != "Bob" {
		t.Fatalf("GetContactBriefByID failed, got %d, %s", idx, brief.Title)
	}
	idx, _ = m.GetContactBriefByID(coretypes.SpaceID("s1"), "unknown")
	if idx != -1 {
		t.Fatalf("expected -1 for unknown, got %d", idx)
	}

	// GetContactBriefByUserID
	cID, uBrief := m.GetContactBriefByUserID("u_bob")
	if cID != "s1_c1" || uBrief.Title != "Bob" {
		t.Fatalf("GetContactBriefByUserID failed, got %s, %s", cID, uBrief.Title)
	}
	cID, _ = m.GetContactBriefByUserID("unknown")
	if cID != "" {
		t.Fatalf("expected empty for unknown, got %s", cID)
	}

	// Updates
	up := m.Updates()
	if len(up) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(up))
	}
	upItem := m.Updates(dbmodels.SpaceItemID("s1_c1"))
	if len(upItem) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(upItem))
	}

	// SetContactBrief
	m.SetContactBrief(coretypes.SpaceID("s1"), "c1", tb) // same brief
	tbChanged := &testBrief{ContactBrief: tb.ContactBrief}
	tbChanged.Gender = dbmodels.GenderFemale
	m.SetContactBrief(coretypes.SpaceID("s1"), "c1", tbChanged)
	if m.Contacts["s1_c1"].Gender != dbmodels.GenderFemale {
		t.Fatalf("expected updated gender, got %v", m.Contacts["s1_c1"].Gender)
	}

	// RemoveContact
	m.RemoveContact(coretypes.SpaceID("s1"), "c1")
	if _, ok := m.Contacts["s1_c1"]; ok {
		t.Fatal("expected s1_c1 to be removed from Contacts")
	}
	// ParentContactBrief when no parent
	idx, _, _ = m.ParentContactBrief()
	if idx != -1 {
		t.Fatalf("expected -1 when no parent, got %d", idx)
	}
}

func TestRemaining19UncoveredBlocks(t *testing.T) {
	// 1. WithMultiSpaceContactIDs lines 75, 78
	multiEmptySpace := &WithMultiSpaceContactIDs{WithContactIDs: WithContactIDs{ContactIDs: []string{"*", "_contact1"}}}
	if err := multiEmptySpace.Validate(); err == nil {
		t.Fatal("expected error on empty spaceID")
	}
	multiEmptyContact := &WithMultiSpaceContactIDs{WithContactIDs: WithContactIDs{ContactIDs: []string{"*", "space1_"}}}
	if err := multiEmptyContact.Validate(); err == nil {
		t.Fatal("expected error on empty contactID")
	}

	// 2. ContactBrief error branches
	bBadGender := &ContactBrief{Type: ContactTypePerson, Title: "T", Gender: "invalid_gender"}
	if err := bBadGender.Validate(); err == nil {
		t.Fatal("expected error on invalid gender")
	}
	bBadCountry := &ContactBrief{Type: ContactTypePerson, Title: "T", OptionalCountryID: with.OptionalCountryID{CountryID: "invalid_country"}}
	if err := bBadCountry.Validate(); err == nil {
		t.Fatal("expected error on invalid country")
	}
	bBadRoles := &ContactBrief{Type: ContactTypePerson, Title: "T", RolesField: with.RolesField{Roles: []string{" untrimmed "}}}
	if err := bBadRoles.Validate(); err == nil {
		t.Fatal("expected error on untrimmed role")
	}
	bBadUID := &ContactBrief{Type: ContactTypePerson, Title: "T", WithUserID: dbmodels.WithUserID{UserID: " untrimmed "}}
	if err := bBadUID.Validate(); err == nil {
		t.Fatal("expected error on untrimmed UserID")
	}

	// 3. getShortTitle all taken (line 173) and GetShortNames empty name (line 190)
	contactsTaken := map[string]*ContactBrief{
		"c1": {ShortTitle: "John"},
		"c2": {ShortTitle: "Doe"},
	}
	bTaken := &ContactBrief{Names: &person.NameFields{FullName: "John Doe"}}
	if bTaken.DetermineShortTitle("", contactsTaken) != "" {
		t.Fatal("expected empty string when all short names are taken")
	}
	_ = GetShortNames("") // triggers strings.Split("", " ") -> name == "" -> continue

	// 4. ContactBase error branches
	cbBriefErr := &ContactBase{ContactBrief: ContactBrief{Type: "invalid"}, Status: ContactStatusActive}
	if err := cbBriefErr.Validate(); err == nil {
		t.Fatal("expected error on invalid ContactBrief")
	}
	cbNoTitleNames := &ContactBase{ContactBrief: ContactBrief{Type: ContactTypePerson, Title: " "}, Status: ContactStatusActive}
	if err := cbNoTitleNames.Validate(); err == nil {
		t.Fatal("expected error on missing name|title in ContactBase")
	}
	cbBadNames := &ContactBase{ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "T", Names: &person.NameFields{FirstName: " bad "}}, Status: ContactStatusActive}
	if err := cbBadNames.Validate(); err == nil {
		t.Fatal("expected error on bad Names in ContactBase")
	}
	cbCompBadGender := &ContactBase{ContactBrief: ContactBrief{Type: ContactTypeCompany, Title: "C", OptionalCountryID: with.OptionalCountryID{CountryID: "US"}, Gender: "bad_gender"}, Status: ContactStatusActive}
	if err := cbCompBadGender.Validate(); err == nil {
		t.Fatal("expected error on company bad gender")
	}
	cbCompBadAge := &ContactBase{ContactBrief: ContactBrief{Type: ContactTypeCompany, Title: "C", OptionalCountryID: with.OptionalCountryID{CountryID: "US"}, AgeGroup: "bad_age"}, Status: ContactStatusActive}
	if err := cbCompBadAge.Validate(); err == nil {
		t.Fatal("expected error on company bad age group")
	}
	cbBadGroups := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "T", Gender: dbmodels.GenderMale, AgeGroup: dbmodels.AgeGroupAdult},
		Status:       ContactStatusActive,
		WithGroupIDs: WithGroupIDs{GroupIDs: []string{" bad "}},
	}
	if err := cbBadGroups.Validate(); err == nil {
		t.Fatal("expected error on bad WithGroupIDs in ContactBase")
	}
	cbBadTz := &ContactBase{
		ContactBrief: ContactBrief{Type: ContactTypePerson, Title: "T", Gender: dbmodels.GenderMale, AgeGroup: dbmodels.AgeGroupAdult},
		Status:       ContactStatusActive,
		WithTimezone: dbmodels.WithTimezone{Timezone: &dbmodels.Timezone{Iana: " bad "}},
	}
	if err := cbBadTz.Validate(); err == nil {
		t.Fatal("expected error on bad Timezone in ContactBase")
	}

	// 5. WithContactsBase.Validate when contact.Validate() fails
	wb := &WithContactsBase[*testBrief]{}
	badBrief := &testBrief{ContactBrief: ContactBrief{Type: "invalid_type"}}
	wb.SetContactBrief("c1", badBrief)
	if err := wb.Validate(); err == nil {
		t.Fatal("expected error when contact.Validate() fails in WithContactsBase")
	}

	// 6. WithSingleSpaceContactsWithoutContactIDs.Validate when brief.Validate() fails
	s := &WithSingleSpaceContactsWithoutContactIDs[*testBrief]{}
	s.AddContact("c1", &testBrief{ContactBrief: ContactBrief{Type: "invalid"}})
	if err := s.Validate(); err == nil {
		t.Fatal("expected error on brief.Validate failure in WithSingleSpaceContactsWithoutContactIDs")
	}

	// 7. ContactBrief.Validate when AccountsOfUser.Validate() fails
	bBadAcc := &ContactBrief{
		Type:           ContactTypePerson,
		Title:          "T",
		AccountsOfUser: appuser.AccountsOfUser{Accounts: []string{"invalid_account"}},
	}
	if err := bBadAcc.Validate(); err == nil {
		t.Fatal("expected error on invalid AccountsOfUser")
	}
}


