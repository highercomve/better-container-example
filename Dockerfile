FROM --platform=$BUILDPLATFORM golang:alpine as builder 

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /app

COPY . /app/

RUN /app/build.sh

FROM scratch

WORKDIR /app

COPY --from=builder /app/bce /app/bce
COPY static /app/static/
COPY templates /app/templates/

EXPOSE 5000

ENTRYPOINT [ "/app/bce" ]
