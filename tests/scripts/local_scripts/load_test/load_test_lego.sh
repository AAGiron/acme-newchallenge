#!/bin/bash

source ../../config.sh

echo -e "Run Pebble!\n"
read -p "Is Pebble already running ? (y/n) " PEBBLE_RUNNING

if [ "$PEBBLE_RUNNING" == "n" ]; then
    exit
fi 

sudo rm -rf ${LEGO_DIR}/.lego
cd ${LEGO_DIR}

# Classic scenario

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
-certalgo P256 \
--loadtestfinalize \
--numthreads 10 \
--loadtestduration 10 \
--loadtestcsv ${PQCACME_TESTS_DIR}/measurements/lego_load_tests_result.csv \
run

# Post-quantum scenario

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}-newchallenge@${SERVER_NAME} \
--http.port ":5002" \
--http \
-a \
-newchallenge \
--pqorderport 10004 \
--pqtls \
--kex Kyber512 \
-k Dilithium2 \
--certalgo Dilithium2 \
--loadtestfinalize \
--numthreads 10 \
--loadtestduration 10 \
--loadtestcsv ${PQCACME_TESTS_DIR}/measurements/lego_load_tests_result.csv \
run