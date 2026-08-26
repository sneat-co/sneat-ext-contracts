// Package facade4contactus declares the contactus contributor interfaces that
// pass the extension-contract-repo ownership test — interfaces whose entire
// signature is expressible in contactus-own (contactusmodels) plus
// foundational/core types, and therefore belong in the published contactus
// contract.
//
// These are caller-satisfied callback signatures: the contactus module provides
// and registers the concrete implementation at bootstrap, while sibling modules
// (spaceus, userus) consume the interface from here so they need not depend on
// contactus DAL/DBO types directly. The registration plumbing and the consumer
// import-path repoint land during the ordered consumer cutover, after this
// contract is published.
//
// Interfaces that fail the ownership test stay in their consumer module:
// invitus's ContactusAccess (which speaks invitus/spaceus types) and the
// linkage RelatedDboFactory (which speaks dal4spaceus/dbo4linkage types) are
// deliberately NOT declared here, to preserve the zero-other-extension-deps
// invariant.
package facade4contactus
