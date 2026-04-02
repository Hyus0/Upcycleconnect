FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY api.go ./
COPY app/ ./app/
COPY db/ ./db/
COPY models/ ./models/
RUN go build -o api-go .

FROM alpine:3.21
WORKDIR /app
COPY --from=build /app/api-go /app/api-go
EXPOSE 8081
CMD ["/app/api-go"]
