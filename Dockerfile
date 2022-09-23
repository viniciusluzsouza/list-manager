FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /list-manager-app /app/cmd/app/main.go

EXPOSE 8080

CMD [ "/list-manager-app" ]
