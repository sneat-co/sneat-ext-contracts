# @sneat/extension-splitus-contract

Splitus extension contract library.

## Bill contract version 1

`ICreateSplitusBillV1Request` carries exact major-unit decimal strings,
client-stable bill identity, and separate paid/owed contact allocations. The
host must call `assertCreateSplitusBillV1Recorder` with its trusted
authenticated identity before accepting the recorder audit claim.

Parse untrusted API responses with `parseCreateSplitusBillV1Response`,
`parseGetSplitusBillV1Response`, or `parseListSplitusBillsV1Response`. This
rejects numeric money JSON, invalid posting proof, unsafe identifiers, and
unbounded list/detail arrays before a runtime renders them. Contact IDs are
resolved through the host's real Contactus data; they are not display labels.

The earlier split DTOs and `ISplitusService` remain exported as deprecated
compatibility surfaces until every consumer has moved to
`ISplitusBillServiceV1`.

The matching storage-neutral Go host contract is module
`github.com/sneat-co/sneat-ext-contracts/splitus`, package
`contract4splitus`.

## Provenance

Migrated from sneat-co/debtus (npm continuity from @sneat/extension-splitus-contract@0.2.0).
