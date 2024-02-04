FROM golang:1.21 as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod tidy

RUN go build -o main ./db_experiments/hazelcast/exp/pessimistic_lock_example.go

CMD ["./main"]