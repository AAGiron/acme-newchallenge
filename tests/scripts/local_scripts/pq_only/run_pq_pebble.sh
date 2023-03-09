#!/bin/bash
source ../../config.sh

cd ${PEBBLE_DIR}

go run cmd/pebble/main.go \
-pqtls \
-kex Kyber512 \
-rootdir ${WRAPPED_CERT_TESTS_DIR}/root_ca \
-timingcsv ${WRAPPED_CERT_TESTS_DIR}/measurements/pebble_issuance_time.csv \
-ocspresponsepath ${WRAPPED_CERT_TESTS_DIR}/resources/ocsp_response \
-rootSig Dilithium2 \
-interSig Dilithium2 \
-issuerSig Dilithium2