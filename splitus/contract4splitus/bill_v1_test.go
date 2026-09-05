package contract4splitus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func exact(value string) ExactDecimalString { return ExactDecimalString(value) }

func electricityBill() CreateBillV1Request {
	provider := "Grid Energy"
	expected := exact("80.00")
	standing := exact("15.00")
	previous := &PreviousComparableBillV1{
		BillID: "electricity-2026-07", ActualAmount: exact("85.00"), Comparison: ExpectedActualIncreased,
	}
	return CreateBillV1Request{
		ContractVersion: BillContractVersion,
		SpaceID:         "housemates-space",
		BillID:          "electricity-2026-08",
		RecorderUserID:  "recorder-user",
		BillKind:        BillKindUtility,
		Currency:        "EUR",
		ActualAmount:    exact("90.00"),
		PaidAllocations: []BillAllocationV1{
			{AllocationID: "paid-alex", ContactID: "alex-contact", Amount: exact("90.00")},
		},
		OwedAllocations: []BillAllocationV1{
			{AllocationID: "owed-alex", ContactID: "alex-contact", Amount: exact("30.00")},
			{AllocationID: "owed-bea", ContactID: "bea-contact", Amount: exact("30.00")},
			{AllocationID: "owed-cam", ContactID: "cam-contact", Amount: exact("30.00")},
		},
		Utility: &UtilityDetailsV1{
			UtilityKind:  UtilityKindElectricity,
			ProviderName: &provider,
			Period:       BillingPeriodV1{StartDate: "2026-08-01", EndDate: "2026-08-31"},
		},
		RecurringOccurrence: &RecurringOccurrenceV1{
			HappeningID: "monthly-electricity", OccurrenceID: "2026-08",
			ExpectedAmount: &expected, StandingChargeAmount: &standing,
			ExpectedComparison: ExpectedActualIncreased,
			PreviousComparable: previous,
		},
	}
}

func TestCreateBillV1AcceptsEUR90ForThreeHousemates(t *testing.T) {
	request := electricityBill()
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if request.ActualAmount != exact("90.00") || len(request.OwedAllocations) != 3 {
		t.Fatalf("validated bill = %#v", request)
	}
}

func TestCreateBillV1ReconcilesCustomAllocationsExactly(t *testing.T) {
	request := electricityBill()
	request.PaidAllocations = []BillAllocationV1{
		{AllocationID: "paid-alex", ContactID: "alex-contact", Amount: exact("70.00")},
		{AllocationID: "paid-bea", ContactID: "bea-contact", Amount: exact("20.00")},
	}
	request.OwedAllocations = []BillAllocationV1{
		{AllocationID: "owed-alex", ContactID: "alex-contact", Amount: exact("10.00")},
		{AllocationID: "owed-bea", ContactID: "bea-contact", Amount: exact("20.00")},
		{AllocationID: "owed-cam", ContactID: "cam-contact", Amount: exact("60.00")},
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	request.OwedAllocations[2].Amount = exact("59.99")
	if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}

func TestCreateBillV1KeepsRecorderSeparateFromPayer(t *testing.T) {
	request := electricityBill()
	if request.RecorderUserID == request.PaidAllocations[0].ContactID {
		t.Fatal("test input did not keep recorder user and payer contact identities separate")
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if request.RecorderUserID != "recorder-user" || request.PaidAllocations[0].ContactID != "alex-contact" {
		t.Fatalf("validation changed recorder or payer: %#v", request)
	}
	if err := request.ValidateAuthenticatedRecorder("impersonated-user"); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("ValidateAuthenticatedRecorder() error = %v, want ErrInvalidRequest", err)
	}
	if err := request.ValidateAuthenticatedRecorder("recorder-user"); err != nil {
		t.Fatalf("ValidateAuthenticatedRecorder() error = %v", err)
	}
}

func TestCreateBillV1RejectsUnsafeIDs(t *testing.T) {
	tests := []struct {
		name string
		set  func(*CreateBillV1Request)
	}{
		{name: "slash in Space", set: func(r *CreateBillV1Request) { r.SpaceID = "spaces/foreign" }},
		{name: "relative Space", set: func(r *CreateBillV1Request) { r.SpaceID = ".." }},
		{name: "reserved bill", set: func(r *CreateBillV1Request) { r.BillID = "__reserved__" }},
		{name: "padded bill", set: func(r *CreateBillV1Request) { r.BillID = " bill" }},
		{name: "control in recorder", set: func(r *CreateBillV1Request) { r.RecorderUserID = "user\x00id" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := electricityBill()
			test.set(&request)
			if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
			}
		})
	}
}

func TestDecodeCreateBillV1RejectsNumericJSONAmount(t *testing.T) {
	payload := `{
		"contractVersion":1,
		"spaceID":"housemates-space",
		"billID":"electricity-2026-08",
		"recorderUserID":"recorder-user",
		"billKind":"general",
		"currency":"EUR",
		"actualAmount":90,
		"paidAllocations":[],
		"owedAllocations":[]
	}`
	if _, err := DecodeCreateBillV1Request(strings.NewReader(payload)); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("DecodeCreateBillV1Request() error = %v, want ErrInvalidRequest", err)
	}
}

