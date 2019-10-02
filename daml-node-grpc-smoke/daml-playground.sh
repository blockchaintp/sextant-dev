#!/bin/bash

export PATH=$PATH:$HOME/.daml/bin

daml build
daml sandbox $PWD/.daml/dist/daml-node-ledger-api-1.0.0.dar
