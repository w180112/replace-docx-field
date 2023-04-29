set -ex
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr
openssl x509 -signkey server.key -in server.csr -req -days 365 -out server.crt
certs_dir_count=$(ls -l | grep certs | wc -l)
if [ $certs_dir_count -eq 0 ]; then
	mkdir certs
fi
mv server.crt certs/
mv server.key certs/
