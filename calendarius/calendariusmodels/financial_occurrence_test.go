package calendariusmodels

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFinancialOccurrenceQueryValidate(t *testing.T) {
	valid := FinancialOccurrenceQuery{SpaceID: "space1", HappeningID: "h1", FromDate: "2026-09-01", ToDate: "2026-12-02", Timezone: "Europe/Dublin"}
	require.NoError(t, valid.Validate())
	valid.HappeningID = ""
	require.NoError(t, valid.Validate())
	require.Error(t, valid.ValidateForResolve())

	invalid := []FinancialOccurrenceQuery{
		{},
		{SpaceID: "bad/space", HappeningID: "h1", FromDate: "2026-09-01", ToDate: "2026-09-02"},
		{SpaceID: "space1", HappeningID: "bad/happening", FromDate: "2026-09-01", ToDate: "2026-09-02"},
		{SpaceID: "space1", HappeningID: "h1", FromDate: "2026-02-30", ToDate: "2026-03-01"},
		{SpaceID: "space1", HappeningID: "h1", FromDate: "2026-09-02", ToDate: "2026-09-01"},
		{SpaceID: "space1", HappeningID: "h1", FromDate: "2026-09-01", ToDate: "2026-12-03"},
		{SpaceID: "space1", HappeningID: "h1", FromDate: "2026-09-01", ToDate: "2026-09-02", Timezone: "not/a-zone"},
	}
	for _, query := range invalid {
		require.Error(t, query.Validate(), query)
	}
}

func TestValidateFinancialOccurrenceID(t *testing.T) {
	require.NoError(t, ValidateFinancialOccurrenceID("slot1@2026-09-01"))
	for _, value := range []string{"", ".", "../other", "bad/control\n"} {
		require.Error(t, ValidateFinancialOccurrenceID(value))
	}
}
