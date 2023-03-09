#!/bin/bash

source ../../config.sh

echo -e "\nDeleting server certificates...\n"

sudo rm -rf ${LEGO_DIR}/.lego/*

cd ${LEGO_DIR}

echo -e "Ordering LEGO to issue a post-quantum certificate\n"

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}@${SERVER_NAME} \
--http.port ":5002" \
--http \
--certlabel pq_cert \
-a \
--pqtls \
--kex Kyber512 \
-k Dilithium2 \
--certalgo Dilithium3 \
--timingcsv ${WRAPPED_CERT_TESTS_DIR}/measurements/lego_issuance_time.csv \
--serverip ${SERVER_NAME} \
run

cd ${WRAPPED_CERT_TESTS_DIR}/src

go run launch_servers.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-pqtls \
-cert ${LEGO_DIR}/.lego/certificates/pq_cert/${SERVER_NAME}.crt \
-key ${LEGO_DIR}/.lego/certificates/pq_cert/${SERVER_NAME}.key \
-ocspresponsepath ${WRAPPED_CERT_TESTS_DIR}/resources/ocsp_response \
-measurements ${WRAPPED_CERT_TESTS_DIR}/measurements

