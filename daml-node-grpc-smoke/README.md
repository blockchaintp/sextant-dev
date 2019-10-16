# Overview

A simple node command line app to verify that the version of node-grpc binding is working and that you are confident that it will work with sextant-api.

## Pre-requisite

1. Install daml sdk, please see doc [https://docs.daml.com/getting-started/installation.html](https://docs.daml.com/getting-started/installation.html)
1. Install node versions between 10.16.2 and 11.14.0. Do not install node v12

## Test against sandbox

1. Open a terminal and setup the sandbox `ops.sh setup`
1. Open another terminal and run the test scripts `ops.sh test`
1. If you see the following result, it means the grpc binding is working correctly:

```shell
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
