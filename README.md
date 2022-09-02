webhook-conversion-service
==========================

Fixes up your webhooks. Impossible becomes possible ;)

Created originally for purpose to fix invalid URL addresses injected by Gitlab in webhook body. ArgoCD was not able to recognize from which repository the webhook came from.
The case was in communication between Gitlab and ArgoCD over the internal network, when both services are behind a VPN.

**Features:**
- Replacing text in body of a request and response
- Multiple webhooks/endpoints support
- Configuration in a YAML file

Usage
-----

```bash
webhook-conversion-service --config ./example-config.yaml --listen ":8080"
```
