#!/bin/sh
set -e

echo "Building for $TARGETPLATFORM" 
export CGO_ENABLED=0 GOOS=linux 

case "$TARGETPLATFORM" in
	"linux/arm/v6"*)
		export GOARCH=arm GOARM=6
		;;
	"linux/arm/v7"*)
		export GOARCH=arm GOARM=7
		;;
	"linux/arm64"*)
		export GOARCH=arm64 GOARM=7
		;;
	"linux/386"*)
		export GOARCH=386
		;;
	"linux/amd64"*)
		export GOARCH=amd64
		;;
	"linux/mips"*)
		export GOARCH=mips
		;;
	"linux/mipsle"*)
		export GOARCH=mipsle
		;;
	"linux/mips64"*)
		export GOARCH=mips64
		;;
	"linux/mips64le"*)
		export GOARCH=mips64le
		;;
	"linux/riscv64"*)
		export GOARCH=riscv64
		;;
	*)
		echo "Unknown machine type: $machine"
		echo "Building using host architecture"
esac

go mod download
go build -ldflags="-s -w -extldflags=-static" -o bce -v .
