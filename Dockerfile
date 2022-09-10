FROM alpine:3.16 as builder
ADD .build/webhook-conversion-service /usr/bin/
RUN chmod +x /usr/bin/webhook-conversion-service
RUN apk add ca-certificates


FROM scratch
COPY --from=builder /usr/bin/webhook-conversion-service /usr/bin/
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /etc/ca-certificates.conf /etc/ca-certificates.conf
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

EXPOSE 8080
USER 10080

ENTRYPOINT ["/usr/bin/webhook-conversion-service"]
