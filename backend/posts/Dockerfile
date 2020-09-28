FROM golang:latest

ADD . /go/src/github.com/BearCloud/fa20-project-dev/backend/posts

WORKDIR /go/src/github.com/BearCloud/fa20-project-dev/backend/posts

RUN go mod download

RUN go build -o main .

EXPOSE 81

ENTRYPOINT [ "./main" ]
