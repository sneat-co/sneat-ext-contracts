package dto4schoolus

import (
	"encoding/json"
	"testing"
)

// TestMarkAttendanceRequestJSONRoundTrip freezes the wire field names.
func TestMarkAttendanceRequestJSONRoundTrip(t *testing.T) {
	in := `{"studentUID":"u123","date":"2026-09-01","status":"late","note":"bus delay"}`
	var r MarkAttendanceRequest
	if err := json.Unmarshal([]byte(in), &r); err != nil {
		t.Fatal(err)
	}
	if r.StudentUID != "u123" || r.Date != "2026-09-01" || r.Status != AttendanceStatusLate {
		t.Fatalf("unexpected decode: %+v", r)
	}
}

// TestEnrolStudentRequestJSONRoundTrip freezes the enrolment wire shape.
func TestEnrolStudentRequestJSONRoundTrip(t *testing.T) {
	in := `{"studentUID":"u456","studentName":"Ann","guardianUID":"g789"}`
	var r EnrolStudentRequest
	if err := json.Unmarshal([]byte(in), &r); err != nil {
		t.Fatal(err)
	}
	if r.StudentUID != "u456" || r.StudentName != "Ann" || r.GuardianUID != "g789" {
		t.Fatalf("unexpected decode: %+v", r)
	}
}

// TestEnumValues pins the string values of the contract enums.
func TestEnumValues(t *testing.T) {
	cases := map[string]string{
		string(SchoolTypeInternational):        "international",
		string(AttendanceStatusExcused):        "excused",
		string(EnrolmentStatusGraduated):       "graduated",
		string(SchoolRoleVicePrincipal):        "vice_principal",
		string(SchoolRoleTransportCoordinator): "transport_coordinator",
	}
	for got, want := range cases {
		if got != want {
			t.Fatalf("enum value mismatch: got %q want %q", got, want)
		}
	}
}
