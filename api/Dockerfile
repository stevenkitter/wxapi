FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY ./api/app .
CMD ["./app"]