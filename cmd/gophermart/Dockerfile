FROM golang:1.22.7
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /gophermart ./cmd/gophermart/main.go

FROM alpine
WORKDIR /
COPY --from=0 /gophermart /gophermart
EXPOSE 8080
ENTRYPOINT ["/gophermart"]