#build stage
FROM golang:1.26.5@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS builder

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
FROM debian:bookworm-slim@sha256:60eac759739651111db372c07be67863818726f754804b8707c90979bda511df AS build-release-stage

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