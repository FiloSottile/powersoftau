# Introduction

This Gitian descriptor allows deterministic builds of the powersoftau Go implementation.

# HOWTO

```
git clone -b vbox git@github.com:devrandom/gitian-builder
git clone git@github.com:FiloSottile/powersoftau golang-powersoftau
cd gitian-builder
```

You can use Gitian in either VirtualBox or LXC mode. KVM doesn't yet work with a Debian guest.

## LXC

```
bin/make-base-vm --lxc --suite stretch --distro debian
export USE_LXC=1
```

Notes:

- this doesn't work in Ubuntu 17.10 due to a bug in `lxc-execute` - see https://github.com/lxc/lxc/issues/2028.
- br0 must be set up as described in the gitian-builder README

## VirtualBox

```
bin/make-base-vm --vbox --suite stretch --distro debian
export USE_VBOX=1
```

## Building

```
bin/gbuild ../golang-powersoftau/contrib/gitian/powersoftau.yml
```

