FROM 644152709166.dkr.ecr.eu-west-1.amazonaws.com/usabilla/dev/dockerhub-mirror/alpine:latest

COPY rds_exporter  /bin/
# COPY config.yml           /etc/rds_exporter/config.yml

RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates

EXPOSE      9042
ENTRYPOINT  [ "/bin/rds_exporter", "--config.file=/etc/rds_exporter/config.yml" ]
