# Stage 1: Build the runtime binary
FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY ./runtime /build/runtime

RUN cd runtime && go build -o /build/daxagrid-runtime main.go

# Stage 2: Minimal final image
FROM alpine:3.18

WORKDIR /app

# Copy compiled binary
COPY --from=builder /build/daxagrid-runtime /usr/bin/daxagrid-runtime

# Expose public function port and internal ops port
EXPOSE 8080 36365

# Run the runtime directly
CMD ["/usr/bin/daxagrid-runtime"]
