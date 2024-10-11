FROM golang:1.22.0-alpine3.19

WORKDIR /app

COPY . .
COPY ./go.mod .
COPY ./go.sum .

# Download and install the dependencies:
RUN apk update && apk add git
RUN go get -d -v ./...

# Build the go app
RUN go build -o build main.go

CMD ./build
