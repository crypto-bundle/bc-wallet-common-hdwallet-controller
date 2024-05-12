FROM golang:1.22.2-alpine AS gobuild

ENV GO111MODULE on
ENV GOSUMDB off
# add go-base repo to exceptions as a private repository.
ENV GOPRIVATE $GOPRIVATE,github.com/crypto-bundle

# add private github token
RUN apk add --no-cache git openssh build-base && \
    mkdir -p -m 0700 ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    git config --global url."git@github.com".insteadOf "https://github.com/"

WORKDIR /src

# Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever).
COPY go.* ./

COPY . .

# Compile! Should only compile our sources since everything else is precompiled.
ARG RACE=-race
ARG CGO=1
ARG RELEASE_TAG="v0.0.0-00000000-100500"
ARG COMMIT_ID="0000000000000000000000000000000000000000"
ARG SHORT_COMMIT_ID="00000000"
ARG BUILD_NUMBER="100500"
ARG BUILD_DATE_TS="1713280105"
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=ssh \
    mkdir -p /src/bin && \
    GOOS=linux CGO_ENABLED=${CGO} go build ${RACE} -trimpath -race -installsuffix cgo \
        -gcflags all=-N \
        -o ./bin/api \
        -ldflags "-linkmode external -extldflags -static -s -w \
                -X 'main.BuildDateTS=${BUILD_DATE_TS}' \
                -X 'main.BuildNumber=${BUILD_NUMBER}' \
                -X 'main.ReleaseTag=${RELEASE_TAG}' \
                -X 'main.CommitID=${COMMIT_ID}' \
                -X 'main.ShortCommitID=${SHORT_COMMIT_ID}'" \
        ./cmd/api

FROM scratch

# Import the user and group files from the build stage.
#COPY --from=gobuild /etc/group /etc/passwd /etc/

ENV APP_ROOT /opt/appworker
ENV PATH /opt/appworker

COPY --from=gobuild /src/bin $APP_ROOT

EXPOSE 8080
USER appworker
CMD ["/opt/appworker/api"]