FROM golang:alpine AS builder

COPY . /app
WORKDIR /app

# Install deps
RUN apk --update add make ca-certificates

# Build again in case anything has changed
RUN make build-release

FROM scratch

# Copy the built static binary from the intermediate image
COPY --from=builder /app/Neo /Neo

CMD ["/Neo"]
