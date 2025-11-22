FROM golang:1.25
WORKDIR /authorization
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server.go
CMD [ "./server" ]