run:
	go run translate.go

build: translate translate.go
	go build -o ./bin/translate translate.go

install:
	bash ./kill.sh
	sudo cp ./bin/translate /usr/local/bin/