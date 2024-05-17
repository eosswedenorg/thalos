FROM alpine:latest
LABEL maintainer="Henrik Hautakoski <henrik.Hautakoski@gmail.com>"
ARG VERSION=1.1.2-rc2
WORKDIR /thalos
ADD --chmod=755 https://github.com/eosswedenorg/thalos/releases/download/v$VERSION/thalos-server-${VERSION}-linux-amd64-musl thalos-server
ENTRYPOINT [ "./thalos-server" ]
