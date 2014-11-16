#!/bin/bash

set -e

echo -n "What is your server's domain?"
read host

echo "# Certificate Authority Configuration"

mkdir -p ca

if [ -f ca/ca.pass ]; then
  echo "Certificate Authority password already exists; not creating"
  echo "(ca/ca.pass)"
else
  echo "Creating CA password (ca/ca.pass)"
  openssl rand -base64 32 > ca/ca.pass
fi

if [ -f ca/ca.srl ]; then
  echo "Certificate Serial exists; not creating. (ca/ca.srl)"
else
  echo "Creating Certificate Serial (ca/ca.srl)"
  echo 01 > ca/ca.srl
fi

if [ -f ca/ca-key.pem ]; then
  echo "Found Certificate Authority key (ca/ca-key.pem). Not creating a new"
  echo "one."
else
  echo "Generating your Certificate Authority key (ca/ca-key.pem)."
  openssl genrsa -aes256 -passout file:ca/ca.pass -out ca/ca-key.pem 2048
fi

openssl rsa -in ca/ca-key.pem -passin file:ca/ca.pass \
  -out ca/ca-key-insecure.pem
trap "rm ca/ca-key-insecure.pem" 1 2

if [ -f ca/ca.pem ]; then
  echo "Found Certificate Authority Certificate (ca/ca.pem). Not creating a"
  echo "new one."
else
  echo "Generating your Certificate Authority's Certificate (ca/ca.pem)"
  echo "   daemon"
  SUBJ="/C=US/ST=Online/O=Docker Remote Builds/CN=$host"
  openssl req -new -x509 -days 365 -subj "$SUBJ" -key ca/ca-key-insecure.pem \
    -out ca/ca.pem
fi
echo ""
echo " --> Certificate Authority section is done. Files are in ca/."
echo ""
echo " *** CAUTION: DO NOT LOSE THESE FILES. ***"
echo ""
echo ""
echo ""

echo "# Server Configuration"
mkdir -p server
if [ -f server/server-key.pem ]; then
  echo "Found a Daemon key (server/server-key.pem). Not creating a new one."
else
  echo "Creating a server key (server/server-key.pem)"
  openssl genrsa -out server/server-key.pem 2048
fi

if [ -f server/server.csr ]; then
  echo "Found a Certificate Signing Request for the the daemon"
  echo "(server/server.csr)."
  echo "Not creating a new one."
else
  echo "Generating the Certificate Signing Request for your Docker Daemon"
  echo "(server/server.csr)."
  SUBJ="/C=US/ST=Online/O=Docker Remote Builds/OU=Server/CN=$host"
  openssl req -new -key server/server-key.pem -subj "$SUBJ" \
    -out server/server.csr
fi

if [ -f server/server-cert.pem ]; then
  echo "Found a signed daemon certificate (server/server-cert.pem). Not"
  echo "creating a new one."
else
  echo "Signing the Docker daemon certificate with the Certificate Authority"
  openssl x509 -req -days 365 -in server/server.csr -CA ca/ca.pem -CAkey \
    ca/ca-key-insecure.pem -CAserial ca/ca.srl -out server/server-cert.pem
fi

./userdata.sh > user-data

echo ""
echo ""
echo " --> Server configuration is complete. The files are in server/"
echo ""
echo ""
echo ""

echo "### What project would you like to create a certificate for?"
echo "Notes:"
echo " - Please use only A-Za-z0-9 and dashes and underscores for the name."
echo -n "Project name:"
read proj

mkdir "client-$proj"
cd "client-$proj"

cp -r "../ca/ca.pem" ".docker.ca"

export PASS=$(openssl rand -base64 32)
openssl genrsa -out ".docker.insecure.key" 2048
SUBJ="/C=US/ST=Online/O=Docker Remote Builds/OU=Client/CN=$proj"
openssl req -subj "$SUBJ" -new -key ".docker.insecure.key" \
  -out "client.csr"
echo "extendedKeyUsage = clientAuth" > extfile.cnf
echo "Signing your project's Certificate Signing Request with the Certificate"
echo "Authority's key. Use the same password from when creating the"
echo "Certificate Authority's key."
openssl x509 -req -in "client.csr" -CA ../ca/ca.pem \
  -CAkey ../ca/ca-key-insecure.pem -CAserial ../ca/ca.srl -out ".docker.crt" \
  -extfile extfile.cnf
openssl rsa -in .docker.insecure.key -out .docker.key -aes256 \
  -passout env:PASS
cat .docker.crt .docker.key > .docker.pem

HOST="$host" ../travis.sh > .travis.yml

rm client.csr extfile.cnf .docker.crt .docker.key .docker.insecure.key \
  ../ca/ca-key-insecure.pem

echo ""
echo " --> Client configuration ($proj) is complete. Re-run this script to"
echo "     make an additional project."
echo ""
echo " --> A CoreOS compatible UserData file has been placed at ./user-data."
echo "     This file can be placed in the user-data field of Amazon Web"
echo "     Services when launching a CoreOS instance. This instance will be"
echo "     setup to run your remote docker host. This is an example file only,"
echo "     and does not appropriately secure or configure or maintain the rest"
echo "     of the system."
echo ""
echo " --> A .travis.yml example (named travis.yml) has been placed in each"
echo "     the project directory. You must include ALL files in the client"
echo "     directory."
echo ""
echo " --> The travis project will not work until you add the password to the"
echo "     .travis.yml. After moving ALL the files in client-$proj/ to your"
echo "     repository, set up the travis CLI, and run:"
echo "         travis encrypt DOCKER_PASS=$PASS"
echo "     from within the project directory."
echo ""
echo " --> Look at this repository at https://github.com/grahamc/taxi for a"
echo "     full example of a repository using this system."
echo ""
echo ""
echo " ... Good luck."
echo ""
echo " - Graham Christensen"
echo ""

unset PASS

