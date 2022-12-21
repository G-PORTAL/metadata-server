FROM gportal/golang:latest

# Import application source
COPY ./ /opt/app-root/src

# Change working directory
WORKDIR /opt/app-root/src

# Build binary for Latency Service
RUN go build -v -o "${APP_ROOT}/metadata-server" cmd/server.go && \
    setcap cap_net_bind_service+ep "${APP_ROOT}/metadata-server"

# Finally delete application source
RUN rm -rf /opt/app-root/src/*

EXPOSE 80/tcp

VOLUME /config

RUN /usr/bin/fix-permissions ${APP_ROOT}

CMD ["/opt/app-root/metadata-server"]