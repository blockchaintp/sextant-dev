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
    "sandbox")
        pushd ./node-ledger
            npm run sandbox
        popd
        ;;
    "daml-rpc")
        echo "To-do"
        ;;
    *)
        echo "$0 setup | sandbox | daml-rpc"
        ;;
esac


exit 0