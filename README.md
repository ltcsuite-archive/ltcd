ltcd
====

[![Build Status](https://travis-ci.org/ltcsuite/ltcd.png?branch=master)](https://travis-ci.org/ltcsuite/ltcd)

ltcd is an alternative full node litecoin implementation written in Go (golang).

This project is currently under active development and is in a Beta state.  It
is extremely stable and has been in production use since October 2013.

It properly downloads, validates, and serves the block chain using the exact
rules (including consensus bugs) for block acceptance as Litecoin Core.  We have
taken great care to avoid ltcd causing a fork to the block chain.  It includes a
full block validation testing framework which contains all of the 'official'
block acceptance tests (and some additional ones) that is run on every pull
request to help ensure it properly follows consensus.  Also, it passes all of
the JSON test data in the Litecoin Core code.

It also properly relays newly mined blocks, maintains a transaction pool, and
relays individual transactions that have not yet made it into a block.  It
ensures all individual transactions admitted to the pool follow the rules
required by the block chain and also includes more strict checks which filter
transactions based on miner requirements ("standard" transactions).

One key difference between ltcd and Litecoin Core is that ltcd does *NOT* include
wallet functionality and this was a very intentional design decision.  See the
blog entry [here](https://blog.conformal.com/ltcd-not-your-moms-bitcoin-daemon)
for more details.  This means you can't actually make or receive payments
directly with ltcd.  That functionality is provided by the
[btcwallet](https://github.com/ltcsuite/btcwallet) and
[Paymetheus](https://github.com/ltcsuite/Paymetheus) (Windows-only) projects
which are both under active development.

## Requirements

[Go](http://golang.org) 1.7 or newer.

## Installation

#### Windows - MSI Available

https://github.com/ltcsuite/ltcd/releases

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Ensure Go was installed properly and is a supported version:

```bash
$ go version
$ go env GOROOT GOPATH
```

NOTE: The `GOROOT` and `GOPATH` above must not be the same path.  It is
recommended that `GOPATH` is set to a directory in your home directory such as
`~/goprojects` to avoid write permission issues.  It is also recommended to add
`$GOPATH/bin` to your `PATH` at this point.

- Run the following commands to obtain ltcd, all dependencies, and install it:

```bash
$ go get -u github.com/Masterminds/glide
$ git clone https://github.com/ltcsuite/ltcd $GOPATH/src/github.com/ltcsuite/ltcd
$ cd $GOPATH/src/github.com/ltcsuite/ltcd
$ glide install
$ go install . ./cmd/...
```

- ltcd (and utilities) will now be installed in ```$GOPATH/bin```.  If you did
  not already add the bin directory to your system path during Go installation,
  we recommend you do so now.

## Updating

#### Windows

Install a newer MSI

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Run the following commands to update ltcd, all dependencies, and install it:

```bash
$ cd $GOPATH/src/github.com/ltcsuite/ltcd
$ git pull && glide install
$ go install . ./cmd/...
```

## Getting Started

ltcd has several configuration options avilable to tweak how it runs, but all
of the basic operations described in the intro section work with zero
configuration.

#### Windows (Installed from MSI)

Launch ltcd from your Start menu.

#### Linux/BSD/POSIX/Source

```bash
$ ./ltcd
```

## IRC

- irc.freenode.net
- channel #litecoin
- [webchat](https://webchat.freenode.net/?channels=litecoin)

## Issue Tracker

The [integrated github issue tracker](https://github.com/ltcsuite/ltcd/issues)
is used for this project.

## Documentation

The documentation is a work-in-progress.  It is located in the [docs](https://github.com/ltcsuite/ltcd/tree/master/docs) folder.

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from the ltcsuite developers.  To
verify the signature perform the following:

- Download the public key from the Conformal website at
  https://opensource.conformal.com/GIT-GPG-KEY-conformal.txt

- Import the public key into your GPG keyring:
  ```bash
  gpg --import GIT-GPG-KEY-conformal.txt
  ```

- Verify the release tag with the following command where `TAG_NAME` is a
  placeholder for the specific tag:
  ```bash
  git tag -v TAG_NAME
  ```

## License

ltcd is licensed under the [copyfree](http://copyfree.org) ISC License.
