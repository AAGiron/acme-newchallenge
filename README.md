# ACME-newchallenge

Repository for the implemenation of the ACME protocol with Post-Quantum Cryptography (PQC), including the "PQ-Transition Challenge" proposal. The source code based on Go-LEGO, Go-Pebble, Go-JOSE, and Go standard library (i.e., the submodules).

# Requirements and Installation

- Install the Go programming language and add it's binary to your PATH environment variable/Instructions available at: https://go.dev/doc/install
- Clone this repository, making sure that all required submodules are downloaded.
- Run the installation script in `tests/scripts/installation/`:
```
./install_local.sh
```

# Execution

Refer to `tests/scripts/local_scripts/` for load-testing and scripts for standard issuance of certificates. Execute Pebble first, and then the client (LEGO).

# Disclaimer

This is a prototype implementation for benchmarking, demonstration and experimentation purposes. 

Suggestions and contributions are welcome!
