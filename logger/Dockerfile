FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD [ "./main" ]