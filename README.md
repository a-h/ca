# Certificate Authority

Follows guide at https://jamielinux.com/docs/openssl-certificate-authority/index.html

The create-root-pair-* and create-intermediate-pair-* tasks set up the Certificate Authority.

The ca-* tasks use the Certificate Authority (CA) to create certificates for use in web servers etc.

## Tasks

### create-root-pair-prepare-directory

https://jamielinux.com/docs/openssl-certificate-authority/create-the-root-pair.html#prepare-the-directory

```bash
mkdir -p ca
mkdir -p ca/certs ca/crl ca/newcerts ca/private
chmod 700 ca/private
touch ca/index.txt
echo 1000 > ca/serial
```

### create-root-pair-prepare-config

https://jamielinux.com/docs/openssl-certificate-authority/create-the-root-pair.html#prepare-the-configuration-file

```bash
cp root.cnf ca/openssl.cnf
```

### create-root-pair-create-root-key

Dir: ca

This would normally be done on an air-gapped machine.

https://jamielinux.com/docs/openssl-certificate-authority/create-the-root-pair.html#create-the-root-key

```bash
cat /dev/random | head -c 10240 | shasum | head -c 40 > ca_username_password.txt
openssl genrsa -aes256 -passout file:ca_username_password.txt -out private/ca.key.pem 4096
```

### create-root-pair-create-certificate

Dir: ca

Note, change "a-h.github.com" to the domain you want to use.

https://jamielinux.com/docs/openssl-certificate-authority/create-the-root-pair.html#create-the-root-certificate

```bash
openssl req -config openssl.cnf \
      -subj "/C=GB/ST=London/L=London/O=Organisation/OU=IT Department/CN=a-h.github.com" \
      -passin file:ca_username_password.txt \
      -key private/ca.key.pem \
      -new -x509 -days 7300 -sha256 -extensions v3_ca \
      -out certs/ca.cert.pem
```

### create-root-pair-verify-root-certificate

Dir: ca

https://jamielinux.com/docs/openssl-certificate-authority/create-the-root-pair.html#verify-the-root-certificate

```bash
openssl x509 -noout -text -in certs/ca.cert.pem
```

### create-intermediate-pair-prepare-directory

```bash
mkdir -p ca/intermediate
mkdir -p ca/intermediate/certs ca/intermediate/crl ca/intermediate/csr ca/intermediate/newcerts ca/intermediate/private
chmod 700 ca/intermediate/private
touch ca/intermediate/index.txt
echo 1000 > ca/intermediate/serial
echo 1000 > ca/intermediate/crlnumber
```

### create-intermediate-pair-prepare-config

```bash
cp intermediate.cnf ca/intermediate/openssl.cnf
```

### create-intermediate-pair-create-intermediate-key

Dir: ca

Note, change "a-h.github.com Intermediate CA" to the domain you want to use.

https://jamielinux.com/docs/openssl-certificate-authority/create-the-intermediate-pair.html#create-the-intermediate-key

```bash
cat /dev/random | head -c 10240 | shasum | head -c 40 > intermediate_username_password.txt
openssl genrsa -aes256 \
      -passout file:intermediate_username_password.txt \
      -out intermediate/private/intermediate.key.pem 4096
chmod 400 intermediate/private/intermediate.key.pem
```

### create-intermediate-pair-create-intermediate-certificate-create-csr

Dir: ca

Note, change "a-h.github.com Intermediate CA" to the domain you want to use.

https://jamielinux.com/docs/openssl-certificate-authority/create-the-intermediate-pair.html#create-the-intermediate-certificate

```bash
openssl req -config intermediate/openssl.cnf -new -sha256 \
      -subj "/C=GB/ST=London/L=London/O=Organisation/OU=IT Department/CN=a-h.github.com Intermediate CA" \
      -passin file:intermediate_username_password.txt \
      -key intermediate/private/intermediate.key.pem \
      -out intermediate/csr/intermediate.csr.pem
```

### create-intermediate-pair-create-intermediate-certificate-sign-csr

Dir: ca

https://jamielinux.com/docs/openssl-certificate-authority/create-the-intermediate-pair.html#create-the-intermediate-certificate

```bash
openssl ca -config openssl.cnf -extensions v3_intermediate_ca \
      -batch \
      -days 3650 -notext -md sha256 \
      -passin file:ca_username_password.txt \
      -in intermediate/csr/intermediate.csr.pem \
      -out intermediate/certs/intermediate.cert.pem
chmod 444 intermediate/certs/intermediate.cert.pem
```

### create-intermediate-pair-verify-intermediate-certificate

Dir: ca

https://jamielinux.com/docs/openssl-certificate-authority/create-the-intermediate-pair.html#verify-the-intermediate-certificate

```bash
openssl verify -CAfile certs/ca.cert.pem \
      intermediate/certs/intermediate.cert.pem
```

### create-intermediate-pair-create-certificate-chain

Dir: ca

https://jamielinux.com/docs/openssl-certificate-authority/create-the-intermediate-pair.html#create-the-certificate-chain-file

```bash
cat intermediate/certs/intermediate.cert.pem \
      certs/ca.cert.pem > intermediate/certs/ca-chain.cert.pem
chmod 444 intermediate/certs/ca-chain.cert.pem
```

### ca-create-key

Creates a key for the domain. This is typically done by the requestor of the certificate. The requestor then creates a Certificate Signing Request (CSR) for the CA to sign.

Dir: ca/intermediate
Inputs: DOMAIN

https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html#create-a-key

```bash
cat /dev/random | head -c 10240 | shasum | head -c 40 > "$DOMAIN_username_password.txt"
openssl genrsa -aes256 \
      -passout file:"$DOMAIN_username_password.txt" \
      -out private/$DOMAIN.key.pem 2048
chmod 400 private/$DOMAIN.key.pem
```

### ca-create-csr

https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html#create-a-certificate

Dir: ca/intermediate
Inputs: DOMAIN

```bash
openssl req -config ./openssl.cnf \
      -subj "/CN=$DOMAIN" \
      -passin file:"$DOMAIN_username_password.txt" \
      -key private/$DOMAIN.key.pem \
      -new -sha256 \
      -out csr/$DOMAIN.csr.pem
```

### ca-sign-csr

https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html#create-a-certificate

Dir: ca/intermediate
Inputs: DOMAIN

```bash
openssl ca -config ./openssl.cnf \
      -batch \
      -extensions server_cert -days 375 -notext -md sha256 \
      -passin file:../intermediate_username_password.txt \
      -in csr/$DOMAIN.csr.pem \
      -out certs/$DOMAIN.cert.pem
chmod 444 certs/$DOMAIN.cert.pem
```

### ca-view-certs

Dir: ca/intermediate

```bash
cat index.txt
```

### ca-verify-cert

Dir: ca/intermediate
Inputs: DOMAIN

```bash
openssl verify -CAfile certs/ca-chain.cert.pem \
      certs/$DOMAIN.cert.pem
```

### ca-deploy-cert

https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html#deploy-the-certificate

Dir: ca/intermediate
Inputs: DOMAIN

```bash
echo -e "Provide the following to the requestor:\n\n" \
 "  certs/ca-chain.cert.pem\n" \
 "  certs/$DOMAIN.cert.pem\n" \
 "  private/$DOMAIN.key.pem\n" \
 "\n" \
 "In the case that the CSR came from a 3rd party, you won't have the private key, they have that themselves, so you can just" \
 "provide the two cert files."
```
