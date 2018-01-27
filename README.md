powersoftau
===========

powersoftau is an independent implementation of the [Powers of Tau](https://z.cash.foundation/blog/powers-of-tau/) MPC ceremony.

It is written in Go, shares no code with the [main Rust implementation](https://github.com/ebfull/powersoftau), and uses the [RELIC](https://github.com/relic-toolkit/relic) toolkit for BLS12-381.

Installation
------------

You will need a recent version of [Go](https://golang.org) (1.9+ is known to work), [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) set and a C compiler.

```
git clone --recursive https://github.com/FiloSottile/powersoftau $GOPATH/src/github.com/FiloSottile/powersoftau
cd $GOPATH/src/github.com/FiloSottile/powersoftau && make
go install github.com/FiloSottile/powersoftau/cmd/taucompute
```

Usage
-----

```
Usage of $GOPATH/bin/taucompute:
  -challenge string
    	path to the challenge file (default "./challenge")
  -next string
    	path to the next challenge file, optional
  -pprof
    	run a profiling server; use ONLY FOR DEBUGGING
  -response string
    	path to the response file (default "./response")
```

To facilitate running multiple implementations of Powers of Tau, you can run `taucompute` with the `-next` flag, which will also write a new challenge file once done with the computation.

NOTE: you will have to submit both response files.
