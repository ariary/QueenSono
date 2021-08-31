before.build:
	go mod download && go mod vendor

build.queensono-sender:
	@echo "build in ${PWD}";go build -o qssender cmd/client/main.go;sudo setcap cap_net_raw+eip qssender

build.queensono-receiver:
	@echo "build in ${PWD}";go build -o qsreceiver cmd/server/main.go;sudo setcap cap_net_raw+eip qsreceiver
	
