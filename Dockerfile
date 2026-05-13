#build stage
FROM golang:1.26.3@sha256:313faae491b410a35402c05d35e7518ae99103d957308e940e1ae2cfa0aac29b AS builder

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
FROM debian:bookworm-slim@sha256:f9c6a2fd2ddbc23e336b6257a5245e31f996953ef06cd13a59fa0a1df2d5c252 AS build-release-stage

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