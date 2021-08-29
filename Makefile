before.build:
	go mod download && go mod vendor

build.queensono-client:
	@echo "build in ${PWD}";go build -o qsclient cmd/client/main.go;sudo setcap cap_net_raw+eip qsclient

build.queensono-server:
	@echo "build in ${PWD}";go build -o qsserver cmd/server/main.go;sudo setcap cap_net_raw+eip qsserver
	
