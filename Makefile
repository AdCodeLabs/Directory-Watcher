build:
	go build main.go

run:
	go run cmd/directoryWatcher/main.go --type local --path ./ --server 0.0.0.0:8097
