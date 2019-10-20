#!/bin/bash
export SDK_VERSION=0.13.29
export COMPOSE_PROJECT_NAME=daml

set -e

COMMAND=$1
SERVICE="daml-playground"

case $COMMAND in
    "sandbox")
        ./daml-playground.sh
        ;;
    "test")
        if [ -z $ENDPOINT_URL ]; then
            echo "export URL to either daml sandbox or daml-on-sawtooth deployment"
            echo "If you are using daml sandbox the url is `localhost` or export ENDPOINT_URL=localhost"
        fi
        if [ -z $ENDPOINT_PORT ]; then
            echo "export PORT to either daml sandbox or daml-on-sawtooth deployment"
            echo "If you are using daml sandbox port is `6865` or export ENDPOINT_URL=6865"
        fi
        pushd ./node-ledger
            npm run sandbox
        popd
        ;;
    *)
        echo "$0 [sandbox | test]"
        ;;
esac


exit 0