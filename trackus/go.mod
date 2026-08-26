// Contract lib for the Trackus extension: value types and ports only, no
// persistence, no facades, no transports. Its only dependency is the
// dependency-free action-specification module, so any host can consume it
// cheaply — keep it that way.
module github.com/sneat-co/sneat-ext-contracts/trackus

go 1.26.0

toolchain go1.27.0

require github.com/sneat-co/sneat-go-core/convospec v0.1.0
