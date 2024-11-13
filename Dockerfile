FROM golang:1.23-alpine

WORKDIR /app

COPY . /app/

RUN apk add build-base
RUN go build -o woningfinder-bin ./cmd/woningfinder

EXPOSE 8080

CMD ["./woningfinder-bin"]
