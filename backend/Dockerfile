FROM golang:1.21 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
go build -a -installsuffix cgo -o /main cmd/api/main.go

EXPOSE 9001

FROM alpine:3
COPY --from=builder main /app/main
COPY config/ /app/config/
ENTRYPOINT ["/app/main"]
