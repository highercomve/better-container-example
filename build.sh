#!/bin/sh
set -e

echo "Building for $TARGETPLATFORM" 
export CGO_ENABLED=0

case "$TARGETPLATFORM" in
	"linux/arm/v6"*)
		export GOOS=linux GOARCH=arm GOARM=6
		;;
	"linux/arm/v7"*)
		export GOOS=linux GOARCH=arm GOARM=7
		;;
	"linux/arm64"*)
		export GOOS=linux GOARCH=arm64 GOARM=7
		;;
	"linux/386"*)
		export GOOS=linux GOARCH=386
		;;
	"linux/amd64"*)
		export GOOS=linux GOARCH=amd64
		;;
	"linux/mips"*)
		export GOOS=linux GOARCH=mips
		;;
	"linux/mipsle"*)
		export GOOS=linux GOARCH=mipsle
		;;
	"linux/mips64"*)
		export GOOS=linux GOARCH=mips64
		;;
	"linux/mips64le"*)
		export GOOS=linux GOARCH=mips64le
		;;
	"linux/riscv64"*)
		export GOOS=linux GOARCH=riscv64
		;;
	*)
		echo "Unknown machine type: $machine"
		echo "Building using host architecture"
esac

go mod download
go build -ldflags="-s -w -extldflags=-static" -o bce -v .
