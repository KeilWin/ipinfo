#!/bin/bash

PROJECT_NAME="ipinfo"

mkdir -p bin
mkdir -p bin/ssl
mkdir -p bin/logs
mkdir -p bin/data
mkdir -p bin/configs

openssl req -x509 -out bin/ssl/$PROJECT_NAME.crt -keyout bin/ssl/$PROJECT_NAME.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

cp configs/$PROJECT_NAME.env.example bin/configs/$PROJECT_NAME.env
