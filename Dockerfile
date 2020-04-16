FROM golang:1.14.2-alpine
RUN mkdir /app
WORKDIR /app
EXPOSE 8090
COPY esmd .
COPY esmd.yml .
CMD ["./esmd", "-l localhost:8090", "-c esmd.yml"]
