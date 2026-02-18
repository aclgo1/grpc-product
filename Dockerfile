FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-product ./cmd/product/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/grpc-product ./

COPY --from=builder /app/.env ./

EXPOSE 50054

ENTRYPOINT ["./grpc-product"]