FROM golang:1.21 as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod tidy
