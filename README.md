# ACME-newchallenge

Repository for the implemenation of the ACME protocol with Post-Quantum Cryptography (PQC), including the "PQ-Transition Challenge" proposal. The source code based on Go-LEGO, Go-Pebble, Go-JOSE, and Go standard library (i.e., the submodules).

# Requirements and Installation

- Install the Go programming language and add it's binary to your PATH environment variable/Instructions available at: https://go.dev/doc/install
- Clone/Download this repository, making sure that all required submodules are downloaded (`git submodule update --init --recursive`)
- Run the installation script in `tests/scripts/installation/`:
```
./install_local.sh
```

# Execution

Refer to `tests/scripts/local_scripts/` for load-testing and scripts for issuance of certificates. Execute Pebble first, and then the client (LEGO).

# Documentation

A documentation in `docs/` provides details on the PQ-Transition Challenge implementation. For additional information, overview, design and results please refer to the paper.

# Disclaimer

This is a prototype implementation for benchmarking, demonstration and experimentation purposes. 

Suggestions and contributions are welcome!

# Known issues

Anonymous github does not allow downloading repositories (nor submodules) so if you want to:
- Look at the source (anonymized): [Go-std](https://anonymous.4open.science/r/go-std-C24A), [go-pebble](https://anonymous.4open.science/r/go-pebble-78DE/), [Go-JOSE](https://anonymous.4open.science/r/go-jose-5555), [Go-LEGO](https://anonymous.4open.science/r/go-lego-2E5F)
- Download the source: [PQTransitionACMEChallenge-Sourcev1.0.zip](https://mega.nz/file/S8khgZ4Z#3b55kBbXonaMPMlz5CKse92FbbsfB4MTeI8CaRilIJE), unzip it, then refer to `tests/scripts/local_scripts/` for installation and execution scripts. The zip file includes everything but the main requirement (a Go installation) still apply.
