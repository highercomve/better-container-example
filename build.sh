#!/bin/sh
set -e

echo "Building for $TARGETPLATFORM" 

case "$TARGETPLATFORM" in
	"linux/arm/v6"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6
		;;
	"linux/arm/v7"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7
		;;
	"linux/arm64"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=7
		;;
	"linux/386"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=386
		;;
	"linux/amd64"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=amd64
		;;
	"linux/mips"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=mips
		;;
	"linux/mipsle"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=mipsle
		;;
	"linux/mips64"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=mips64
		;;
	"linux/mips64le"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=mips64le
		;;
	"linux/riscv64"*)
		export CGO_ENABLED=0 GOOS=linux GOARCH=riscv64
		;;
	*)
		echo "Unknown machine type: $machine"
		exit 1
esac

go mod download
go build -ldflags="-s -w -extldflags=-static" -o better-container-example -v .
