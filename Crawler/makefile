all: dependencies verify build

build: api.go config.go json.go main.go utility.go
	go build -o main *.go 

dependencies:
	go get github.com/BurntSushi/toml/cmd/tomlv

verify: config.toml
	tomlv config.toml