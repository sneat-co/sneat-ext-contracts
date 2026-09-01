---
budgetus-contract: patch
---

fix(budgetus): restore spaceID on list requests

IListRequest and ICreateListRequest lost `extends ISpaceRequest` between
0.1.0 and 0.1.2, dropping spaceID from the whole list-request family even
though the service still sends it and the backend still routes on it.
Consumers could not compile against anything newer than 0.1.0.
