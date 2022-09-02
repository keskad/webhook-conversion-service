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
- Stateless service, could be running easily in HA mode by simply bumping the replicas count

Usage
-----

```bash
webhook-conversion-service --config ./example-config.yaml --listen ":8080"
```

```yaml
endpoints:
    # Example for ArgoCD
    - path: /api/webhook
      targetUrl: http://argocd-server.argocd.svc.cluster.local
      replacements:
          - from: gitlab.example.org
            to: my-instance.gitlab.svc.cluster.local

    # just for testing
    - path: /test
      targetUrl: https://google.com/search
      replacements:
          - from: google
            to: DuckDuckLetsGo
          - from: Google
            to: DuckDuckLetsGo
```

Resources requirements
----------------------

The reverse proxy is very lightweight, probably after some time you would forget that it exists at all.

```yaml
limits:
    cpu: 100m
    memory: 32Mi
requests:
    cpu: 50m
    memory: 16Mi
```

Security
--------

This is a reverse proxy. All headers, body and query string are passed to the upstream, which means the upstream could be manipulated in those ways.
Please be aware of a fact, that potential attacker could access your upstream service using your REVERSE PROXY IP.

**Advices:**
- Control who has access to the service using firewall rules or ingress network policy in Kubernetes
- If its possible use only for internal services

**Example Network Policy:**

```yaml
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-only-gitlab-to-argocd
  namespace: argocd
spec:
    podSelector:
        matchLabels:
            app.kubernetes.io/instance: webhooks-conversion
            app.kubernetes.io/name: webhooks-conversion
    policyTypes:
        - Ingress
        - Egress
    ingress:
        - namespaceSelector:
              matchLabels:
                  name: gitlab
          ports:
              - protocol: TCP
                port: 8080
    egress:
        - to:
            - namespaceSelector:
                  matchLabels:
                      name: argocd
          ports:
              - port: 80
                protocol: TCP
```
