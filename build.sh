go mod tidy

mkdir -p tmp/

echo "building for darwin arm64"
GOOS="darwin" GOARCH="arm64" go build -ldflags="-s -w" -o tmp/build/ceres-darwin cmd/debug_tool/main.go

echo "building for linux amd64"
GOOS="linux" GOARCH="amd64" go build -ldflags="-s -w" -o tmp/build/ceres cmd/debug_tool/main.go

echo "building for windows amd64"
GOOS="windows" GOARCH="amd64" go build -ldflags="-s -w" -o tmp/build/ceres.exe cmd/debug_tool/main.go

echo "done building Ceres debug tool"
