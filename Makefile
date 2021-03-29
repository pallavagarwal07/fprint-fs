main:
	go build main.go

install: main
	sudo mv main /usr/bin/fprint-fs
