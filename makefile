run:
	go run translate.go

build: translate translate.go
	go build translate.go

install:
	bash ./kill.sh
	sudo cp ./translate /usr/local/bin/