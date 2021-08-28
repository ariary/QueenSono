before.build:
	go mod download && go mod vendor

build.queensono-client:
	@echo "build in ${PWD}";go build -o queensono-client cmd/client/main.go;sudo sudo setcap cap_net_raw+eip queensono-client

build.queensono-server:
	@echo "build in ${PWD}";go build -o queensono-server cmd/client/main.go;sudo sudo setcap cap_net_raw+eip queensono-server
	
