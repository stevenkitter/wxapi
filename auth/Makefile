build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
	docker build -t coderxu/wxapi .
	docker push coderxu/wxapi