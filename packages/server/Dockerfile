FROM golang:1.19.2-buster

WORKDIR /server

COPY go.mod .
RUN go mod download
COPY . .

RUN go build -o area_server .

CMD ["./area_server"]