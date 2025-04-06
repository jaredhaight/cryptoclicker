.PHONY: vet
vet:
	go vet .\...

.PHONY: build
build:
	go build -o .\clicker.exe .\cmd\main.go