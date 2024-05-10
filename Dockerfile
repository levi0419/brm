FROM golang:1.19.0 AS builder

# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build -o app

FROM alpine:latest AS production
# We have to copy the output from our
# builder stage to our production stage
COPY --from=builder app .
COPY db/migration ./db/migration

EXPOSE 7070

CMD ./app

