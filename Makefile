build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/bmclapi-report-linux-amd64 .
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/bmclapi-report-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/bmclapi-report-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/bmclapi-report-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/bmclapi-report-darwin-arm64 .