FROM golang:1.19.1-alpine3.16

WORKDIR /app

COPY go.mod ./
RUN go get timbatrading
RUN go mod download

COPY *.go ./

RUN go build -o timbatrading

EXPOSE 8080

CMD [ "/timbatrading-backend" ]