#!/bin/bash
export SDK_VERSION=0.13.29
export COMPOSE_PROJECT_NAME=daml

set -e

COMMAND=$1
SERVICE="daml-playground"

case $COMMAND in
    "setup")
        ./daml-playground.sh
        ;;
    "test")
        pushd ./node-ledger
            npm test
        popd
        ;;
    *)
        echo "$0 setup | test"
        ;;
esac


exit 0