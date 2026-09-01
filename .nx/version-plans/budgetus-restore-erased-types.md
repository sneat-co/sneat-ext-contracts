---
budgetus-contract: patch
---

fix(budgetus): restore types erased to `any`

The 0.1.0 -> 0.1.2 migration replaced five published types with `any`:
IListContext's Dbo generic, IBudgetusSpaceDbo.listGroups, the space
parameters of IBudgetusService.deleteList/getListById, and
IListItemResult.listDto. Contract 0.1.0 contained no `any` at all.

Consumers lost type checking on list DBOs and could not compile against
anything newer than 0.1.0.
