## This dockerfile is primarily for testing. It can be used like the following:
# docker build -t quickpiperaudiobook .
# docker run quickpiperaudiobook /app/examples/lorem_ipsum.txt

FROM --platform=linux/amd64 golang:latest as build

WORKDIR /app

# Copy all the code from the current directory
COPY . .

# Install Go dependencies and build the binary
RUN go mod tidy && \
    go build -o QuickPiperAudiobook .

# Final stage
FROM --platform=linux/amd64 ubuntu:latest

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    # needed for ebook-convert
    calibre \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary
COPY --from=build /app/QuickPiperAudiobook .

# Copy any additional files needed at runtime
COPY --from=build /app/examples /app/examples

# Set the command to run the binary, allowing an argument for the file
ENTRYPOINT ["./QuickPiperAudiobook"]
# Allow passing arguments from CLI; mount a file if needed, or point to a remote URL
CMD [] 
