# This dockerfile can be used to build a binary for use with the QuickPiperAudiobook command. 
# You can use it for testing, or other architectures that don't have a piper build.
# docker build -t quickpiperaudiobook .
# docker run quickpiperaudiobook /app/examples/lorem_ipsum.txt

FROM golang:1.22 AS build

WORKDIR /app

COPY . .

# Install Go dependencies and build the binary
RUN go mod tidy && \
    go build -o QuickPiperAudiobook .

FROM ubuntu:latest

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    # needed for ebook-convert
    calibre \
    ffmpeg \
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
