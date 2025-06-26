# alertmanager-webhook-signal

[![Integration](https://github.com/systemli/alertmanager-webhook-signal/actions/workflows/integration.yml/badge.svg)](https://github.com/systemli/alertmanager-webhook-signal/actions/workflows/integration.yml) [![Quality](https://github.com/systemli/alertmanager-webhook-signal/actions/workflows/quality.yml/badge.svg)](https://github.com/systemli/alertmanager-webhook-signal/actions/workflows/quality.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=systemli_alertmanager-webhook-signal&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=systemli_alertmanager-webhook-signal) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=systemli_alertmanager-webhook-signal&metric=coverage)](https://sonarcloud.io/summary/new_code?id=systemli_alertmanager-webhook-signal) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=systemli_alertmanager-webhook-signal&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=systemli_alertmanager-webhook-signal)

This service listens for [webhook requests by Alertmanager](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config) and forwards the alerts to a Signal group.

It requires the [JSON-RPC service of AsamK/signal-cli](https://github.com/AsamK/signal-cli/wiki/JSON-RPC-service) to send messages to Signal.

## Configuration

The service expects several environment variables to be set. See `.env.example`.

### Alertmanager configuration

Example configuration to use this service as a [webhook receiver](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config) in Alertmanager that receives alerts with severity `critical` (default receiver `admins` receives all alerts):

```yaml
receivers:
- name: admins_mail
  [...]
- name: admins_signal
  webhook_configs:
  - url: http://container:8080/alertmanager
    send_resolved: true

route:
  group_wait: 1m
  group_interval: 5m
  repeat_interval: 4h
  receiver: admins_mail
  routes:
  - match:
      severity: critical
    receiver: admins_signal
    continue: true
  - match:
      severity: critical
    receiver: admins_mail
    continue: true
```
