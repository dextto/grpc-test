# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca.key -out ca.crt -subj "/C=KR/ST=Busan/L=Haeundae/O=TEST CA/OU=Education/CN=*.testca.com/emailAddress=testca@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca.crt -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr -subj "/C=KR/ST=Seoul/L=Gangnam/O=Modusign Server/OU=Computer/CN=*.modisign.co.kr/emailAddress=server@modisign.co.kr"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server.csr -days 60 -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server.crt -noout -text
