#!/bin/bash

source ../../config.sh

echo -e "\nDeleting server certificates...\n"

sudo rm -rf ${LEGO_DIR}/.lego


echo -e "\nOrdering LEGO to issue a post-quantum certificate"


cd ${LEGO_DIR}

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}@${SERVER_NAME} \
--http.port ":5002" \
--http \
-a \
--pqtls \
--kex Kyber512 \
-k Dilithium2 \
--certalgo Dilithium2 \
--timingcsv ${PQCACME_TESTS_DIR}/measurements/lego_issuance_time.csv \
run
