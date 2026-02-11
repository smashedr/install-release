FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/ir"]
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine:latest /tmp /tmp
COPY $TARGETPLATFORM/ir /usr/bin/
ENV DOCKER=true
