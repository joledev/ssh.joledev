.PHONY: build run build-linux deploy clean

build:
	go build -o ssh.joledev.exe .

run: build
	SSH_PORT=2222 ./ssh.joledev.exe

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ssh.joledev .

deploy: build-linux
	@echo "Copy ssh.joledev binary + data/ + posts/ to your VPS"
	@echo "Then run: sudo bash deploy/setup.sh"

clean:
	rm -f ssh.joledev ssh.joledev.exe .ssh .ssh.pub

test-local:
	@echo "Starting server on port 2222..."
	@echo "Connect with: ssh -p 2222 localhost"
	SSH_PORT=2222 go run .
