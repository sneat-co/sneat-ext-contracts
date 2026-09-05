# Debtus Go contract

This module defines the storage-neutral public boundary for posting source
obligations to Debtus and reading their durable status. Debtus owns the
implementation and financial records; callers such as Splitus own their source
revisions and posting intents.

Trusted authority reservations are deliberately absent from serialized request
types. Server composition must bind and validate that evidence separately, and
the provider must match recorder/actor IDs to its authenticated server context;
the serialized IDs grant no authority.

The immutable reconciliation receipt is separate from mutable source status.
Current status exposes Debtus-owned principal, outstanding, repayment and credit
amounts without prescribing how a credit is refunded or settled.

Reconciliation returns explicit pending, posting, applied or attention state;
only applied state carries an immutable receipt. A separate paged activity query
keeps repayments and adjustments visible and links them to stable root activity
and source-line IDs without exposing Debtus storage records.
