build:
	protoc --proto_path=. --micro_out=. --go_out=. ./proto/greeter/greeter.proto
	protoc --proto_path=. --micro_out=. --go_out=. ./proto/wxAuth/wxAuth.proto
