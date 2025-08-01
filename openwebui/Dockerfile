FROM registry.suse.com/bci/golang:1.24 as build

WORKDIR /app

COPY openwebui-client.go ./
RUN go mod init openwebui-client
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /openwebui-client ./openwebui-client.go

# Bundle the application into a scratch image
FROM scratch

COPY --from=build /openwebui-client /bin/openwebui-client

CMD ["/bin/openwebui-client"]
