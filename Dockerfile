FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]
