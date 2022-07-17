FROM golang:1.18.4

WORKDIR testdir

COPY go.mod .

RUN go mod download && go mod verify

COPY . .

RUN go test ./...

FROM golang:1.18.4

WORKDIR /

COPY go.mod .

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/proofofwork ./..

CMD ["proofofwork"]