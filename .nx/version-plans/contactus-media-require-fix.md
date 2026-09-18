---
contactus-contract: patch
---

Require the media contract tag that exists.

contactus/go.mod required `github.com/sneat-co/sneat-ext-contracts/media
v0.1.0`, a sibling tag that was never cut — the media module's only release is
`media/v0.1.1`. `contactus/v0.12.10` shipped carrying that requirement, so every
consumer upgrading to it failed module resolution before compiling:

    go: github.com/sneat-co/sneat-ext-contracts/media@v0.1.0:
        invalid version: unknown revision media/v0.1.0

The requirement was corrected in 72785b5e, but that merge carried no version
plan, so `Prepare release` stopped at `has=false` and no tag was cut. Consumers
still have no resolvable version above 0.12.9 until this releases.

Resulting version: 0.12.11 (patch bump from declared 0.12.10).
