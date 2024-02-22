#!/bin/bash

# Generates a self-signed CA certificate and a private key.
openssl req -new \
	-newkey rsa:4096 \
	-keyout data/ca.key \
	-subj /C=US/ST=NY/L=NY/O=crypto/CN=localhost \
	-x509 \
	-sha256 \
	-days 365 \
	-out data/ca.crt

# Generates a private key for the API.
openssl genrsa -out data/api.key 4096

# Generates an API CSR.
openssl req -new \
	-key data/api.key \
	-subj /C=US/ST=NY/L=NY/O=api/CN=localhost \
	-config ./configs/openssl.cnf \
	-out data/api.csr

# Generates an API certificate signed by the CA and a private key.
openssl x509 \
	-req \
	-in data/api.csr \
	-extfile ./configs/openssl.cnf \
	-extensions v3_req \
	-CA data/ca.crt \
	-CAkey data/ca.key \
	-CAcreateserial \
	-days 365 \
	-sha256 \
	-out data/api.crt

# Verifies validity of certs.
openssl verify -verbose -CAfile data/ca.crt data/api.crt

# # Generates a private key for mq.
# openssl genrsa -out data/mq.key 4096
#
# # Generates a mq CSR.
# openssl req -new \
# 	-key data/mq.key \
# 	-subj /C=US/ST=NY/L=NY/O=mq/CN=mq \
# 	-config ./configs/openssl.cnf \
# 	-out data/mq.csr
#
# # Generates a mq certificate signed by the CA and a private key.
# openssl x509 \
# 	-req \
# 	-in data/mq.csr \
# 	-extfile ./configs/openssl.cnf \
# 	-extensions v3_req \
# 	-CA data/ca.crt \
# 	-CAkey data/ca.key \
# 	-CAcreateserial \
# 	-days 365 \
# 	-sha256 \
# 	-out data/mq.crt
#
# # Verifies validity of certs.
# openssl verify -verbose -CAfile data/ca.crt data/mq.crt

# Chmod so services can read the files.
chmod +rw data/*key
