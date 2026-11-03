package briefs4contactus

import (
	"crypto/rand"
	"encoding/hex"
)

// RotateAuthorityGenerations fences inherited access across membership/binding
// lifecycle changes, including delete/recreate ABA. Persist BOTH returned fields
// in the same transaction as the authority mutation. Empty legacy epochs are
// deliberately not usable evidence for inherited access.
func (v *ContactBase) RotateAuthorityGenerations() {
	var tokens [32]byte
	// On supported Go, crypto/rand.Read fills the buffer and never returns
	// an error; entropy-source failure terminates the process.
	_, _ = rand.Read(tokens[:])
	v.MembershipGeneration = hex.EncodeToString(tokens[:16])
	v.BindingGeneration = hex.EncodeToString(tokens[16:])
}
