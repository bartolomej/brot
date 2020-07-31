echo "Building linux binaries..."
GOOS=windows go build -o bin/brot
echo "Building windows binaries..."
GOOS=linux go build -o bin/brot.exe
