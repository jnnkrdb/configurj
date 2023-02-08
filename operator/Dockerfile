# -----------------------------------------------------
# BUILD STAGE
FROM golang:1.19-alpine AS build
# Create Path-Structure
RUN mkdir -p /build/engine
# Copy Sources
COPY engine /build/engine
# Setting Workdir for BUILD-Process
WORKDIR /build/engine
# SET ENVs for GO BUILD
ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux
# Get GO-Packages
RUN go mod download
# Build the binary
RUN go build -o go_configurj .
# -----------------------------------------------------
# FINAL STAGE
FROM alpine:3.16
RUN mkdir /config
RUN mkdir /app
COPY --from=build /build/engine/go_configurj /app
ENTRYPOINT [ "/app/go_configurj" ]