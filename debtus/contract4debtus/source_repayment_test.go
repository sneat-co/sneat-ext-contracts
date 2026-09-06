package contract4debtus

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func validSourceRepaymentRequest() RecordSourceObligationRepaymentRequest {
	return RecordSourceObligationRepaymentRequest{
		ContractVersion: SourceRepaymentContractVersion,
		Source: SourceRef{
			Namespace: SourceNamespaceSplitus,
			SpaceID:   "house-1",
			RecordID:  "bill-1",
		},
		LineID:       "bea-to-alex",
		ObligationID: "obligation-1",
		Currency:     "EUR",
		AmountMinor:  "3000",
		RepaidAt:     time.Date(2026, 9, 6, 8, 30, 0, 0, time.UTC),
		OperationKey: "bill-1/repayment-1",
		ActorUserID:  "bea-user",
	}
}

func validSourceRepaymentResult(request RecordSourceObligationRepaymentRequest) RecordSourceObligationRepaymentResult {
	canRepay := SourceObligationRepaymentCapability{CanRecordRepayment: true}
	return RecordSourceObligationRepaymentResult{
		ContractVersion: SourceRepaymentContractVersion,
		Source:          request.Source,
		OperationKey:    request.OperationKey,
		Obligation: SourceObligationStatus{
			LineID: request.LineID, ObligationIDs: []string{request.ObligationID},
			Debtor:   ContactRef{SpaceID: request.Source.SpaceID, ContactID: "bea"},
			Creditor: ContactRef{SpaceID: request.Source.SpaceID, ContactID: "alex"},
			Currency: request.Currency, PrincipalMinor: 6000, OutstandingMinor: 3000,
			RepaidMinor: 3000, Status: SettlementStatusPartSettled,
			RepaymentCapability: &canRepay,
		},
		Repayment: SourceObligationActivity{
			ActivityID: "repayment-1", RootActivityID: request.ObligationID,
			LineIDs: []string{request.LineID}, Kind: SourceActivityRepayment,
			From:     ContactRef{SpaceID: request.Source.SpaceID, ContactID: "bea"},
			To:       ContactRef{SpaceID: request.Source.SpaceID, ContactID: "alex"},
			Currency: request.Currency, AmountMinor: 3000,
			ActorUserID: request.ActorUserID, OccurredAt: request.RepaidAt,
		},
	}
}

func TestRecordSourceObligationRepaymentRequestValidate(t *testing.T) {
	request := validSourceRepaymentRequest()
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got, err := request.AmountMinor.MinorUnits(); err != nil || got != 3000 {
		t.Fatalf("MinorUnits() = %d, %v; want 3000, nil", got, err)
	}
}

func TestSourceRepaymentWholeSecondBrowserWireExample(t *testing.T) {
	data, err := os.ReadFile("testdata/source_repayment_v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden struct {
		Request struct {
			RepaidAt string `json:"repaidAt"`
		} `json:"request"`
		Response struct {
			Repayment struct {
				OccurredAt string `json:"occurredAt"`
			} `json:"repayment"`
		} `json:"response"`
	}
	if err = json.Unmarshal(data, &golden); err != nil {
		t.Fatal(err)
	}
	instant := time.Date(2026, 9, 6, 8, 30, 0, 0, time.UTC)
	formatted, err := FormatSourceRepaymentBrowserTime(instant)
	if err != nil {
		t.Fatal(err)
	}
	if formatted != golden.Request.RepaidAt || formatted != golden.Response.Repayment.OccurredAt {
		t.Fatalf("browser time = %q; request=%q response=%q", formatted, golden.Request.RepaidAt, golden.Response.Repayment.OccurredAt)
	}
	defaultJSON, err := json.Marshal(instant)
	if err != nil {
		t.Fatal(err)
	}
	if string(defaultJSON) == `"`+formatted+`"` || string(defaultJSON) != `"2026-09-06T08:30:00Z"` {
		t.Fatalf("default time JSON = %s; test no longer proves the interop hazard", defaultJSON)
	}
}

func TestRecordSourceObligationRepaymentRequestRejectsInvalidShape(t *testing.T) {
	tests := map[string]func(*RecordSourceObligationRepaymentRequest){
		"unsupported contract version": func(r *RecordSourceObligationRepaymentRequest) { r.ContractVersion = 2 },
		"unsafe source namespace":      func(r *RecordSourceObligationRepaymentRequest) { r.Source.Namespace = "Splitus" },
		"padded source record":         func(r *RecordSourceObligationRepaymentRequest) { r.Source.RecordID = " bill-1" },
		"control in line ID":           func(r *RecordSourceObligationRepaymentRequest) { r.LineID = "line\n1" },
		"invalid UTF-8 obligation ID":  func(r *RecordSourceObligationRepaymentRequest) { r.ObligationID = string([]byte{0xff}) },
		"lowercase currency":           func(r *RecordSourceObligationRepaymentRequest) { r.Currency = "eur" },
		"zero amount":                  func(r *RecordSourceObligationRepaymentRequest) { r.AmountMinor = "0" },
		"leading-zero amount":          func(r *RecordSourceObligationRepaymentRequest) { r.AmountMinor = "03000" },
		"overflow amount":              func(r *RecordSourceObligationRepaymentRequest) { r.AmountMinor = "9223372036854775808" },
		"missing timestamp":            func(r *RecordSourceObligationRepaymentRequest) { r.RepaidAt = time.Time{} },
		"non-UTC timestamp": func(r *RecordSourceObligationRepaymentRequest) {
			r.RepaidAt = time.Date(2026, 9, 6, 9, 30, 0, 0, time.FixedZone("plus-one", 3600))
		},
		"sub-millisecond timestamp": func(r *RecordSourceObligationRepaymentRequest) {
			r.RepaidAt = r.RepaidAt.Add(time.Nanosecond)
		},
		"missing operation key": func(r *RecordSourceObligationRepaymentRequest) { r.OperationKey = "" },
		"missing trusted actor": func(r *RecordSourceObligationRepaymentRequest) { r.ActorUserID = "" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			request := validSourceRepaymentRequest()
			mutate(&request)
			if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
			}
		})
	}
}

