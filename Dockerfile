FROM golang:1.18.4

WORKDIR /

COPY . .
RUN go build -v -o /usr/local/bin/proofofwork ./..

CMD ["proofofwork"]