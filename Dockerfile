FROM scratch
ARG TARGETARCH
ADD /linux/${TARGETARCH} /
ENTRYPOINT ["/demodata_service"]
