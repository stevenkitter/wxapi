FROM alphine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY ./app .
CMD ["./app"]