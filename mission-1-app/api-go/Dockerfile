FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o api-go .

FROM alpine:3.21
WORKDIR /app
COPY --from=build /app/api-go /app/api-go
EXPOSE 8081
CMD ["/app/api-go"]
