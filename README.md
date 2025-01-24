# Golang JSON objects Vault

Golang microservice to store JSON entities and access it

The service is a JSON object store with an HTTP interface. Saved objects are placed in RAM, it is possible to set the retention time to the object.
Storage enriches by source file while start-up and save current state before shutting down.


## API

Service API presented in `docs` folder and includes such requests:

### Objects

```
PUT /objects/{uid} - to save new obj
```

```
GET /objects/{uid} - to get existing obj
```

### Metrics

Service allows to receive Prometheus metrics

```
GET /metrics - proetheus metrics
```

### Kubernetes

```
GET /probes/readiness - check service readiness
```

```
GET /probes/liveness - check service liveness
```

### Docker

Service allows to use Dockerfile:

```
docker build -t json-vault .
```

```
docker run -d -p 8081:8081 json-vault
```

@Author Andrei Gerasimov
@2025
@end
1