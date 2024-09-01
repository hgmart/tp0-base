#!/bin/bash

IMAGE="alpine:latest" 
CONTAINER_NAME="test_echo_server"
DOCKER_NETWORK="tp0_testing_net"
SERVER_ADDRESS="server:12345"
SENT_MESSAGE="test_message"
SUCCESS_MESSAGE="action: test_echo_server | result: success"
FAILURE_MESSAGE="action: test_echo_server | result: fail"

# El siguiente código envía a través de la red DOCKER_NETWORK definida
# un mensaje que espera de regreso, en caso correcto muestra el mensaje
# SUCCESS_MESSAGE y en caso contrario FAILURE_MESSAGE.
docker run --rm --network $DOCKER_NETWORK --name $CONTAINER_NAME $IMAGE \
    sh -c "echo $SENT_MESSAGE \
    | nc $SERVER_ADDRESS \
    | grep -q $SENT_MESSAGE && (echo \"$SUCCESS_MESSAGE\"; exit 0)  || (echo \"$FAILURE_MESSAGE\"; exit 1)"