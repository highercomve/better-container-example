FROM golang:alpine

WORKDIR /app
COPY . /app/
RUN go mod download; \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o bce -v .

ENTRYPOINT [ "/app/bce" ]

