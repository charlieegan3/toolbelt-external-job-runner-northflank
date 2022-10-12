# toolbelt-external-job-runner-northflank

This is an extension for [toolbelt](https://github.com/charlieegan3/toolbelt) which provides the running of external
jobs in a northflank project.

Example config to run `cmd/run.go`:

```go
northflank:
  token: xxx.xxx...
job:
  project_id: example
  job_id: job-1
  command: echo 1
  env:
	KEY: VALUE  
      
```