### Getting started:
* Run the container with `docker compose up`.
* Start the project with `go run app/main.go` and see the results in the console.
* Benchmarking the project with:
```bash
make bench
```

### Structure:
1.  `app/` - place for handlers,tests and the `main.go` file;
2. `business/` - place for layers of our business logic
   1. `core/` - business logic and models for every use case divided in layers (`db` - `store` -`core`)
   2. `data/` - all initial schemas and initDB function
   3. `sys/` - initialization of our database + CRUD
   4. `web/` - place for middlewares like `cashe`

### Benchmark results:
```bash
BenchmarkSourceCampaigns-8   10000	132289 ns/op	0.00549 avg_response_time_ms	30.4 max_response_time_ms	0.00167 min_response_time_ms
PASS
```
It's easy to understand that initially the result has a bigger response time, due to the result not cashed. After cashing, the handling gives much better results.

Thanks for the review!