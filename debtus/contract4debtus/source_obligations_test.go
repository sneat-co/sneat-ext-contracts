package contract4debtus

import (
	"errors"
	"math"
	"testing"
)

func validRequest(t *testing.T) ReconcileSourceObligationsRequest {
	t.Helper()
	r := ReconcileSourceObligationsRequest{
		Source:                   SourceRef{Namespace: SourceNamespaceSplitus, SpaceID: "space-1", RecordID: "expense-1"},
		RecorderUserID:           "recorder-1",
		ExpectedPreviousRevision: 3,
		NewRevision:              4,
		OperationKey:             "expense-1/revision-4",
		DesiredLines: []ObligationLine{
			{LineID: "line-b", Debtor: ContactRef{SpaceID: "space-1", ContactID: "bob"}, Creditor: ContactRef{SpaceID: "space-1", ContactID: "alice"}, Currency: "EUR", AmountMinor: 1250},
			{LineID: "line-a", Debtor: ContactRef{SpaceID: "space-1", ContactID: "carol"}, Creditor: ContactRef{SpaceID: "space-1", ContactID: "alice"}, Currency: "EUR", AmountMinor: 999},
		},
	}
	digest, err := r.CanonicalInputDigest()
	if err != nil {
		t.Fatal(err)
	}
	r.InputDigest = digest
	return r
}

func TestReconcileSourceObligationsRequestValidate(t *testing.T) {
	r := validRequest(t)
	if err := r.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsMalformedAndCrossSpaceIdentities(t *testing.T) {
	tests := map[string]func(*ReconcileSourceObligationsRequest){
		"empty recorder":       func(r *ReconcileSourceObligationsRequest) { r.RecorderUserID = "" },
		"padded operation key": func(r *ReconcileSourceObligationsRequest) { r.OperationKey = " padded" },
		"bad namespace":        func(r *ReconcileSourceObligationsRequest) { r.Source.Namespace = "Splitus" },
		"debtor cross space":   func(r *ReconcileSourceObligationsRequest) { r.DesiredLines[0].Debtor.SpaceID = "foreign" },
		"creditor cross space": func(r *ReconcileSourceObligationsRequest) { r.DesiredLines[0].Creditor.SpaceID = "foreign" },
		"same party":           func(r *ReconcileSourceObligationsRequest) { r.DesiredLines[0].Debtor.ContactID = "alice" },
		"lowercase currency":   func(r *ReconcileSourceObligationsRequest) { r.DesiredLines[0].Currency = "eur" },
		"fraction substitute":  func(r *ReconcileSourceObligationsRequest) { r.DesiredLines[0].AmountMinor = 0 },
		"revision gap":         func(r *ReconcileSourceObligationsRequest) { r.NewRevision = 5 },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			r := validRequest(t)
			mutate(&r)
			digest, err := r.CanonicalInputDigest()
			if err != nil {
				t.Fatal(err)
			}
			r.InputDigest = digest
			if err = r.Validate(); !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
			}
		})
	}
}

func TestValidateRejectsDuplicateLinesAndOverflow(t *testing.T) {
	t.Run("duplicate", func(t *testing.T) {
		r := validRequest(t)
		r.DesiredLines[1].LineID = r.DesiredLines[0].LineID
		r.InputDigest, _ = r.CanonicalInputDigest()
		if err := r.Validate(); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
		}
	})
	t.Run("same currency total overflow", func(t *testing.T) {
		r := validRequest(t)
		r.DesiredLines[0].AmountMinor = math.MaxInt64
		r.DesiredLines[1].AmountMinor = 1
		r.InputDigest, _ = r.CanonicalInputDigest()
		if err := r.Validate(); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
		}
	})
}

func TestCanonicalInputDigestIsDeterministicAndSensitive(t *testing.T) {
	r := validRequest(t)
	want := r.InputDigest
	r.DesiredLines[0], r.DesiredLines[1] = r.DesiredLines[1], r.DesiredLines[0]
	got, err := r.CanonicalInputDigest()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("digest changed with line order: got %s, want %s", got, want)
	}

	r.DesiredLines[0].AmountMinor++
	changed, _ := r.CanonicalInputDigest()
	if changed == want {
		t.Fatal("digest did not change with exact minor-unit amount")
	}
	r.InputDigest = want
	if err = r.Validate(); !errors.Is(err, ErrDigestMismatch) {
		t.Fatalf("Validate() error = %v, want ErrDigestMismatch", err)
	}
}

func TestCanonicalDigestHasFixedLowercaseSHA256Shape(t *testing.T) {
	r := validRequest(t)
	if len(r.InputDigest) != 64 {
		t.Fatalf("digest length = %d, want 64", len(r.InputDigest))
	}
	for _, c := range r.InputDigest {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			t.Fatalf("digest contains non-lowercase-hex character %q", c)
		}
	}
}
