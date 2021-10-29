## jaeger-k8s deploy

### Prerequisites 

* Kubernetes cluster ready
* Jaeger architecture understood
* App for test
* All yaml config files

### Config files
all config file could be found in [jaeger-k8s](../.deploy/jaeger-k8s)

```shell script
.
├── cassandra-job.yml
├── cassandra-service.yml
├── cassandra-stateful-set.yml
├── configmap.yml
├── jaeger-production-agent-ds.yml
├── jaeger-production-collector-deploy.yml
├── jaeger-production-collector-svc.yml
├── jaeger-production-collector-zipkin-svc.yml
├── jaeger-production-query-deploy.yml
├── jaeger-production-query-svc.yml
└── prometheus
    ├── Kubernetes-Pod-Resources.json
    ├── grafana.yaml
    ├── namespace.yaml
    └── prometheus.yaml
```

### Deploy in sequence

1. cassandra deploy: [ElasticSearch is also supported] 

    * apply `cassandra-stateful-set.yml` `cassandra-service.yml`
    
2. cassandra schema job to create `jaeger` tables:

    * apply `cassandra-job.yml`
    
3. deploy jaeger-collector:
    
    * apply `jaeger-production-collector-deploy.yml`
    * apply `jaeger-production-collector-svc.yml`
    * apply `jaeger-production-collector-zipkin-svc.yml` [OPTIONAL for zipkin]
    
    **NOTICE**: 

    * jaeger has breaking changes at `v1.18`, you need to read.
    
4. deploy the jaeger-agent: [OPTIONAL, if you don't know why, you should read jaeger document in detail]
    
    * apply `jaeger-production-agent-ds.yml`
    
    **NOTICE**: **jaeger-client** should use `6381/UDP` to report spans to jaeger-agent.
    
5. deploy the jaeger-query
    
    * apply `jaeger-production-query-deploy.yml`
    * apply `jaeger-production-query-svc.yml`
    
    **NOTICE**: jaeger-query would log span to jaeger-agent, default `localhost:6381`, so there's error log in `query`.
    you can modify the `flag args` or `ENV` of jaeger-query to let this work.
    
6. monitor jaeger
    
    have a read at first: [monitor jaeger in v1.18](https://www.jaegertracing.io/docs/1.18/monitoring/)
    
    * deploy prometheus and grafana
    * prepare your dashboard
    
    