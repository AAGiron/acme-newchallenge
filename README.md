# ACME-newchallenge

Repository for the implemenation of the ACME protocol with Post-Quantum Cryptography (PQC), including the "PQ-Transition Challenge" proposal. The source code is based on Go-LEGO, Go-Pebble, Go-JOSE, and Go standard library (i.e., the submodules).

# Requirements and Installation

- Install the Go programming language and add it's binary to your PATH environment variable/Instructions available at: https://go.dev/doc/install
- Clone/Download this repository, making sure that all required submodules are downloaded (`git submodule update --init --recursive`)
- Run the installation script in `tests/scripts/installation/`:
```
./install_local.sh
```

# Execution

Refer to `tests/scripts/local_scripts/` for load-testing and scripts for issuance of certificates. Open two terminals, execute Pebble first, and then the client (LEGO). 

The certificates are normally stored in `go-lego/.lego/certificates/` (check LEGO output logs). Measurements are given as CSV files in `tests/measurements/`, but you have to give flags for specific metrics. 

# Documentation

A documentation in `docs/` provides details on the PQ-Transition Challenge implementation. For additional information, overview, design and results please refer to the paper (to appear).

Tested on a (fresh) Ubuntu 22.10 LTS [multipass](https://multipass.run/) [modified](https://multipass.run/docs/modify-an-instance) instance, with 2GB memory and 20 GB disk. 

# Further configurations

One can change Pebble's server certificate by a PQC one. See `go-pebble/test/certs/` for examples. They can be configured before launching `./pebble.sh` in `go-pebble/test/config/pebble-config.json`, and then you can change the algorithms in the example scripts (`pebble.sh` and `lego.sh` or `lego-newchallenge.sh`).

# Disclaimer

This is a prototype implementation for benchmarking, demonstration and experimentation purposes. 

Suggestions and contributions are welcome!

# Known issues

If an error like `Temporary naming resolution failure` appears when testing a certificate issuance, please check if your `/etc/hosts` contains the line `${IP_SERVER}   ${SERVER_NAME}`. The installation script adds it but some VM instances can flush out such a configuration.  

If an error like `Get "https://127.0.0.1:14000/dir": x509: certificate signed by unknown authority` appears please check if Pebble's TLS certificate is trusted (it should be, after the installation script, nevertheless you can add it yourself in your running instance `sudo cp <go-pebble-dir>/test/certs/pebble.minica.pem /etc/ssl/certs/pebble.minica.crt && sudo update-ca-certificates`).