func TestDecodeRecordSourceObligationRepaymentRequestRejectsNumericAmountAndUnknownFields(t *testing.T) {
	valid := `{"contractVersion":1,"source":{"namespace":"splitus","spaceID":"house-1","recordID":"bill-1"},"lineID":"line-1","obligationID":"obligation-1","currency":"EUR","amountMinor":"3000","repaidAt":"2026-09-06T08:30:00.000Z","operationKey":"repayment-1"}`
	request, err := DecodeRecordSourceObligationRepaymentRequest(strings.NewReader(valid), "bea-user")
	if err != nil || request.AmountMinor != "3000" || request.ActorUserID != "bea-user" {
		t.Fatalf("Decode() = %+v, %v", request, err)
	}
	encoded, err := json.Marshal(request)
	if err != nil || bytes.Contains(encoded, []byte("actorUserID")) || bytes.Contains(encoded, []byte("bea-user")) {
		t.Fatalf("trusted actor leaked into browser JSON: %s, %v", encoded, err)
	}

	for name, payload := range map[string]string{
		"numeric amount": strings.Replace(valid, `"3000"`, `3000`, 1),
		"overflow":       strings.Replace(valid, `"3000"`, `"9223372036854775808"`, 1),
		"browser actor":  valid[:len(valid)-1] + `,"actorUserID":"attacker"}`,
		"unknown field":  valid[:len(valid)-1] + `,"authority":true}`,
		"trailing JSON":  valid + `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, decodeErr := DecodeRecordSourceObligationRepaymentRequest(bytes.NewBufferString(payload), "bea-user"); !errors.Is(decodeErr, ErrInvalidRequest) {
				t.Fatalf("Decode() error = %v, want ErrInvalidRequest", decodeErr)
			}
		})
	}
}

func TestSourceObligationRepaymentCapabilityValidate(t *testing.T) {
	valid := []SourceObligationRepaymentCapability{
		{CanRecordRepayment: true},
		{RepaymentUnavailableReason: SourceRepaymentReadOnlyRole},
		{RepaymentUnavailableReason: SourceRepaymentNotParty},
		{RepaymentUnavailableReason: SourceRepaymentNotOutstanding},
		{RepaymentUnavailableReason: SourceRepaymentLedgerPartyUnresolved},
		{RepaymentUnavailableReason: SourceRepaymentLinkedEventLimitReached},
	}
	for _, capability := range valid {
		if err := capability.Validate(); err != nil {
			t.Fatalf("Validate(%+v) error = %v", capability, err)
		}
	}
	for _, invalid := range []SourceObligationRepaymentCapability{
		{},
		{CanRecordRepayment: true, RepaymentUnavailableReason: SourceRepaymentNotParty},
		{RepaymentUnavailableReason: "serverMessage"},
	} {
		if err := invalid.Validate(); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("Validate(%+v) error = %v, want ErrInvalidRequest", invalid, err)
		}
	}
}

func TestRecordSourceObligationRepaymentResultValidateFor(t *testing.T) {
	request := validSourceRepaymentRequest()
	result := validSourceRepaymentResult(request)
	if err := result.ValidateFor(request); err != nil {
		t.Fatalf("ValidateFor() error = %v", err)
	}

	tests := map[string]func(*RecordSourceObligationRepaymentResult){
		"different source":        func(r *RecordSourceObligationRepaymentResult) { r.Source.RecordID = "bill-2" },
		"different operation key": func(r *RecordSourceObligationRepaymentResult) { r.OperationKey = "repayment-2" },
		"missing capability":      func(r *RecordSourceObligationRepaymentResult) { r.Obligation.RepaymentCapability = nil },
		"different line":          func(r *RecordSourceObligationRepaymentResult) { r.Obligation.LineID = "line-2" },
		"missing obligation":      func(r *RecordSourceObligationRepaymentResult) { r.Obligation.ObligationIDs = []string{"obligation-2"} },
		"wrong activity root":     func(r *RecordSourceObligationRepaymentResult) { r.Repayment.RootActivityID = "obligation-2" },
		"wrong actor":             func(r *RecordSourceObligationRepaymentResult) { r.Repayment.ActorUserID = "alex-user" },
		"wrong repayment from":    func(r *RecordSourceObligationRepaymentResult) { r.Repayment.From.ContactID = "cam" },
		"wrong repayment to":      func(r *RecordSourceObligationRepaymentResult) { r.Repayment.To.ContactID = "cam" },
		"wrong amount":            func(r *RecordSourceObligationRepaymentResult) { r.Repayment.AmountMinor = 2999 },
		"wrong timestamp": func(r *RecordSourceObligationRepaymentResult) {
			r.Repayment.OccurredAt = r.Repayment.OccurredAt.Add(time.Second)
		},
		"sub-millisecond timestamp": func(r *RecordSourceObligationRepaymentResult) {
			r.Repayment.OccurredAt = r.Repayment.OccurredAt.Add(time.Nanosecond)
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			candidate := validSourceRepaymentResult(request)
			mutate(&candidate)
			if err := candidate.ValidateFor(request); !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("ValidateFor() error = %v, want ErrInvalidRequest", err)
			}
		})
	}
}
