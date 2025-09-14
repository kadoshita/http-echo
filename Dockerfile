# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Build stage
FROM golang:1.21-alpine AS builder

# TARGETOS and TARGETARCH are set automatically when --platform is provided.
ARG TARGETOS
ARG TARGETARCH
ARG PRODUCT_VERSION
ARG BIN_NAME=http-echo

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -a \
    -o=${BIN_NAME} \
    -ldflags="-s -w -X 'github.com/kadoshita/http-echo/version.Version=${PRODUCT_VERSION}'" \
    -trimpath \
    -buildvcs=false \
    .

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot as default

# TARGETOS and TARGETARCH are set automatically when --platform is provided.
ARG TARGETOS
ARG TARGETARCH
ARG PRODUCT_VERSION
ARG BIN_NAME
ENV PRODUCT_NAME=$BIN_NAME

LABEL name="http-echo" \
      maintainer="kadoshita" \
      vendor="kadoshita" \
      version=$PRODUCT_VERSION \
      release=$PRODUCT_VERSION \
      licenses="MPL-2.0" \
      summary="A test webserver that echos a response. You know, for kids." \
      org.opencontainers.image.title="http-echo" \
      org.opencontainers.image.description="A test webserver that echos a response. You know, for kids." \
      org.opencontainers.image.source="https://github.com/kadoshita/http-echo" \
      org.opencontainers.image.url="https://github.com/kadoshita/http-echo" \
      org.opencontainers.image.documentation="https://github.com/kadoshita/http-echo" \
      org.opencontainers.image.vendor="kadoshita" \
      org.opencontainers.image.licenses="MPL-2.0" \
      org.opencontainers.image.version=$PRODUCT_VERSION

COPY --from=builder /app/$BIN_NAME /http-echo
COPY LICENSE /usr/share/doc/$PRODUCT_NAME/LICENSE.txt

EXPOSE 5678/tcp

ENV ECHO_TEXT="hello-world"

ENTRYPOINT ["/http-echo"]
