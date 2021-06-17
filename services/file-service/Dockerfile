FROM golang:1.16-alpine as base

RUN apk add upx

# Set working directory
WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

# Get dependencies
RUN go mod download

# Copy source files
COPY . .

FROM base as dev

RUN echo 'go run /usr/src/app/cmd/main.go' >> /start_app

RUN chmod +x /start_app

CMD ["/start_app"]

FROM base as builder

# Set working directory
WORKDIR /usr/src/app

# Build binary
RUN go build -ldflags="-s -w" -o dist/main cmd/main.go

# Compress binary
RUN upx dist/main

FROM alpine:3

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/dist/main main

EXPOSE 4000

CMD [ "./main" ]