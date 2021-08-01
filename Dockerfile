FROM golang:1.16

WORKDIR /app

COPY ./server .

RUN go mod download

RUN go get github.com/go-delve/delve/cmd/dlv

CMD [ "go", "run", "." ]