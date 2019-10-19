#!/bin/bash

export PATH=$PATH:$HOME/.daml/bin

if [ -d ${PWD}/dist ]; then
    rm -rf ${PWD}/dist
fi

daml build --project-root $PWD -o $PWD/dist/daml-node-ledger-api.dar
daml sandbox