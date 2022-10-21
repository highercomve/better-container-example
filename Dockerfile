FROM golang:alpine as builder 

WORKDIR /app

COPY . /app/

RUN go mod download; \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o bce -v .

FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=builder /app/bce /app/bce
COPY static /app/static/
COPY templates /app/templates/

EXPOSE 5000

ENTRYPOINT [ "/app/bce" ]
