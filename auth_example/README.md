https://medium.com/opentelemetry/securing-your-opentelemetry-collector-1a4f9fa5bd6f

## Generate SSL Keys 

```bash
$ cfssl genkey -initca ca-csr.json | cfssljson -bare ca
$ cfssl gencert -ca ca.pem -ca-key ca-key.pem cert-csr.json | cfssljson -bare cert
```

## Start Cloak OAuth Server 

```
docker run -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:11.0.2
```
- Generate tokens (following instructions) from above article

## Start Agent and Collector 

```bash
otelcol --config config.collector.yaml --log-level=debug
```

```bash
otelcol --config config.agent.yaml --log-level=debug --metrics-addr=:8898
```







