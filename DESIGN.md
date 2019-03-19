# Making tcp-info and paris-traceroute data joinable with experiment data

author       | Peter Boothe <pboothe@google.com>
-------------+----------------
last updated | 2019-03-19
status       | approved

## Introduction

We want to join data derived from tcp-info, from paris-traceroute, and from
the original experiments. We need the solution to be unambiguous, because
fuzzy matching via timestamps represents a usability challenge for M-Lab data
users. We would like the solution to be simple, because complex is harder to
maintain and simpler makes our data more accessible to the world and
ourselves.

## Background

Joining M-Lab's experiment data with our background-generated per-connection
data is currently extremely difficult. Users are forced to perform fuzzy
matching based on network 5-tuples and the timestamp of a connection. This
mostly works, but represents a pretty high complexity bar that a researcher
needs to clear to even begin deriving value out of combining our data
streams. Furthermore, when the fuzzy matching fails, it is never clear
whether the matching failed due to the fuzzy-match being overly-restrictive
or whether the data is simply missing. By providing a unique universal
identifier string for every socket on the platform, we can eliminate both the
computational difficulty of the fuzzy matching and eliminate any questions
about whether matching data exists or not.

## Solution

Generate a globally-unique UUID for every TCP socket and record that UUID
with every measurement. We suggest using the UUID in the filename of the
results.

From every TCP socket connection we can generate a unique string, consisting
of all of:

- server hostname
- server boot time (in seconds-since-epoch, expressed as a decimal rounded to the nearest second)
- socket TCP cookie (also in hex, a fixed-digit hexadecimal encoding of a 64 bit number)

all joined with underscores. e.g.: `ndt.mlab3.atl05.measurement-lab.org_1548788619_00000000000084FF`

We need the TCP cookie because it is a socket UUID for a given server
instance. We need the boot time and the server hostname, because those
uniquely specify the server instance. A hostname by itself is not enough - it
appears that the TCP cookie resets to zero at boot, so we also need to know
when the server booted up to uniquely specify the namespace for the cookie.
Together, those uniquely specify a socket, and in the spirit of not including
extra information we don't need, we call that good enough and only use those
three items.

An experiment SHOULD generate this string and save it as metadata for all
connected TCP sockets that it uses for measurement. The easiest form of
metadata to use is the filename, so consider using the generated string as a
filename, as long as you can be confident that the hostname is safe to use as
a filename.

M-Lab has written a library which does this mapping of socket to string. We
recommend (but do not require) that everyone use that library instead of
rolling their own solution.

Because of how we are constructing the UUID, it has a nice internal structure
which MAY be taken advantage of. However, that structure IS NOT guaranteed.
The only thing we guarantee about our UUIDs is their uniqueness, not their
internal structure. Production code SHOULD NOT depend on any aspect of the
UUID's structure besides its uniqueness.

### Server Hostname

This should be equivalent to the output of a correctly-functioning hostname
command.

### Server Boot Time

Unfortunately, there is no unambiguous way to convert the floating-point
difference between `/proc/uptime` and `time.Now()` into an integer.
Therefore, we recommend that implementations do the conversion once,
write the result to a well-known location, and that all libraries read from
that location. BSD users may, at this point, casually note that they have
`sysctl kern.boottime` and do not suffer from this limitation in Linux.

### Socket TCP Cookie

TCP cookies were added to the Linux kernel in 2017. They can be discovered in
C using `getsockopt` with the appropriate options or via the `ss -e` command.

## Library implementation

A Go-based implementation can be found in <https://github.com/m-lab/uuid/>

## Example integrations

- [tcp-info](https://github.com/m-lab/tcp-info)
- [ndt-server](https://github.com/m-lab/ndt-server)
- [traceroute-caller](https://github.com/m-lab/traceroute-caller)