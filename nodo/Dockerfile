FROM golang:alpine as builder

# Set necessary environmet variables needed for our image
ENV  \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
ENV PORT=$PORT
# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
EXPOSE $PORT

# Command to run when starting the container
CMD ["/dist/main"]