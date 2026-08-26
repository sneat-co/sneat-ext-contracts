---
listus-contract: patch
---

feat(listus-contract): add listus Go module extracted from ext-listus backend/v0.0.8

Adds the listus Go module (github.com/sneat-co/sneat-ext-contracts/listus) extracted
from sneat-co/ext-listus, `backend/` @ tag `backend/v0.0.8`
(1ad9bfa1cf995a116cd78de2a09741ee05c1f6e0). Lockstep npm+Go release: resulting version
0.0.6 (npm listus-contract 0.0.5 -> 0.0.6, first Go tag listus/v0.0.6). Consumers move
from the old paths — backend/v0.0.8 (sneat-bots, sneat-go/sneat-cli) and backend/v0.0.5
(the listus product repo) — to listus/v0.0.6: a module-path restart, founder-accepted.
