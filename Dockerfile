FROM gportal/pipeline:latest AS protobuf

COPY ./ /opt/app-root/src

WORKDIR /opt/app-root/src

RUN make proto

FROM gportal/golang:latest as builder

# Import application source
COPY ./ /opt/app-root/src
COPY --from=protobuf /opt/app-root/src/pkg/proto/ /opt/app-root/src/pkg/proto/

# Change working directory
WORKDIR /opt/app-root/src

# Build binary for Latency Service
RUN go build -v -o "${APP_ROOT}/metadata-server" cmd/server.go

FROM gportal/golang:latest

# Import application source
COPY --from=builder /opt/app-root/metadata-server /opt/app-root/metadata-server

RUN setcap cap_net_bind_service+ep "${APP_ROOT}/metadata-server"

EXPOSE 80/tcp

ENV METADATA_SERVER_CONFIG=/data/config.yaml

VOLUME /data

RUN /usr/bin/fix-permissions ${APP_ROOT}

CMD ["/opt/app-root/metadata-server"]
