# syntax = docker/dockerfile:experimental
FROM golang:1.19.2 AS build
WORKDIR /src
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /movie-backend ./app/*.go

FROM scratch AS linux-amd64
COPY --from=build /movie-backend /movie-backend-linux-amd64

FROM scratch AS darwin-amd64
COPY --from=build /movie-backend /movie-backend-darwin-amd64

FROM scratch AS windows-amd64
COPY --from=build /movie-backend /movie-backend.exe