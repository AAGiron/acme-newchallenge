#!/bin/bash

source config.sh

cd ${WRAPPED_CERT_TESTS_DIR}/src

echo -e "Client: Executing PKIELP 1º pre-quantum handshake\n"

go run launch_client.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-ipclient ${CLIENT_IP} \
-pebbleroot ${DOCKER_VOLUME}/root_ca/root_ca_pebble.pem \
-pskdir ${DOCKER_VOLUME}/storage/client_psk_db.csv \
-wrappedcert \
-pqtls \
-tspath ${DOCKER_VOLUME}/storage/truststore.jks \
-prescenario \
-wrapalgo Ascon80pq \
-syncclient \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements

echo -e "Client: Executing PKIELP 2º pre-quantum handshake\n"

go run launch_client.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-pebbleroot ${DOCKER_VOLUME}/root_ca/root_ca_pebble.pem \
-pskdir ${DOCKER_VOLUME}/storage/client_psk_db.csv \
-wrappedcert \
-pqtls \
-tspath ${DOCKER_VOLUME}/storage/truststore.jks \
-prescenario \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements

echo -e "Client: Executing PKIELP post-quantum handshake\n"

go run launch_client.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-pebbleroot ${DOCKER_VOLUME}/root_ca/root_ca_pebble.pem \
-pskdir ${DOCKER_VOLUME}/storage/client_psk_db.csv \
-wrappedcert \
-pqtls \
-postscenario \
-tspath ${DOCKER_VOLUME}/storage/truststore.jks \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements
