FROM golang:1.16-alpine

WORKDIR /src

RUN wget https://raw.githubusercontent.com/eficode/wait-for/v2.1.2/wait-for
RUN chmod 755 ./wait-for

COPY . .

RUN go mod download

RUN go build -o ./bin/chatapp ./cmd/chatapp
