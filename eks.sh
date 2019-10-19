#!/bin/bash

COMMAND=$1

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
    echo "set env CLUSTER_NAME"
    exit 1
  fi
}

function checkRegion() {

  if [ -z $AWS_REGION ]; then
    echo "Missing region name"
    echo "set env AWS_REGION"
    exit 1
  fi

}

function checkDockerCred() {

  if [ -z $DOCKER_USER ]; then
    echo "Missing docker user name"
    echo "set env DOCKER_USER"
    exit 1
  fi

  if [ -z $DOCKER_PW ]; then
    echo "Missing docker password"
    echo "set env DOCKER_PW"
    exit 1
  fi

  if [ -z $DOCKER_EMAIL ]; then
    echo "Missing docker email"
    echo "set env DOCKER_EMAIL"
    exit 1
  fi
}

function createCred() {
  kubectl create secret docker-registry regcred --docker-server=https://dev.catenasys.com:8083/ --docker-username=$DOCKER_USER --docker-password='${DOCKER_PW}' --docker-email=$DOCKER_EMAIL
}

prequisite
case "$COMMAND" in
    create-cred)
        checkDockerCred
        createCred
        ;;
    url)
        aws eks describe-cluster --name $CLUSTER_NAME --region $AWS_REGION --output json | jq '.cluster.endpoint'
        ;;
    *)
        echo "$0 [ create-cred | url ]"
        ;;
esac