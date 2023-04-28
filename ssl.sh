set -ex
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr
openssl x509 -signkey server.key -in server.csr -req -days 365 -out server.crt
