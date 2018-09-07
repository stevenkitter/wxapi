FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY ./auth/app .
CMD ["./app"]