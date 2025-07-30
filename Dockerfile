# 1) Builder stage
FROM golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy EVERYTHING (including your generated stubs under api/ and controllers/)
COPY . .

# Build the operator binary
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS:-linux} \
    GOARCH=${TARGETARCH} \
    go build -a -o manager cmd/main.go

# 2) Final distroless image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
