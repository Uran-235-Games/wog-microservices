FROM alpine:latest

WORKDIR /app

COPY server.exe .

RUN apk add --no-cache ca-certificates

EXPOSE 1488

CMD ["./server"]
