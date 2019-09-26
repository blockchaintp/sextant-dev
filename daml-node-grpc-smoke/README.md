# Overview

A simple node command line app to invoke daml ledger

# Pre-requisite

1. Install docker
1. Install C++ conpiler (on macOS install Xcode)
1. Install node v11.x.x (ideally 11.14.0). Do not install node v12

# To test against sandbox

1. `cd` into node-ledger
1. Setup the sandbox `ops.sh setup`
1. Run the test scripts `ops.sh test`
1. Teardown the sandbox `ops.sh teardown`
1. You should see the following:

```
TAP version 13
# get client id should return an id associated with the ledger
ok 1 should be equal
# allocating new parties should not throw error
ok 2 should be equivalent
# allocating same parties twice it should throw an Error number
ok 3 should be truthy
# get list of parties should return an array of parties details
ok 4 should be truthy

1..4
# tests 4
# pass  4

# ok
```