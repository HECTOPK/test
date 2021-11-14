.PHONY: cert client server proto

cert:
	openssl req -x509 -newkey rsa:4096 -nodes -keyout cert/ca-key.pem -out cert/ca-cert.pem -sha256 -days 36500 \
		-subj "/C=US/ST=California/L=Oakland/O=test/OU=test/CN=test1.com/emailAddress=test1@gmail.com"

	openssl req -newkey rsa:4096 -nodes -keyout cert/server-key.pem -out cert/server-req.pem \
		-subj "/C=US/ST=California/L=Oakland/O=test/OU=test/CN=test2.com/emailAddress=test2@gmail.com"

	openssl x509 -req -in cert/server-req.pem -days 36500 -CA cert/ca-cert.pem -CAkey cert/ca-key.pem -CAcreateserial -out cert/server-cert.pem -extfile cert/server-ext.cnf

	openssl req -newkey rsa:4096 -nodes -keyout cert/client-key.pem -out cert/client-req.pem \
		-subj "/C=US/ST=California/L=Oakland/O=test/OU=test/CN=test3.com/emailAddress=test3@gmail.com"

	openssl x509 -req -in cert/client-req.pem -days 36500 -CA cert/ca-cert.pem -CAkey cert/ca-key.pem -CAcreateserial -out cert/client-cert.pem


client:
	go build -o client client.go

server:
	go build -o server server.go

proto:
	protoc --go_out=api/ --go_opt=paths=source_relative \
        --go-grpc_out=api/ --go-grpc_opt=paths=source_relative \
        api.proto