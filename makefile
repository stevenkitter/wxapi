build:
	protoc --proto_path=. --micro_out=. --go_out=. ./proto/wxAuth/wxAuth.proto
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./api/app ./api
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./auth/app ./auth
