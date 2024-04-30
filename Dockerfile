# ----------- STEP 1: Build executable -----------
    FROM golang:alpine AS builder

    # Install git
    RUN apk update && apk add --no-cache git tzdata ca-certificates && update-ca-certificates
    
    WORKDIR /app
    COPY . .
    
    # Fetch dependencies
    RUN go get -d -v
    # Build executable
    RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main main.go
    
    # ----------- STEP 2: Build small image ----------- 
    FROM scratch
    WORKDIR /app
    
    # Import user files, certificates, and timezone info
    COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
    COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
    COPY --from=builder /etc/passwd /etc/passwd
    COPY --from=builder /etc/group /etc/group
    # Import and base config executable
    COPY .env.example /app/.env
    COPY --from=builder /app/bin/main /app/bin/main
    # Set timezone
    ENV TZ=Asia/Jakarta

    # Expose port
    EXPOSE 8080

    # Run executable
    ENTRYPOINT [ "/app/bin/main" ]