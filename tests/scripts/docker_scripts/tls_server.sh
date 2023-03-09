#!/bin/bash

source config.sh

CURRENT_DIR=$PWD

echo -e "\nCleaning server PSK database and deleting server certificates...\n"

rm -rf ${LEGO_DIR}/.lego

cd ${DOCKER_VOLUME}/storage/
rm -f client_psk_db.csv server_psk_db.csv truststore.jks

cd ${LEGO_DIR}

echo -e "Ordering LEGO to issue a classic certificate for the server's 1º pre-quantum handshake\n"

go run cmd/lego/main.go \
-s https://${PEBBLE_IP}:14000/dir \
-d ${SERVER_NAME} \
-m ${SERVER_NAME}@${SERVER_NAME} \
--http.port ":5002" \
--http \
--certlabel normal_cert \
-a \
--pqtls \
--kex Kyber512 \
-k ec256 \
--certalgo P256 \
--timingcsv ${DOCKER_VOLUME}/measurements/lego_issuance_time.csv \
--synclego \
--serverip ${SERVER_NAME} \
run

cd ${WRAPPED_CERT_TESTS_DIR}/src

echo -e "Server: Executing PKIELP 1º pre-quantum handshake\n"

go run launch_servers.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-wrapdir ${LEGO_DIR}/.lego/certificates \
-pskdir ${DOCKER_VOLUME}/storage/server_psk_db.csv \
-wrappedcert \
-pqtls \
-prescenario \
-ocspresponsepath ${DOCKER_VOLUME}/resources/ocsp_response \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements

echo -e "Ordering LEGO to issue a wrapped certificate for the server's 2º pre-quantum handshake\n"

cd ${WRAPPED_CERT_TESTS_DIR}/scripts/local_scripts/pkielp_simulation

python3 read_server_psk_db.py ${DOCKER_VOLUME}/storage ${LEGO_DIR} ${DOCKER_VOLUME}/measurements/lego_issuance_time.csv ${PEBBLE_IP} ${SERVER_NAME}

# The server only stores Cert PSKs until it issues the corresponding wrapped certificates
> ${DOCKER_VOLUME}/storage/server_psk_db.csv

cd ${WRAPPED_CERT_TESTS_DIR}/src

echo -e "Server: Executing PKIELP 2º pre-quantum handshake\n"

go run launch_servers.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-wrapdir ${LEGO_DIR}/.lego/certificates \
-wrappedcert \
-pqtls \
-prescenario \
-ocspresponsepath ${DOCKER_VOLUME}/resources/ocsp_response \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements

echo -e "Server: Executing PKIELP post-quantum handshake\n"

go run launch_servers.go common.go stats_tls.go \
-ipserver ${SERVER_NAME} \
-wrapdir ${LEGO_DIR}/.lego/certificates \
-wrappedcert \
-pqtls \
-ocspresponsepath ${DOCKER_VOLUME}/resources/ocsp_response \
-port 4433 \
-measurements ${DOCKER_VOLUME}/measurements
