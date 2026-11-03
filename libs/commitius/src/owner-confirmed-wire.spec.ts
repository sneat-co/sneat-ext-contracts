// Generated from the owning Go amendment flow with fictional test data.
// TestSelfReportAmendAtomicallyPreservesHistoryAndReplacesExpiry emits this
// exact command + selected read + history when IBETICAN_OWNER_CHANGE_FIXTURE_PROOF
// is set. Source JSON SHA256: c89ee3c22cfade29bf7882c0f4c31c448afe2f974d2443c0b05c1c80515feaa6
import {describe,expect,expectTypeOf,it} from 'vitest';
import type {IGoalTermsChangeOutcome,IGoalTermsChangeView,IGoalChangeHistory} from './index';

const ownerConfirmedWire = {
  "outcome": {
    "changeID": "owner_chg_e29adcdda71cbc92195d59f61a0c6ac59e516b0eced77da6e01b9b882f64d2c8",
    "baseVersion": 1,
    "proposedVersion": 2,
    "status": "owner-confirmed"
  },
  "selected": {
    "changeID": "owner_chg_e29adcdda71cbc92195d59f61a0c6ac59e516b0eced77da6e01b9b882f64d2c8",
    "baseVersion": 1,
    "proposedVersion": 2,
    "status": "owner-confirmed",
    "requestedAt": "2026-07-18T01:00:00Z",
    "closedAt": "2026-07-18T01:00:00Z",
    "base": {
      "version": 1,
      "effectiveVersion": 2,
      "schemaVersion": 1,
      "supported": true,
      "title": "Meditate ten minutes each day",
      "deadline": "2026-07-20T00:00:00Z",
      "definition": {
        "schemaVersion": 1,
        "conjunction": "all",
        "evaluation": {
          "kind": "every-period",
          "startAt": "2026-07-17T00:00:00Z",
          "timeZone": "UTC",
          "periodUnit": "day",
          "anchor": "local-start-date",
          "partialPeriods": "include",
          "requiredPeriodCount": 3
        },
        "clauses": [
          {
            "id": "meditate",
            "kind": "numeric",
            "metric": "duration",
            "comparator": "gte",
            "target": 600,
            "unitOfMeasure": "second"
          }
        ]
      },
      "calendar": {
        "normalizationVersion": "goal-calendar-v1",
        "kind": "every-period",
        "normalizedStartAt": "2026-07-17T00:00:00Z",
        "normalizedDeadline": "2026-07-20T00:00:00Z",
        "timeZone": "UTC",
        "localStartLabel": "2026-07-17T00:00:00Z",
        "localDeadlineLabel": "2026-07-20T00:00:00Z",
        "periodUnit": "day",
        "anchor": "local-start-date",
        "partialPeriods": "include",
        "requiredPeriodCount": 3
      },
      "calendarDigest": "6fe4448fe9d1068f1fec60cfdcdac97154da77bdfd895694fd24aa84f52d86ba",
      "totalPeriodCount": 3,
      "periodPage": {
        "offset": 0,
        "periods": [
          {
            "id": "p0001",
            "localLabel": "2026-07-17T00:00:00Z — 2026-07-18T00:00:00Z",
            "startAt": "2026-07-17T00:00:00Z",
            "endAt": "2026-07-18T00:00:00Z",
            "partial": false,
            "intervals": [
              {
                "startAt": "2026-07-17T00:00:00Z",
                "endAt": "2026-07-18T00:00:00Z"
              }
            ]
          }
        ],
        "hasMore": true
      }
    },
    "proposed": {
      "version": 2,
      "effectiveVersion": 2,
      "schemaVersion": 1,
      "supported": true,
      "title": "My meaningful revised goal",
      "deadline": "2026-07-20T00:00:00Z",
      "definition": {
        "schemaVersion": 1,
        "conjunction": "all",
        "evaluation": {
          "kind": "every-period",
          "startAt": "2026-07-17T00:00:00Z",
          "timeZone": "UTC",
          "periodUnit": "day",
          "anchor": "local-start-date",
          "partialPeriods": "include",
          "requiredPeriodCount": 3
        },
        "clauses": [
          {
            "id": "meditate",
            "kind": "numeric",
            "metric": "duration",
            "comparator": "gte",
            "target": 600,
            "unitOfMeasure": "second"
          }
        ]
      },
      "calendar": {
        "normalizationVersion": "goal-calendar-v1",
        "kind": "every-period",
        "normalizedStartAt": "2026-07-17T00:00:00Z",
        "normalizedDeadline": "2026-07-20T00:00:00Z",
        "timeZone": "UTC",
        "localStartLabel": "2026-07-17T00:00:00Z",
        "localDeadlineLabel": "2026-07-20T00:00:00Z",
        "periodUnit": "day",
        "anchor": "local-start-date",
        "partialPeriods": "include",
        "requiredPeriodCount": 3
      },
      "calendarDigest": "6fe4448fe9d1068f1fec60cfdcdac97154da77bdfd895694fd24aa84f52d86ba",
      "totalPeriodCount": 3,
      "periodPage": {
        "offset": 0,
        "periods": [
          {
            "id": "p0001",
            "localLabel": "2026-07-17T00:00:00Z — 2026-07-18T00:00:00Z",
            "startAt": "2026-07-17T00:00:00Z",
            "endAt": "2026-07-18T00:00:00Z",
            "partial": false,
            "intervals": [
              {
                "startAt": "2026-07-17T00:00:00Z",
                "endAt": "2026-07-18T00:00:00Z"
              }
            ]
          }
        ],
        "hasMore": true
      }
    },
    "requiredApprovalCount": 0,
    "approvalCount": 0,
    "canApprove": false,
    "canReject": false,
    "canWithdraw": false
  },
  "history": {
    "effectiveVersion": 2,
    "items": [
      {
        "changeID": "owner_chg_e29adcdda71cbc92195d59f61a0c6ac59e516b0eced77da6e01b9b882f64d2c8",
        "baseVersion": 1,
        "proposedVersion": 2,
        "requestedAt": "2026-07-18T01:00:00Z",
        "closedAt": "2026-07-18T01:00:00Z",
        "status": "owner-confirmed",
        "supported": true
      }
    ],
    "hasMore": false
  }
} as const;

describe('actual owner-confirmed wire values',()=>{
 it('admits the owning command, selected read and history without pretending referee consent',()=>{
  expectTypeOf(ownerConfirmedWire.outcome).toExtend<IGoalTermsChangeOutcome>();
  expectTypeOf(ownerConfirmedWire.selected).toExtend<IGoalTermsChangeView>();
  expectTypeOf(ownerConfirmedWire.history).toExtend<IGoalChangeHistory>();
  expect(ownerConfirmedWire.outcome.status).toBe('owner-confirmed');
  expect(ownerConfirmedWire.selected.status).toBe('owner-confirmed');
  expect(ownerConfirmedWire.history.items[0].status).toBe('owner-confirmed');
  expect(ownerConfirmedWire.selected.requiredApprovalCount).toBe(0);
  expect(ownerConfirmedWire.selected.approvalCount).toBe(0);
  expect(ownerConfirmedWire.selected.canApprove).toBe(false);
  expect(ownerConfirmedWire.selected.canReject).toBe(false);
  expect(ownerConfirmedWire.selected.canWithdraw).toBe(false);
 });
});
