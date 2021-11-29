FROM golang:1.17-buster
RUN apt-get update && \
  apt-get install -y \
    build-essential gcc libc6-dev-i386 gcc-8-arm-linux-gnueabi gcc-arm-linux-gnueabi