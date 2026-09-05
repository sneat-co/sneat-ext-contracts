# Debtus Go contract

This module defines the storage-neutral public boundary for posting source
obligations to Debtus and reading their durable status. Debtus owns the
implementation and financial records; callers such as Splitus own their source
revisions and posting intents.

Trusted authority reservations are deliberately absent from serialized request
types. Server composition must bind and validate that evidence separately.
