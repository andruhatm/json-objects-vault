FROM golang:1.19-alpine AS build

RUN apk add --no-cache git

WORKDIR /json-objects-vault

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/app .

EXPOSE 8081

CMD ["./out/app"]