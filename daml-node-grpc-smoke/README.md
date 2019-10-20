# Overview

A simple node command line app to verify that the version of node-grpc binding is working and that you are confident that it will work with sextant-api.

## Pre-requisite

1. Install daml sdk, please see doc [https://docs.daml.com/getting-started/installation.html](https://docs.daml.com/getting-started/installation.html)

2. Install node versions between 10.16.2 and 11.14.0. Do not install node v12

## Test against sandbox

1. Navigate `cd` into the folder `<location of sextant-dev project>/daml-node-grpc-smoke`

2. Open a terminal and setup the sandbox `ops.sh sandbox`. You will see something similar to this running continuously in the terminal:

```shell
WARNING: Using an outdated version of the DAML SDK in project.
To migrate to the latest DAML SDK, please set the sdk-version
field in daml.yaml to 0.13.31

Compiling daml-node-ledger-api to a DAR.
Created /Users/paulsitoh/workspace/sextant-project/sextant-dev/daml-node-grpc-smoke/dist/daml-node-ledger-api.dar.
WARNING: Using an outdated version of the DAML SDK in project.
To migrate to the latest DAML SDK, please set the sdk-version
field in daml.yaml to 0.13.31

Sandbox verbosity changed to INFO
DAML LF Engine supports LF versions: 0, 0.dev, 1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.dev; Transaction versions: 1, 2, 3, 4, 5, 6, 7, 8; Value versions: 1, 2, 3, 4, 5
Starting plainText server
listening on localhost:6865
   ____             ____
  / __/__ ____  ___/ / /  ___ __ __
 _\ \/ _ `/ _ \/ _  / _ \/ _ \\ \ /
/___/\_,_/_//_/\_,_/_.__/\___/_\_\

Initialized sandbox version 100.13.29 with ledger-id = sandbox-28b1e775-0431-43d5-b3cc-c509b38e71ad, port = 6865, dar file = List(), time mode = Static, ledger = in-memory, daml-engine = {}
```

3. Set the following environment variables:

```shell
   export ENDPOINT_URL=localhost
   export ENDPOINT_PORT=6685
```

4. Open another terminal and run the test scripts `ops.sh test`

5. If you see the following result, it means the grpc binding is working correctly:

```shell
TAP version 13
# get client and it should return an id
ok 1 should be equal
# get client id should return an id associated with the ledger
ok 2 should be equal
# allocating new parties should not throw error
ok 3 should be equivalent
# allocating same parties twice it should throw an Error number
ok 4 should be truthy
# get list of parties should return an array of parties details
ok 5 should be truthy
# Upload Dar
Before upload --> []
After upload ---> [{"packageId":"b2b6909e12fcc8581f0e2ae94917fd79c221a7e6ad728b01268e0b6664abd810","packageSize":21614,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"1d739910b47df63ee080c1a1df6e5f983f9c5a0573fc0a7c2c20d7592b96cb8d","packageSize":913231,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"0299c0c9da6616c4d2212969050c0c4c0c3085cac657aa3cf3563a37b7522cdb","packageSize":1021661,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}}]

1..5
# tests 5
# pass  5

# ok
```

6. Return to the sandbox terminal and shut it down to terminate the sandbox.

**NOTE:** If you rerun the test before shutting down the sandbox you will see something similar to this:

```shell
Before upload --> [{"packageId":"b2b6909e12fcc8581f0e2ae94917fd79c221a7e6ad728b01268e0b6664abd810","packageSize":21614,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"1d739910b47df63ee080c1a1df6e5f983f9c5a0573fc0a7c2c20d7592b96cb8d","packageSize":913231,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"0299c0c9da6616c4d2212969050c0c4c0c3085cac657aa3cf3563a37b7522cdb","packageSize":1021661,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}}]
After upload ---> [{"packageId":"b2b6909e12fcc8581f0e2ae94917fd79c221a7e6ad728b01268e0b6664abd810","packageSize":21614,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"1d739910b47df63ee080c1a1df6e5f983f9c5a0573fc0a7c2c20d7592b96cb8d","packageSize":913231,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}},{"packageId":"0299c0c9da6616c4d2212969050c0c4c0c3085cac657aa3cf3563a37b7522cdb","packageSize":1021661,"sourceDescription":"","knownSince":{"wrappers_":null,"arrayIndexOffset_":-1,"array":[],"pivot_":1.7976931348623157e+308,"convertedFloatingPointFields_":{}}}]
```

In the sandbox terminal you will see something similar to this:

```shell
Ignoring duplicate upload of package b2b6909e12fcc8581f0e2ae94917fd79c221a7e6ad728b01268e0b6664abd810. Existing package: PackageDetails(21614,1970-01-01T00:00:00Z,None), new package: PackageDetails(21614,1970-01-01T00:00:00Z,None)
Ignoring duplicate upload of package 1d739910b47df63ee080c1a1df6e5f983f9c5a0573fc0a7c2c20d7592b96cb8d. Existing package: PackageDetails(913231,1970-01-01T00:00:00Z,None), new package: PackageDetails(913231,1970-01-01T00:00:00Z,None)
Ignoring duplicate upload of package 0299c0c9da6616c4d2212969050c0c4c0c3085cac657aa3cf3563a37b7522cdb. Existing package: PackageDetails(1021661,1970-01-01T00:00:00Z,None), new package: PackageDetails(1021661,1970-01-01T00:00:00Z,None)
```

In implies that the GRPC node is working properly and the response suggest you have uploaded a duplicate daml archive.

## Test against your daml-on-sawtooth deployment

1. Navigate `cd` into the folder `<location of sextant-dev project>/daml-node-grpc-smoke`

2. Set the following environment variables:

```shell
   export ENDPOINT_URL=<URL to your deployment>
   export ENDPOINT_PORT=<port of your deployment>
```

You will need to configure your deployment to faciliate access via the URL/PORT mentioned above. Please refer to the appropriate documentation.

3. Open a terminal and run the command `Ops.sh test`
