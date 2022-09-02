webhook-conversion-service
==========================

Fixes up your webhooks. Impossible becomes possible ;)

Created originally for purpose to fix invalid URL addresses injected by Gitlab in webhook body. ArgoCD was not able to recognize from which repository the webhook came from.
The case was in communication between Gitlab and ArgoCD over the internal network, when both services are behind a VPN.

**Features:**
- Replacing text in body of a request and response
- Multiple webhooks/endpoints support
- Configuration in a YAML file
- Non-root, distroless container image

Usage
-----

```bash
webhook-conversion-service --config ./example-config.yaml --listen ":8080"
```

Security
--------

This is a reverse proxy. All headers, body and query string are passed to the upstream, which means the upstream could be manipulated in those ways.
Please be aware of a fact, that potential attacker could access your upstream service using your REVERSE PROXY IP.

**Advices:**
- Control who has access to the service using firewall rules or ingress network policy in Kubernetes
- If its possible use only for internal services
