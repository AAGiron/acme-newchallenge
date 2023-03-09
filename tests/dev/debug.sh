#!/bin/bash

rm -rf temp_client
rm -rf temp_server

mkdir temp_client
mkdir temp_server

cp -r ../src/{csv,db,graphs,scripts,common.go,launch_client.go,stats_tls.go} temp_client

cp -r ../src/{csv,db,graphs,scripts,common.go,launch_servers.go,stats_tls.go} temp_server