func TestExactDecimalRejectsOverflow(t *testing.T) {
	if _, err := exact("92233720368547758.08").MinorUnits(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("MinorUnits() error = %v, want ErrInvalidRequest", err)
	}
	request := electricityBill()
	request.ActualAmount = exact("92233720368547758.07")
	request.PaidAllocations = []BillAllocationV1{
		{AllocationID: "paid-alex", ContactID: "alex-contact", Amount: exact("92233720368547758.07")},
		{AllocationID: "paid-bea", ContactID: "bea-contact", Amount: exact("0.01")},
	}
	if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}

func TestDecodeCreateBillV1RejectsExpectedWithoutActual(t *testing.T) {
	request := electricityBill()
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	payload := strings.Replace(string(data), `"actualAmount":"90.00",`, "", 1)
	if _, err := DecodeCreateBillV1Request(strings.NewReader(payload)); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("DecodeCreateBillV1Request() error = %v, want ErrInvalidRequest", err)
	}
}

func TestCreateBillV1BindsExpectedComparisonToActual(t *testing.T) {
	request := electricityBill()
	request.RecurringOccurrence.ExpectedComparison = ExpectedActualDecreased
	if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}

func TestCreateBillV1BindsPreviousComparableToActual(t *testing.T) {
	request := electricityBill()
	request.RecurringOccurrence.PreviousComparable.ActualAmount = exact("95.00")
	if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}

func TestListBillsV1RequestIsBounded(t *testing.T) {
	request := ListBillsV1Request{
		ContractVersion: BillContractVersion,
		SpaceID:         "housemates-space",
		PageSize:        MaxBillListPageSize,
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	request.PageSize++
	if err := request.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}

func appliedBillResponse() CreateBillV1Response {
	request := electricityBill()
	digest := strings.Repeat("a", 64)
	lineID := "bea-to-alex"
	return CreateBillV1Response{
		ContractVersion: BillContractVersion,
		Bill: BillV1{
			CreateBillV1Request: request,
			Revision:            "1",
			Posting: BillPostingV1{
				Status: PostingStatusApplied, OperationKey: "post-electricity-2026-08", InputDigest: digest,
				Receipt: &BillPostingReceiptV1{
					ReceiptID: "debtus-receipt-1", OperationKey: "post-electricity-2026-08", InputDigest: digest, Revision: "1",
					ObligationLines: []DebtusReceiptLineV1{{LineID: lineID, ObligationIDs: []string{"debt-bea-alex"}}},
				},
			},
			Debtus: &DebtusStatusV1{
				Status: DebtusSettlementStatusUnsettled,
				Obligations: []DebtusObligationV1{{
					LineID: lineID, ObligationIDs: []string{"debt-bea-alex"},
					DebtorContactID: "bea-contact", CreditorContactID: "alex-contact", Currency: "EUR",
					PrincipalAmount: exact("30.00"), OutstandingAmount: exact("30.00"), RepaidAmount: exact("0.00"), CreditAmount: exact("0.00"),
					Status:           DebtusSettlementStatusUnsettled,
					SettlementTarget: DebtusSettlementTargetV1{Route: SettlementRoute, SpaceID: request.SpaceID, SourceNamespace: DebtusSourceNamespace, SourceRecordID: request.BillID, LineID: &lineID},
				}},
				SettlementTarget: DebtusSettlementTargetV1{Route: SettlementRoute, SpaceID: request.SpaceID, SourceNamespace: DebtusSourceNamespace, SourceRecordID: request.BillID},
			},
			CreatedAt: "2026-09-05T12:00:00Z",
			UpdatedAt: "2026-09-05T12:00:01Z",
		},
	}
}

func TestCreateBillV1ResponseValidatesPostingAndDebtusProjection(t *testing.T) {
	response := appliedBillResponse()
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	response.Bill.Posting.Receipt.Revision = "2"
	if err := response.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() mismatched receipt error = %v, want ErrInvalidRequest", err)
	}
}

func TestCreateBillV1ResponseRejectsTimestampAndUnboundedAttention(t *testing.T) {
	response := appliedBillResponse()
	response.Bill.UpdatedAt = "not-a-time"
	if err := response.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() timestamp error = %v, want ErrInvalidRequest", err)
	}

	response = appliedBillResponse()
	badCode := BillAttentionCode("raw server error text")
	response.Bill.Posting = BillPostingV1{
		Status: PostingStatusAttention, OperationKey: "post-electricity-2026-08", InputDigest: strings.Repeat("a", 64), AttentionCode: &badCode,
	}
	response.Bill.Debtus = nil
	if err := response.Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() attention code error = %v, want ErrInvalidRequest", err)
	}
}

func TestListBillsV1ResponseRejectsItemsPastPageBound(t *testing.T) {
	item := BillListItemV1{
		ContractVersion: BillContractVersion, SpaceID: "housemates-space", BillID: "electricity-2026-08",
		BillKind: BillKindGeneral, Currency: "EUR", ActualAmount: exact("90.00"), OwnPaidAmount: exact("0.00"), OwnOwedAmount: exact("30.00"),
		PostingStatus: PostingStatusApplied, CreatedAt: "2026-09-05T12:00:00Z",
	}
	response := ListBillsV1Response{ContractVersion: BillContractVersion, PageSize: 1, Items: []BillListItemV1{item, item}}
	if err := response.Validate(1); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRequest", err)
	}
}
