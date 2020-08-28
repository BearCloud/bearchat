FROM golang:latest

ADD . .

RUN go mod download

RUN go build -o main .

EXPOSE 80

ENTRYPOINT [ "./main" ]