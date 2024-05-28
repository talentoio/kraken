FROM golang:1.21-alpine

WORKDIR /app

# Copying all the files
COPY go.mod go.sum ./

# Installing dependencies
RUN go mod download

COPY . .

RUN go build -o kraken ./cmd

RUN ls

EXPOSE 8080

CMD [ "./kraken" ]