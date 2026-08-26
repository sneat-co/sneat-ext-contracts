// Package dto4schoolus holds the frozen wire DTOs of the schoolus extension
// (product: School Portal) — the contract surface shared between the Go
// backend and other repos. The TypeSpec source of truth is
// ../../typespec/api4schoolus.tsp.
package dto4schoolus

import "time"

// SchoolType enumerates the kind of school.
type SchoolType string

const (
	SchoolTypeLanguage      SchoolType = "language"
	SchoolTypeMusic         SchoolType = "music"
	SchoolTypeTutoring      SchoolType = "tutoring"
	SchoolTypePrimary       SchoolType = "primary"
	SchoolTypeSecondary     SchoolType = "secondary"
	SchoolTypeInternational SchoolType = "international"
	SchoolTypeDance         SchoolType = "dance"
	SchoolTypeDriving       SchoolType = "driving"
	SchoolTypeSports        SchoolType = "sports"
	SchoolTypeTraining      SchoolType = "training"
)

// AttendanceStatus is the attendance mark for a student on a given date.
type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
	AttendanceStatusLate    AttendanceStatus = "late"
	AttendanceStatusExcused AttendanceStatus = "excused"
)

// EnrolmentStatus is the lifecycle of a student's enrolment in a class.
type EnrolmentStatus string

const (
	EnrolmentStatusActive    EnrolmentStatus = "active"
	EnrolmentStatusWithdrawn EnrolmentStatus = "withdrawn"
	EnrolmentStatusGraduated EnrolmentStatus = "graduated"
)

// SchoolRole is a person's role within a school space.
type SchoolRole string

const (
	SchoolRoleOwner                SchoolRole = "owner"
	SchoolRolePrincipal            SchoolRole = "principal"
	SchoolRoleVicePrincipal        SchoolRole = "vice_principal"
	SchoolRoleAdministrator        SchoolRole = "administrator"
	SchoolRoleOfficeStaff          SchoolRole = "office_staff"
	SchoolRoleTeacher              SchoolRole = "teacher"
	SchoolRoleAssistantTeacher     SchoolRole = "assistant_teacher"
	SchoolRoleCoach                SchoolRole = "coach"
	SchoolRoleStudent              SchoolRole = "student"
	SchoolRoleParent               SchoolRole = "parent"
	SchoolRoleGuardian             SchoolRole = "guardian"
	SchoolRoleAccountant           SchoolRole = "accountant"
	SchoolRoleNurse                SchoolRole = "nurse"
	SchoolRoleLibrarian            SchoolRole = "librarian"
	SchoolRoleTransportCoordinator SchoolRole = "transport_coordinator"
	SchoolRoleExternalTutor        SchoolRole = "external_tutor"
	SchoolRoleGuest                SchoolRole = "guest"
)

// CreateSchoolRequest creates a school in the administrator's business space.
// POST /v0/api4schoolus/spaces/{spaceID}/schools
type CreateSchoolRequest struct {
	Title string     `json:"title"`
	Type  SchoolType `json:"type"`
	// Location is a free-text location, e.g. "Dublin, IE".
	Location string `json:"location,omitempty"`
	// Description is shown on the public school page.
	Description string `json:"description,omitempty"`
}

// CreateSchoolResponse identifies the new school so the UI can navigate to it.
type CreateSchoolResponse struct {
	ID string `json:"id"`
	// PublicPath is the site-relative public school page path, e.g. "/school/abc".
	PublicPath string `json:"publicPath"`
}

// SchoolPublic is the unauthenticated public projection of a school.
// GET /v0/api4schoolus/schools/{schoolID}
type SchoolPublic struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Type        SchoolType `json:"type"`
	Location    string     `json:"location,omitempty"`
	Description string     `json:"description,omitempty"`
}

// Term is an academic term (period the school operates classes over).
type Term struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// StartDate and EndDate are calendar dates in "YYYY-MM-DD" form.
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// SchoolClass is a class (group of enrolled students) within a school.
type SchoolClass struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// TermID is the term this class runs in.
	TermID string `json:"termID"`
	// TeacherUID is the UID of the teacher who owns the class.
	TeacherUID string `json:"teacherUID,omitempty"`
	// Room is a room or venue label, e.g. "Room 2B".
	Room string `json:"room,omitempty"`
}

// CreateClassRequest creates a class within a school.
// POST /v0/api4schoolus/spaces/{spaceID}/schools/{schoolID}/classes
type CreateClassRequest struct {
	Title      string `json:"title"`
	TermID     string `json:"termID"`
	TeacherUID string `json:"teacherUID,omitempty"`
	Room       string `json:"room,omitempty"`
}

// CreateClassResponse identifies the new class.
type CreateClassResponse struct {
	ID string `json:"id"`
}

// EnrolStudentRequest enrols a student into a class.
// POST /v0/api4schoolus/spaces/{spaceID}/classes/{classID}/enrolments
// The server stamps EnrolledAt and defaults Status to "active".
type EnrolStudentRequest struct {
	// StudentUID is the UID of the student being enrolled.
	StudentUID string `json:"studentUID"`
	// StudentName is the display name captured at enrolment time.
	StudentName string `json:"studentName"`
	// GuardianUID is the UID of the enrolling guardian, if the student is a minor.
	GuardianUID string `json:"guardianUID,omitempty"`
}

// Enrolment is a student's enrolment record in a class.
type Enrolment struct {
	ID          string          `json:"id"`
	ClassID     string          `json:"classID"`
	StudentUID  string          `json:"studentUID"`
	StudentName string          `json:"studentName"`
	GuardianUID string          `json:"guardianUID,omitempty"`
	Status      EnrolmentStatus `json:"status"`
	EnrolledAt  time.Time       `json:"enrolledAt"`
}

// EnrolStudentResponse identifies the new enrolment.
type EnrolStudentResponse struct {
	ID     string          `json:"id"`
	Status EnrolmentStatus `json:"status"`
}

// TimetableEntry is one scheduled slot in a class timetable.
// GET /v0/api4schoolus/spaces/{spaceID}/classes/{classID}/timetable
type TimetableEntry struct {
	ID      string `json:"id"`
	ClassID string `json:"classID"`
	// Subject is the subject or lesson label, e.g. "Algebra".
	Subject string `json:"subject"`
	// Weekday is the day of week, 0=Sunday .. 6=Saturday.
	Weekday int `json:"weekday"`
	// StartTime and EndTime are local times in "HH:MM" form.
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Room      string `json:"room,omitempty"`
}

// MarkAttendanceRequest marks attendance for one student on a date.
// POST /v0/api4schoolus/spaces/{spaceID}/classes/{classID}/attendance
// The server stamps MarkedByUID (caller uid) and MarkedAt.
type MarkAttendanceRequest struct {
	StudentUID string `json:"studentUID"`
	// Date is the calendar date the mark applies to, "YYYY-MM-DD".
	Date   string           `json:"date"`
	Status AttendanceStatus `json:"status"`
	// Note is an optional free-text note, e.g. reason for an excused absence.
	Note string `json:"note,omitempty"`
}

// AttendanceRecord is a recorded attendance mark.
type AttendanceRecord struct {
	ID         string           `json:"id"`
	ClassID    string           `json:"classID"`
	StudentUID string           `json:"studentUID"`
	Date       string           `json:"date"`
	Status     AttendanceStatus `json:"status"`
	Note       string           `json:"note,omitempty"`
	// MarkedByUID is the UID of the person (teacher) who recorded the mark.
	MarkedByUID string    `json:"markedByUID"`
	MarkedAt    time.Time `json:"markedAt"`
}
