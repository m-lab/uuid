[![Travis Build Status](https://travis-ci.org/m-lab/uuid.svg?branch=master)](https://travis-ci.org/m-lab/uuid) | [![Coverage Status](https://coveralls.io/repos/m-lab/uuid/badge.svg?branch=master)](https://coveralls.io/github/m-lab/uuid?branch=master) | [![GoDoc](https://godoc.org/github.com/m-lab/uuid?status.svg)](https://godoc.org/github.com/m-lab/uuid) | [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/uuid)](https://goreportcard.com/report/github.com/m-lab/uuid)

# Socket UUIDs for the M-Lab platform

This allows us to generate a globally unique ID for any TCP socket. When we
say globally, we really mean globally - it should be impossible to have two
machines generate the same UUID.

The only case the uniqueness of the UUID could be violated is if two machines
have the same hostname and booted up at the exact same second in time, but it is
bad practice to give machines the same hostname (so don't).

The design of the UUIDs and this system for creating them can be found in
`[./DESIGN.md]`.