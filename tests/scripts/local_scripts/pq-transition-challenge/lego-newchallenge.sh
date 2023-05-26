#!/bin/bash

source ../../config.sh

echo -e "\nDeleting server certificates...\n"

sudo rm -rf ${LEGO_DIR}/.lego/*

cd ${LEGO_DIR}

echo -e "Ordering LEGO to issue the first certificate"

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}@${SERVER_NAME} \
--http.port ":5002" \
--http \
-a \
--pqtls \
--kex Kyber512 \
-k ec256 \
--certalgo P256 \
--timingcsv ${PQCACME_TESTS_DIR}/measurements/lego_issuance_time.csv \
run

echo -e "\nOrdering LEGO to issue the second post-quantum certificate"


cd ${LEGO_DIR}

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}-newchallenge@${SERVER_NAME} \
--http.port ":5002" \
--http \
-a \
-newchallenge \
--pqorderport 10003 \
--pqtls \
--kex Kyber512 \
-k Dilithium2 \
--certalgo Dilithium2 \
--timingcsv ${PQCACME_TESTS_DIR}/measurements/lego_issuance_time.csv \
run
