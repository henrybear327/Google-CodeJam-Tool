all: dependencies verify build

build: main.go config.go
	go build -o main main.go config.go

dependencies:
	go get github.com/BurntSushi/toml/cmd/tomlv

verify: config.toml
	tomlv config.toml