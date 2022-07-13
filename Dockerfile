# -----------------------------------------------------
# BUILD STAGE
FROM golang:1.18-alpine AS build
# Create Path-Structure
RUN mkdir -p /build/handler
RUN mkdir /build/probes
# Copy Sources
COPY go.mod /build/
COPY go.sum /build/
COPY main.go /build/
COPY probes/ /build/probes/
COPY handler/ /build/handler/
# Setting Workdir for BUILD-Process
WORKDIR /build
# Get GO-Packages
RUN go get github.com/gin-contrib/cors
RUN go get github.com/gin-gonic/gin
RUN go get k8s.io/client-go/kubernetes
RUN go get k8s.io/client-go/rest
RUN go get k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get k8s.io/api/core/v1
RUN go get k8s.io/client-go/tools/clientcmd
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o go_configurj .
# -----------------------------------------------------
# FINAL STAGE
FROM alpine:3.16
RUN mkdir /configs
RUN mkdir /app
COPY --from=build /build/go_configurj /app
ENTRYPOINT [ "/app/go_configurj" ]