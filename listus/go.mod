// Contract lib for the Listus extension: value types and ports only, no
// persistence, no facades, no transports. Its only dependency is the
// dependency-free action-specification module, so any host can consume it
// cheaply — keep it that way.
module github.com/sneat-co/sneat-ext-contracts/listus

go 1.26.0

toolchain go1.27.0

require (
	github.com/dal-go/dalgo v0.79.4
	github.com/sneat-co/sneat-go-core v0.69.0
	github.com/sneat-co/sneat-go-core/convospec v0.1.2
)

require (
	github.com/RoaringBitmap/roaring/v2 v2.27.0 // indirect
	github.com/bits-and-blooms/bitset v1.24.6 // indirect
	github.com/dal-go/record v0.1.3 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/validation v0.0.13 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
