rm -rf ./build/server || true
rm -rf ./build/static || true
rm -rf ./build/templates || true

go mod tidy
env GOOS=linux GOARCH=arm64 go build -o build/server_arm
go build -o build/server_linux

cp -rf ./templates build/templates
cp -rf ./static build/static