VERSION=0.0.5

cbatch-linux-amd64: 
	GOOS=linux GOARCH=amd64 go build -o cbatch-linux-amd64

build:
	go build -o cbatch

release: cbatch-linux-amd64 config/config.toml config/bootstrap.tmpl.sh
	tar cvf cbatch-${VERSION}-linux-amd64.tar cbatch-linux-amd64 config/config.toml config/bootstrap.tmpl.sh 

clean: 
	-rm cbatch-${VERSION}-linux-amd64.tar cbatch-linux-amd64 cbatch