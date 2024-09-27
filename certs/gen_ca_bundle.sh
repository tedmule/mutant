#!/usr/bin/env bash

SAN=san.cnf
SERVER_NAME=server
ROOT_NAME=rootCA

# Generate server key and csr 
openssl req -newkey rsa:2048 -nodes -sha256 -keyout ${SERVER_NAME}.key -out ${SERVER_NAME}.csr -config san.cnf

# Sign server certificate using CA
openssl x509 -req -in ${SERVER_NAME}.csr -CA ${ROOT_NAME}.crt -CAkey ${ROOT_NAME}.key -out ${SERVER_NAME}.crt -CAcreateserial -days 3650 -extfile $SAN -extensions v3_req

# generate ca bundle
cat rootCA.crt | base64 -w 0
