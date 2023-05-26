#!/bin/bash
source ../../config.sh

cd ${PEBBLE_DIR}

go run cmd/pebble/main.go \
-pqtls \
-kex Kyber512 \
-newchallenge \
-rootdir ${PQCACME_TESTS_DIR}/root_ca \
-timingcsv ${PQCACME_TESTS_DIR}/measurements/pebble_issuance_time.csv \
-ocspresponsepath ${PQCACME_TESTS_DIR}/resources/ocsp_response \
-rootSig ECDSA-P256 \
-interSig ECDSA-P256 \
-issuerSig ECDSA-P256 \
--pqorderport 10004 \
--pqorderroot Dilithium2 \
--pqorderissuer Dilithium2 \ 

