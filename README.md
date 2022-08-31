webhook-conversion-service
==========================

Fixes up your webhooks. Impossible becomes possible ;)

Created originally for purpose to fix invalid URL addresses injected by Gitlab in webhook body, so the ArgoCD was not able to recognize from which repository the webhook comes from.
The case was in communication between Gitlab and ArgoCD over the internal network, when both services are behind an internal network.

**Features:**
- Replacing text in body of a request and response
- Multiple webhooks/endpoints support
- Configuration in a YAML file
