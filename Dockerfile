FROM scratch

ADD .build/webhook-conversion-service /usr/bin/

EXPOSE 8080
USER 10080

ENTRYPOINT ["webhook-conversion-service"]
