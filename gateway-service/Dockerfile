FROM golang:1.25.5-alpine

WORKDIR /app

RUN apk add --no-cache make

RUN go install github.com/air-verse/air@latest

COPY go.* .

RUN go mod download

COPY . .

COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 80

CMD ["entrypoint.sh"]