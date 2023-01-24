build:
	go build -o bin/sync ./cmd/sync
compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/sync-linux-amd64 ./cmd/sync
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/sync-darwin-amd64 ./cmd/sync
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/sync-windows-amd64.exe ./cmd/sync
