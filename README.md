# superman-detector

## Considerations/Improvements
- An openapi spec is provided to generate client libraries against the API. Ideally behavioral/integration level tests
would be added and run against a running API. Tests should also be validating the response structure that is in the expected
format in openapi.

- Abstracting out the ORM. Gorm was used to communicate to the database and while this is passed via composition to
the login handler, the queries are highly coupled to GORM. This should be abstracted out further, since a true dev/prod
environment would be using different RDBS (Postgres, Mysql) or NoSql alternatives (Dynamo, Datastore).

- Radius uncertainty, was not implemented but would have most likely added additional field related to
Propagation of Error.

- Caching, although the application is measuring the percentile performance of runtime geo ip lookups and
travel distance. These could be possibly be cached using memoization(increased memory usage) or a caching layer
LRU (redis). GeoIP rows could also have been inserted into an already created sqlite database as well.


## Running
As [12 Factor](https://12factor.net/) suggests configuration should be stored in the
environment. For local purposes and running in an orchestrator setting sensitive information
can be provided in other manners (volumeMounts, secrets (k8, aws ssm, aws sm), etc). Below
are settings that configure the runtime and defaults.

[Environment Variable Settings/Lookup](./internal/pkg/settings/settings.go)

```bash
DETECTOR_API_SERVER_READ_TIMEOUT_SECONDS=5s
DETECTOR_API_SERVER_WRITE_TIMEOUT_SECONDS=10s
DETECTOR_API_SERVICE_PORT=8080
DETECTOR_API_SPEED_THRESHOLD_MILES_PER_HOUR=500
DETECTOR_API_GEO_IP_DATABASE_FILE_LOCATION=./GeoLite2-City.mmdb
DETECTOR_API_STASTD_NAMESPACE=superman-detector
DETECTOR_API_STATSD_ADDRESS=localhost:8125
DETECTOR_API_STATSD_BUFFER_POOL_SIZE=1000
DETECTOR_API_SQL_DIALECT=sqlite3
DETECTOR_API_SQLC_CONNECTION_STRING=/tmp/superman.db
```

### Running executable/binary
```bash
make superman-detector && ./superman-detector
```

### Running in docker
```bash
make build.docker && docker run -d superman-detector
```

### Running in docker-compose api w/ telemetry
```bash
docker-compose up
```

## Metrics
Application metrics are send via statsd over UDP. Docker-compose is setup to have a telegraf
sidecar listening for statsd with tag support. Outputs are available as a stream, and Prometheus
scraping /metrics.

## Tests
```bash
make test
```

## External Libraries
- [datadog statsd](github.com/DataDog/datadog-go)
- [govalidator](github.com/asaskevich/govalidator)
- [uuid](github.com/google/uuid)
- [gorm](github.com/jinzhu/gorm)
- [envconfig](github.com/kelseyhightower/envconfig)
- [echo](github.com/labstack/echo/v4)
- [geoip2-golang](github.com/oschwald/geoip2-golang)
- [testify](github.com/stretchr/testify)
- [haversine](github.com/umahmood/haversine)

## Resources/References
- Go 1.13
- Goland IDE (gofmt, govet, golangci-lint)
- Docker 19.03.5
- Docker-Compose 1.24.1
- Influxdata Telegraf 1.13.2