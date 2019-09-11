#!/bin/bash

COMMAND=$1
CLUSTER_NAME=$2
REGION=$3

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
    echo "Missing cluster name"
    exit 1
  fi
}

function checkRegion() {

  if [ -z $REGION ]; then
    echo "Missing region name"
    exit 1
  fi

}

function create() {
  eksctl create cluster --name $CLUSTER_NAME --node-type m5.large --nodes 4 --region $REGION
}

function delete() {
  eksctl delete cluster --name $CLUSTER_NAME -r $REGION -w
}

prequisite
case "$COMMAND" in
    create)
        checkClusterName
        checkRegion
        create
        ;;
    delete)
        checkClusterName
        checkRegion
        delete
        ;;
    *)
        echo "$0 (create | delete) name region"
        ;;
esac
