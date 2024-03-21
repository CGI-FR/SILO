FROM gcr.io/distroless/base
ARG BIN
COPY /bin/silo /silo
CMD ["/silo"]
