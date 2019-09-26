#!/bin/bash
export SDK_VERSION=0.13.20
export COMPOSE_PROJECT_NAME=daml

set -e

COMMAND=$1
SERVICE="daml-playground"

case $COMMAND in
    "setup")
        docker-compose -f ../ledger-smoke-test.yaml up -d
        docker container ls
        ;;
    "test")
        npm test
        ;;
    "teardown")
        docker rm -f $SERVICE
        docker container ls
        ;;
    *)
        echo "$0 setup | test |  teardown "
        ;;
esac


exit 0