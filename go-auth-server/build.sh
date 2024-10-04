GOOS=linux GOARCH=arm64 go build -o bin/go-auth-server
docker build . -t my-auth-server