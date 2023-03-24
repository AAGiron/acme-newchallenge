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

Refer to `tests/scripts/local_scripts/` for load-testing and scripts for issuance of certificates. Open two terminals, execute Pebble first, and then the client (LEGO).

The certificates are normally stored in `go-lego/.lego/certificates/` (check LEGO output logs). Measurements are given as CSV files in `measurements`, but you have to give flags for specific metrics. 

# Documentation

A documentation in `docs/` provides details on the PQ-Transition Challenge implementation. For additional information, overview, design and results please refer to the paper.

Tested on a (fresh) Ubuntu LTS [multipass](https://multipass.run/) [modified](https://multipass.run/docs/modify-an-instance) instance, with 2GB memory and 20 GB disk. 

# Further configurations

One can change Pebble's server certificate by a PQC one. See `go-pebble/test/certs/` for examples. They can be configured before launching `./pebble.sh` in `go-pebble/test/config/pebble-config.json`, and then you can change the algorithms in the example scripts (`pebble.sh` and 

# Disclaimer

This is a prototype implementation for benchmarking, demonstration and experimentation purposes. 

Suggestions and contributions are welcome!

# Known issues

Anonymous github does not allow downloading big repositories (nor submodules) so if you want to:
- Look at the source of the submodules (anonymized): [Go-std](https://anonymous.4open.science/r/go-std-C24A), [go-pebble](https://anonymous.4open.science/r/go-pebble-78DE/), [Go-JOSE](https://anonymous.4open.science/r/go-jose-5555), [Go-LEGO](https://anonymous.4open.science/r/go-lego-2E5F). A look in the `docs/` might give directions.
- Download the source: [PQTransitionACMEChallenge-Sourcev1.0.zip](https://mega.nz/file/S8khgZ4Z#3b55kBbXonaMPMlz5CKse92FbbsfB4MTeI8CaRilIJE), unzip it, then refer to `tests/scripts/local_scripts/` for installation and execution scripts. The zip file includes everything but the main requirement (a Go installation) still apply.

If an error like `Temporary naming resolution failure` appears when testing a certificate issuance, please check if your `/etc/hosts` contains the line `${IP_SERVER}   ${SERVER_NAME}`. The installation script adds it but some VM instances can flush out such a configuration.  
