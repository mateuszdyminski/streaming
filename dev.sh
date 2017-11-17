#!/bin/bash

usage() {
	cat <<EOF
Usage: $(basename $0) <command>

Wrappers around core binaries:
    run                    Runs the streaming in development mode.
    build                  Builds the streaming.
    buildPush              Builds docker image and pushes it to DockerHub.
EOF
	exit 1
}

VERSION=$(git describe --always)
LAST_COMMIT_USER="$(tr -d '[:space:]' <<<"$(git log -1 --format=%cn)<$(git log -1 --format=%ce)>")"
LAST_COMMIT_HASH=$(git log -1 --format=%H)
LAST_COMMIT_TIME=$(git log -1 --format=%cd --date=format:'%Y-%m-%d_%H:%M:%S')

BINARY=streaming-$VERSION

CMD="$1"
shift
case "$CMD" in
	build)
	    LDFLAGS="-X main.appVersion=$VERSION -X main.lastCommitTime=$LAST_COMMIT_TIME -X main.lastCommitHash=$LAST_COMMIT_HASH -X main.lastCommitUser=$LAST_COMMIT_USER -X main.buildTime=$(date -u +%Y-%m-%d_%H:%M:%S)"
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/$BINARY -a -tags netgo .
	;;
	run)
		LDFLAGS="-X main.appVersion=$VERSION -X main.lastCommitTime=$LAST_COMMIT_TIME -X main.lastCommitHash=$LAST_COMMIT_HASH -X main.lastCommitUser=$LAST_COMMIT_USER -X main.buildTime=$(date -u +%Y-%m-%d_%H:%M:%S)"
		go run -ldflags "$LDFLAGS" main.go --config=config/local.toml
	;;
	buildPush)
		docker build -t mateuszdyminski/streaming:latest -t mateuszdyminski/streaming:$VERSION --build-arg APP_VERSION=$VERSION . && docker push mateuszdyminski/streaming
	;;
	*)
		usage
	;;
esac