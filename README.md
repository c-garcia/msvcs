Files for the microservices workshop hold at VTS media.

## tooling

* Go (1.9 or newer)
* `make`
* `docker` (17.12) and docker-compose 
* your favourite IDE or editor

## directory structure

* Make sure you clone this repo on your `$GOPATH/src/notif2`
* Directories
  * `notif2`: everything and `.go` files for entities
    * `test`: unit tests
    * `component`: component tests
    * `e2e`: e2e tests (with `e2e` build tag)
    * `ctrl`: Controller Implementations and UC interfaces
    * `uc`: Use Case Implementations
    * `svcs`: Services (application or domain)
    * `cmd`: executable files (as per the Go conventions)
    * `infra`: infrastructure configuration
    * `docker`: custom container definition
    * `vendor`: vendoring directory. Manage it with `dep`
    
## Infrastructure components

* SMTP server(mailhog)
  * `mail`: 1025 (SMTP)
  * `mail`: 8025 (Web admin and API). Exposed to 8025
* RabbitMQ (Multi-instance)
  * `rabbitmq`: 5672 (AMQP) via HAPROXY. Exposed to 5672
  * Admin ports (`admin`/`admin`) are exposed to random local
    ports. Check with `docker-compose ps`
* Consul
  * `consul`: 8500 (Admin interface) exposed to 8500
* User Service Double: 
  `users`:9999. Exposed to 9999
* Grafana:
  * `grafana`:3000 (Admin interface with `admin`/`admin`). Exported to 3000
* Statsd:
  * `telegraf`:8125 udp. Exposed to 8125 udp
* InfluxDB
  * `influxdb`. No exposed ports to docker host
    
## environments

### development

* designed for fast feedback while working in your IDE/editor
* run through `-devel` make targets
* all services exposed to local ports
* `notif` service is not included. You need to manually start it.

### test

* run through `-test` make targets
* designed for being able to run any test in any supporting infra, your
  workstation or a CI environment
* it runs all services, including the most recently built version of `notif` 
* ideal to run the e2e tests
* in a real environment, the artifact should be able to move
  unmodified through different environments. To maximise learning, we
  have not been super strict with this

### note: the WIP tag

* When you are working with a specific test, you can mark it with the
  `wip` build tag to limit the build and test process to this single test

## make targets

* `make help`: shows a list of available targets
* `make test`: runs unit and component tests
* `make e2e`: runs end to end tests (ensure a environment is up)
* `make *-devel`: starts, stops or gets the status of the devel environment
* `make *-all`: starts, stops or gets the status of the test environment
* `make notif`: builds the service binary
