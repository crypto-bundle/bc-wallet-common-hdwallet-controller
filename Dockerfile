FROM golang:1.19-alpine AS gobuild

ENV GO111MODULE on
ENV GOSUMDB off
# add go-base repo to exceptions as a private repository.
ENV GOPRIVATE $GOPRIVATE,github.com/crypto-bundle

# add private github token
RUN apk add --no-cache git openssh build-base && \
    mkdir -p -m 0700 ~/.ssh && \
    ssh-keyscan gitlab.heronodes.io >> ~/.ssh/known_hosts && \
    git config --global url."git@github.com".insteadOf "https://github.com/"

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

ENV APP_ROOT /opt/appworker
ENV PATH /opt/appworker

COPY --from=gobuild /src/bin $APP_ROOT

EXPOSE 8080
USER appworker
CMD ["/opt/appworker/api"]