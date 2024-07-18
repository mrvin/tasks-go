build:
	go build -ldflags '-w -s' -o bin/hh-client-go
certgen:
	openssl req -nodes -x509 -newkey rsa:4096 \
		-keyout cert/clientKey.pem -out cert/clientCert.pem -days 365 \
		-subj "/C=RU/L=Saint Petersburg/O=hh-client-go Corp./OU=hh-client-go/CN=hh-client-go.com/emailAddress=v.v.vinogradovv@gmail.com"

.PHONY: build certgen
