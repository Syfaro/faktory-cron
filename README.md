# faktory-cron

A very basic scheduler for running Faktory jobs.

```yaml
faktory:
  url: tcp://faktory-host:7419
jobs:
  - name: every minute
    every: "* * * * *"
    job_type: job_name
    queue: optional_queue
    args:
      - arg1
      - arg2
    custom:
      key: value
```
