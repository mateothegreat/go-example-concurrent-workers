# Concurrent Worker Scheduling

This is a simple example of concurrent worker scheduling using goroutines and channels.

## Usage

This scheduler is configured to run 10 jobs concurrently with a total of 10 workers.

See [main.go#L12-L17](main.go#L14-L19) for the settings:

```go
var (
  // maxConcurrent is the maximum number of concurrent workers.
  maxConcurrent = 10
  // totalJobs is the total number of jobs to be processed.
  totalJobs = 10
)
```

See [main.go#L30-L35](main.go#L30-L33) for where you can add your business logic.

## Example

```bash
$ go run main.go
2024/08/08 18:50:23 Worker 7 started, waiting for jobs...
2024/08/08 18:50:23 Worker 5 started, waiting for jobs...
2024/08/08 18:50:23 Worker 0 started, waiting for jobs...
2024/08/08 18:50:23 Worker 8 started, waiting for jobs...
2024/08/08 18:50:23 Worker 2 started, waiting for jobs...
2024/08/08 18:50:23 Worker 1 started, waiting for jobs...
2024/08/08 18:50:23 Worker 4 started, waiting for jobs...
2024/08/08 18:50:23 Worker 3 started, waiting for jobs...
2024/08/08 18:50:23 Worker 9 started, waiting for jobs...
2024/08/08 18:50:23 Worker 6 started, waiting for jobs...
2024/08/08 18:50:23 Result received: Worker 0 processed job: 18:50:23
2024/08/08 18:50:23 Result received: Worker 2 processed job: 18:50:23
2024/08/08 18:50:23 Result received: Worker 5 processed job: 18:50:23
2024/08/08 18:50:23 Result received: Worker 7 processed job: 18:50:23
2024/08/08 18:50:24 Result received: Worker 8 processed job: 18:50:23
2024/08/08 18:50:25 Result received: Worker 9 processed job: 18:50:23
2024/08/08 18:50:25 Result received: Worker 4 processed job: 18:50:23
2024/08/08 18:50:25 Result received: Worker 3 processed job: 18:50:23
2024/08/08 18:50:25 Result received: Worker 1 processed job: 18:50:23
2024/08/08 18:50:26 All jobs processed. Shutting down...
```
