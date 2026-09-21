# zktf-sdk-go

A Go SDK for interacting with the zktf network (server-side).

## Simulator

`simulator` runs a simulated network, device and verifier in process, for tests
that would otherwise need a deployed backend.

It links `libzktf_sim` in addition to `libzktf_sdk`. Because cgo link flags are
per-package, that only applies to binaries that actually import
`zktf-sdk-go/simulator` — depending on the SDK alone needs neither the
simulator package nor its native library.

`scripts/fetch-native.sh` downloads both archives and writes the cgo
environment to `.env`.
