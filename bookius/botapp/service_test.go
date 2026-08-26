package botapp

import "testing"

func TestInvoiceValidate(t *testing.T) {
	valid := Invoice{
		CommitmentID: "outing:option:user",
		AmountCents:  1250,
		Currency:     "EUR",
		Description:  "Lunch",
	}

	tests := []struct {
		name    string
		invoice Invoice
		wantErr bool
	}{
		{name: "valid", invoice: valid},
		{name: "missing commitment ID", invoice: Invoice{AmountCents: valid.AmountCents, Currency: valid.Currency, Description: valid.Description}, wantErr: true},
		{name: "zero amount", invoice: Invoice{CommitmentID: valid.CommitmentID, Currency: valid.Currency, Description: valid.Description}, wantErr: true},
		{name: "missing currency", invoice: Invoice{CommitmentID: valid.CommitmentID, AmountCents: valid.AmountCents, Description: valid.Description}, wantErr: true},
		{name: "missing description", invoice: Invoice{CommitmentID: valid.CommitmentID, AmountCents: valid.AmountCents, Currency: valid.Currency}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.invoice.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
