#!/bin/bash

COMMAND=$1
CLUSTER_NAME=$2

function prequisite() {
  aws --version
  if [ "$?" == "127" ]; then
    exit "aws cli not installed. Please install to proceed."
  fi

  eksctl version
  if [ "$?" == "124" ]; then
    exit "eksctl cli is not installed. Please install to proceed."
  fi
}

function checkClusterName() {

  if [ -z $CLUSTER_NAME ]; then
    echo "Missing cluter name"
    exit 1
  fi
}

function create() {
  eksctl create cluster --name $CLUSTER_NAME --node-type m5.large --nodes 4  
}

function delete() {
  eksctl delete cluster --name $CLUSTER_NAME
}

prequisite
case "$COMMAND" in
    create)
        checkClusterName
        create
        ;;
    delete)
        checkClusterName
        delete
        ;;
    *)
        echo "$0 (create | delete) name"
        ;;
esac