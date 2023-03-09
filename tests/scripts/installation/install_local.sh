#!/bin/bash

read -p "What will be the TLS client's IP? " CLIENT_IP
read -p "What will be the TLS server's NAME (It cannot be the IP address)? " SERVER_NAME
read -p "What will be the ACME server's IP? " PEBBLE_IP

cd ../..

PQCACME_TESTS_DIR=$PWD

cd ../

cd go-std
GO_STD_DIR=$PWD

cd ../liboqs
LIBOQS_DIR=$PWD

cd ../liboqs-go
LIBOQS_GO_DIR=$PWD

cd ${PQCACME_TESTS_DIR}/../..

cd acme-server/go-pebble
PEBBLE_DIR=$PWD

cd ${PQCACME_TESTS_DIR}/..

cd go-lego

LEGO_DIR=$PWD

cd ${PQCACME_TESTS_DIR}

echo -e "\nInstalling liboqs dependencies...\n"

sudo apt install astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml git pkg-config

echo -e "\nCompiling Liboqs\n"

cd ${LIBOQS_DIR}

mkdir build && cd build
cmake -GNinja -DBUILD_SHARED_LIBS=ON ..
sudo ninja
sudo ninja install

cd ${PQCACME_TESTS_DIR}

echo -e "Writing environment variables to ~/.profile and ~/.bashrc\n"

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:${LIBOQS_GO_DIR}/.config

echo "
export LD_LIBRARY_PATH=\$LD_LIBRARY_PATH:/usr/local/lib
export PKG_CONFIG_PATH=\$PKG_CONFIG_PATH:${LIBOQS_GO_DIR}/.config" | tee -a ~/.profile ~/.bashrc


echo -e "\nCompiling Go\n"

cd ${GO_STD_DIR}/src
./make.bash

export PATH=${GO_STD_DIR}/bin:$PATH

echo "export PATH=${GO_STD_DIR}/bin:\$PATH
" | tee -a ~/.profile ~/.bashrc

echo -e "Adding Pebble root CA's certificates to the system trusted CAs list\n"

cd ${PQCACME_TESTS_DIR}/../pebble_https_root_ca
sudo cp * /usr/local/share/ca-certificates
sudo update-ca-certificates

echo -e "Appending '127.0.0.1 ${SERVER_NAME}' to /etc/hosts\n"
sudo echo "127.0.0.1 ${SERVER_NAME}" | sudo tee -a /etc/hosts

echo "PQCACME_TESTS_DIR=${PQCACME_TESTS_DIR}
LEGO_DIR=${LEGO_DIR}
PEBBLE_DIR=${PEBBLE_DIR}
GO_STD_DIR=${GO_STD_DIR}
SERVER_NAME=${SERVER_NAME}
CLIENT_IP=${CLIENT_IP}
PEBBLE_IP=${PEBBLE_IP}
" >> ${PQCACME_TESTS_DIR}/scripts/config.sh

echo -e "Installation completed. Restart your system for the changes to take effect\n"
