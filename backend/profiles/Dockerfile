FROM golang:latest

ADD . /go/src/github.com/BearCloud/fa20-project-dev/backend/profile

WORKDIR /go/src/github.com/BearCloud/fa20-project-dev/backend/profile

RUN go mod download

RUN go build -o main .

EXPOSE 82

ENTRYPOINT [ "./main" ]
