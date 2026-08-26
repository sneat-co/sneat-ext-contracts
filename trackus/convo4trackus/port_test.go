package convo4trackus

import "testing"

func TestAddEntryContractCarriesEntryID(t *testing.T) {
	const entryID = "CommittedMeBot:chat-123:message-456"
	request := AddEntryRequest{
		SpaceID:   "space-1",
		TrackerID: "_push_ups",
		Value:     IntValue(20),
		EntryID:   entryID,
		Date:      "2026-07-26",
	}
	if request.EntryID != entryID {
		t.Fatalf("request EntryID = %q, want %q", request.EntryID, entryID)
	}

	result := EntryResult{
		EntryID:      request.EntryID,
		TrackerID:    request.TrackerID,
		TrackerTitle: "Push-ups",
		Total:        20,
		IsSummable:   true,
	}
	if result.EntryID != entryID {
		t.Fatalf("result EntryID = %q, want %q", result.EntryID, entryID)
	}
}

func TestAddEntryContractKeepsLegacyKeyedLiteralsCompatible(t *testing.T) {
	// Existing callers using keyed literals do not need to opt into
	// idempotency immediately: the added fields have useful zero values.
	request := AddEntryRequest{
		SpaceID:   "space-1",
		TrackerID: "_push_ups",
		Value:     IntValue(20),
		Date:      "2026-07-26",
	}
	if request.EntryID != "" {
		t.Fatalf("legacy request EntryID = %q, want empty", request.EntryID)
	}

	result := EntryResult{
		TrackerID:    request.TrackerID,
		TrackerTitle: "Push-ups",
		Total:        20,
		IsSummable:   true,
	}
	if result.EntryID != "" {
		t.Fatalf("legacy result EntryID = %q, want empty", result.EntryID)
	}
}
