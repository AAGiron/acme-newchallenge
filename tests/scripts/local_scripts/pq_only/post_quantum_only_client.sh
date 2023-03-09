#!/bin/bash

source ../../config.sh

cd ${WRAPPED_CERT_TESTS_DIR}/src

go run launch_client.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-ipclient ${CLIENT_IP} \
-pqtls \
-pebbleroot ${WRAPPED_CERT_TESTS_DIR}/root_ca/root_ca_pebble.pem \
-measurements ${WRAPPED_CERT_TESTS_DIR}/measurements

