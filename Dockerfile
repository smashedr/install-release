FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/install-release"]
COPY $TARGETPLATFORM/install-release /usr/bin/
