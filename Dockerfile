FROM scratch
ARG PORT=8080
COPY ./esmd /go/bin/esmd
COPY ./esmd.yml /go/bin/esmd.yml
ENV ESM_LISTEN ${PORT}
EXPOSE ${PORT}
CMD ["/go/bin/esmd", "-c esmd.yml"]
