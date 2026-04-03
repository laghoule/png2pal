##############################
FROM golang:1.26-alpine AS build

ARG VERSION "devel"
ARG GIT_COMMIT ""

WORKDIR /src

RUN --mount=type=bind,source=.,target=.  \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.version=$VERSION' -X 'main.gitCommit=$GIT_COMMIT'" -o /tmp/png2pal main.go

##############################
FROM scratch

ARG VERSION

LABEL org.opencontainers.image.title="png2pal" \
  org.opencontainers.image.vendor="laghoule" \
  org.opencontainers.image.licenses="GPLv3" \
  org.opencontainers.image.version="${VERSION}" \
  org.opencontainers.image.description="Converts paletted PNG images to binary tileset format. Extracts tiles with configurable dimensions and spacing." \
  org.opencontainers.image.url="https://github.com/laghoule/png2pal/README.md" \
  org.opencontainers.image.source="https://github.com/laghoule/png2pal" \
  org.opencontainers.image.documentation="https://github.com/laghoule/png2pal/README.md"

USER 65534

COPY --link --from=build /tmp/png2pal /bin/png2pal

ENTRYPOINT ["/bin/png2pal"]
