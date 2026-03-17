#build stage
FROM golang:1.26.1@sha256:e2ddb153f786ee6210bf8c40f7f35490b3ff7d38be70d1a0d358ba64225f6428 AS builder

RUN apt-get update && apt-get install -y \
    libzmq3-dev \
    pkg-config \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN chmod +x ./scripts/build.sh

ARG CI=false
ARG VERSION=""
ARG COMMIT=""

RUN ./scripts/build.sh \
    $( [ "$CI" = "true" ] && echo "--ci" ) \
    $( [ -n "$VERSION" ] && echo "--version \"$VERSION\"" ) \
    $( [ -n "$COMMIT" ] && echo "--commit \"$COMMIT\"" )


#runtime stage
FROM debian:bookworm-slim@sha256:f06537653ac770703bc45b4b113475bd402f451e85223f0f2837acbf89ab020a AS build-release-stage

RUN apt-get update && apt-get install -y \
    libzmq5 \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -u 1001 appuser && \
    mkdir -p /data && \
    mkdir -p /home/appuser/.local/share/cogmoteGO/experiments && \
    chown -R appuser:appuser /data /home/appuser/.local

ENV XDG_DATA_HOME=/home/appuser/.local/share

COPY --from=builder --chown=appuser:appuser /app/build/linux-amd64/cogmoteGO /usr/local/bin/

EXPOSE 9012

WORKDIR /data
VOLUME ["/data", "/home/appuser/.local/share"]

USER appuser

ENTRYPOINT ["cogmoteGO"]