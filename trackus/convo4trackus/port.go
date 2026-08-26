// Package convo4trackus declares the Trackus extension's conversational
// capability: the typed actions it understands, the vocabulary it is listening
// for, the deterministic rules that interpret it, and the port a host calls to
// execute them.
//
// Everything here is persistence-free. The port is satisfied in the Trackus
// implementation repo (trackus/backend/botapi) over facade4trackus; a host binds
// that implementation to these declarations. Nothing in this package may import
// DALgo, a facade, or a transport.
package convo4trackus

import "context"

// TrackByKindContact is the tracking dimension a conversational client always
// uses: a measurement belongs to the person who reported it.
//
// facade4trackus.AddEntryRequest.Validate rejects an empty TrackByKind or
// TrackByID, so the port implementation MUST supply both. Passing the acting
// user's ID as TrackByID lets the facade resolve it to their Contactus contact,
// which is what allows a measurement to land in a space that has never used
// Trackus without any setup step.
const TrackByKindContact = "contact"

// Value is a tracker measurement. Exactly one of Int or Float carries the
// number, matching the tracker's declared value type: a rep count is an
// integer, a distance or a weight is a float. Keeping them separate mirrors the
// upstream entry model instead of forcing every measurement through float64.
type Value struct {
	Int     int
	Float   float64
	IsFloat bool
}

// IntValue returns v as an integer measurement.
func IntValue(n int) Value { return Value{Int: n} }

// FloatValue returns v as a fractional measurement.
func FloatValue(f float64) Value { return Value{Float: f, IsFloat: true} }

// AddEntryRequest records one measurement against one tracker.
type AddEntryRequest struct {
	SpaceID   string
	TrackerID string
	Value     Value
	// Date is an optional ISO date (YYYY-MM-DD) for a backdated measurement.
	// Empty means now.
	//
	// Date sets the entry's reported timestamp, not its identity. EntryID owns
	// retry identity independently of the reported date.
	Date string
	// EntryID is an optional caller-supplied idempotency key. Reusing the same
	// non-empty ID for the same tracker and tracked person must overwrite the
	// same entry rather than create a duplicate. Empty asks the implementation
	// to generate a collision-resistant ID.
	EntryID string
}

// EntryResult describes a recorded measurement.
type EntryResult struct {
	TrackerID    string
	TrackerTitle string
	// Unit is the tracker's unit of measure, if it declares one.
	Unit string
	// Total is the tracker's running total when it is a summable tracker
	// (press-ups, kilometres); zero for absolute trackers such as weight.
	Total float64
	// IsSummable reports whether Total is meaningful.
	IsSummable bool
	// EntryID is the exact stored entry identifier. A conforming implementation
	// returns the caller-supplied ID unchanged or the generated ID when the
	// request omitted it. It may be empty only when talking to a legacy
	// implementation that predates this field.
	EntryID string
}

// ValueType enumerates the measurement shapes a new tracker can take. It is a
// deliberately narrow subset of the upstream tracker model — the shapes a
// conversation can create without a configuration screen.
type ValueType string

const (
	ValueTypeInt   ValueType = "int"
	ValueTypeFloat ValueType = "float"
)

// NumberKind says whether measurements accumulate (repetitions, distance) or
// replace one another (body weight).
type NumberKind string

const (
	NumberKindSummable NumberKind = "summable"
	NumberKindAbsolute NumberKind = "absolute"
)

// CreateTrackerRequest creates a tracker the standard set does not cover, such
// as running distance. It is only ever issued after explicit user confirmation:
// a typo must not permanently add a tracker to a space.
type CreateTrackerRequest struct {
	SpaceID string
	// Title is the human name, e.g. "Running".
	Title string
	// Unit is the unit of measure, e.g. "km". Optional.
	Unit       string
	ValueType  ValueType
	NumberKind NumberKind
}

// TrackerResult identifies a created or fetched tracker.
type TrackerResult struct {
	TrackerID string
	Title     string
	Unit      string
}

// Tracker is the presentation-safe view of a tracker and its recent entries.
type Tracker struct {
	TrackerID string
	Title     string
	Unit      string
	Entries   []Entry
}

// Entry is one recorded measurement in a tracker view.
type Entry struct {
	// When is a formatted timestamp; the application service formats it, so a
	// conversational client never does date presentation itself.
	When  string
	Value Value
}

// Port is the Trackus capability a host binds to the declared actions. It is
// implemented in the Trackus repo over facade4trackus and called only from
// conversational handlers.
type Port interface {
	// AddEntry records a measurement. For a STANDARD tracker ID the upstream
	// facade materialises the tracker on first use, so this succeeds in a space
	// with no Trackus data. For an unknown ID it returns an error for which
	// IsTrackerNotFound reports true, which is the runtime's signal to offer
	// tracker creation instead of surfacing a raw error.
	AddEntry(ctx context.Context, userID string, request AddEntryRequest) (EntryResult, error)

	// CreateTracker creates a non-standard tracker and returns its generated ID.
	CreateTracker(ctx context.Context, userID string, request CreateTrackerRequest) (TrackerResult, error)

	// GetTracker returns a tracker and its recent entries.
	GetTracker(ctx context.Context, userID, spaceID, trackerID string) (Tracker, error)

	// IsTrackerNotFound reports whether err means "that tracker does not exist",
	// so the runtime can branch without matching on error strings.
	IsTrackerNotFound(err error) bool
}
