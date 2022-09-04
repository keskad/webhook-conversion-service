FROM alpine:3.6 as builder
ADD .build/webhook-conversion-service /usr/bin/
RUN chmod +x /usr/bin/webhook-conversion-service


FROM scratch
COPY --from=builder /usr/bin/webhook-conversion-service /usr/bin/

EXPOSE 8080
USER 10080

ENTRYPOINT ["/usr/bin/webhook-conversion-service"]
