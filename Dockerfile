FROM golang:1.16

WORKDIR /app

COPY . .

RUN go get github.com/go-delve/delve/cmd/dlv

CMD [ "go", "run", "." ]