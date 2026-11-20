
# ==========================================
# STAGE 1: BUILD THE GO APPLICATION
# ==========================================

FROM golang:1.26-alpine AS builder

# Set the working directory inside the builder
WORKDIR /src

# Copy the Go application source code
COPY app/ ./

# Compile the application into a Linux executable
RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/quicknotes .

# ==========================================
# STAGE 2: CREATE THE LIGHTWEIGHT IMAGE
# ==========================================

FROM alpine:3.22

# Create a non-root user for security
RUN addgroup -S app && adduser -S -G app app

# Copy only the compiled executable
COPY --from=builder /out/quicknotes /usr/local/bin/quicknotes

# Run the application as a non-root user
USER app

# Document the application port
EXPOSE 8080

# Start the API
ENTRYPOINT ["/usr/local/bin/quicknotes"]
