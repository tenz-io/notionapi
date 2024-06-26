module notionapi-example

go 1.20

require (
	github.com/joho/godotenv v1.5.1
	github.com/tenz-io/gokit/httpcli v1.5.1
	github.com/tenz-io/gokit/logger v1.5.0
	github.com/tenz-io/notionapi v0.0.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/tenz-io/gokit/monitor v1.5.0 // indirect
	github.com/tenz-io/gokit/tracer v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/tenz-io/notionapi => ./..
