# -------------------------
# 1) Builder Stage
# -------------------------
    FROM golang:1.23 AS builder

    # Label the stage and maintainer
    LABEL maintainer="ahmadsa@stud.ntnu.no"
    LABEL stage="builder"
    
    # Set the working directory inside the container where all commands will be run.   
     WORKDIR /go/src/assignment-2
    
    # Copy go.mod and go.sum first to leverage caching
    COPY go.mod go.sum /go/src/assignment-2/
    RUN go mod download
    
    # Copy the remaining project files
    COPY cmd ./cmd
    COPY constants ./constants
    COPY firebase ./firebase
    COPY handlers ./handlers
    COPY mock_data ./mock_data
    COPY services ./services
    COPY structs ./structs
    COPY tools ./tools
    
    # Build a static Go binary named "assignment-2"
    RUN CGO_ENABLED=0 GOOS=linux go build -o assignment-2 ./cmd
    
    # -------------------------
    # 2) Final Stage
    # -------------------------
    FROM alpine:latest
    
    # Install dependencies such as tzdata (for time zone support) and curl (for healthchecks)
    RUN apk add --no-cache tzdata curl
    
    # Copy only the built binary from the builder stage
    COPY --from=builder /go/src/assignment-2/assignment-2 .
    
    COPY assignment-2-firebasekey.json .
    
    # Expose the port your application listens on (e.g., 8080)
    EXPOSE 8080
    
    # Define a health check if you want Docker to verify your service is alive
    HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
      CMD ["curl", "-f", "http://localhost:8080/dashboard/v1/status/"]
    
    # Set the container's entrypoint to the compiled binary
    ENTRYPOINT ["./assignment-2"]
    