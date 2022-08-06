FROM golang:1.17-alpine AS gobuild

ENV GO111MODULE on
ENV GOSUMDB off
# add go-base repo to exceptions as a private repository.
ENV GOPRIVATE $GOPRIVATE,github.com/cryptowize-tech

# add private github token
ARG GITHUB_TOKEN
RUN apk add --no-cache bash git openssh build-base
RUN if [ -z "$GITHUB_TOKEN"  ] ; then \
    echo 'GITHUB_TOKEN not provided, please use docker build --build-arg GITHUB_TOKEN="xxxx"' \
    ; else git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/" \
    ; fi

WORKDIR /src

# Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever).
COPY go.* ./

COPY . .

# Compile! Should only compile our sources since everything else is precompiled.
ARG RACE=-race
ARG CGO=1
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /src/bin && \
    GOOS=linux CGO_ENABLED=${CGO} go build ${RACE} -v -installsuffix cgo -o ./bin/api -ldflags "-linkmode external -extldflags -static -s -w" ./cmd/api

FROM scratch

# Import the user and group files from the build stage.
#COPY --from=gobuild /etc/group /etc/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ROOT /opt/appworker
ENV PATH /opt/appworker

COPY --from=gobuild /src/bin $APP_ROOT

EXPOSE 8080
USER appworker
CMD ["/opt/appworker/api"]