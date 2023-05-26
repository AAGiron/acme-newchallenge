#!/bin/bash

source ../../config.sh

cd ${PEBBLE_DIR}

go run cmd/pebble/main.go \
-pqtls \
-kex Kyber512 \
-newchallenge \
-rootdir ${PQCACME_TESTS_DIR}/root_ca \
-rootSig ECDSA-P256 \
-interSig ECDSA-P256 \
-issuerSig ECDSA-P256 \
-rootSig ECDSA-P256 \
-interSig ECDSA-P256 \
-issuerSig ECDSA-P256 \
--pqorderroot Dilithium2 \
--pqorderissuer Dilithium2 \
--pqorderport 10004 \
-loadtestfinalize