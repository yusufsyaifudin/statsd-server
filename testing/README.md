# Testing Using Datadog Statsd Agent

## Run Datadog

```shell
docker run --env-file .env --cgroupns host --pid host -p 8125:8125/udp datadog/dogstatsd:7
```