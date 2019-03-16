# Cab Data

## Introduction

**cabdata** implements an API for querying taxi trip counts by medallion and date. Results are cached to improve performance; requests can bypass the cache and also the cache can be flushed.
 
## Installation

- clone this repository to a location outside `$GOPATH` (eg `~/gi/cabdata)`. Or enable go modules under `$GOPATH` with `export GO111MODULE=on`

- [download and install mysql server, shell and workbench](https://dev.mysql.com/downloads/mysql/)

### import data

- import separately provided data. For example on OSX:

```
/usr/local/mysql/bin/mysql -p -u root -e 'create database ny_cab_data'
/usr/local/mysql/bin/mysql -p -u root ny_cab_data < ny_cab_data_cab_trip_data_full.sql
/usr/local/mysql/bin/mysql -p -u root ny_cab_data
```

- check data imported ok, by querying using cli or console:

`select count(*) from cab_trip_data;`

should return '73488'

`select hack_license from cab_trip_data where medallion = 'D7D598CD99978BD012A87A76A7C891B7';`

 should return:
 
 ```
| 82F90D5EFE52FDFD2FDEC3EAD6D5771D |
| 8E57A362C89B897A6913609FA8B416C4 |
| 8E57A362C89B897A6913609FA8B416C4 |
```

- create user 'api' with password 'secret' that has select permissions on the table `cab_trip_data`  

### go modules

- use go modules to install required modules -- `go build`

## Usage

- run the tests:

`go test`

- build the server and run it:

`go run main.go db.go service.go &`

- the exercise also asks for a client to be written, I believe this an excessive request for an exercise, the service can by queried using curl:

- flush the cache:

```
curl -XPOST http://localhost:12345/flush
"{result:ok}"
```

- query based on medallions and date:

```
curl -XPOST -d'{"medallions":["CFC043F3E41A505744D0FF5E63D007DD","2B1A06E9228B7278227621EF1B879A1D"],"date":"2013-12-01"}' http://localhost:12345/count
[{"medallion":"CFC043F3E41A505744D0FF5E63D007DD","count":2},{"medallion":"2B1A06E9228B7278227621EF1B879A1D","count":4}]
```

- query based on medallions and date, also requesting a fresh copy:

```
curl -XPOST -d'{"medallions":["CFC043F3E41A505744D0FF5E63D007DD","2B1A06E9228B7278227621EF1B879A1D"],"date":"2013-12-01","fresh":true}' http://localhost:12345/count
[{"medallion":"CFC043F3E41A505744D0FF5E63D007DD","count":2},{"medallion":"2B1A06E9228B7278227621EF1B879A1D","count":4}]
```

- the service is simple and will only respond with `null` for invalid requests

```
curl -XPOST -d'{"medallions":["CFC043F3E41A505744D0FF5E63D007DD","2B1A06E9228B7278227621EF1B879A1D"]}' http://localhost:12345/count
null
```

- unhandled endpoints will receive a 404:

```
curl -XPOST -d'{"medallions":["CFC043F3E41A505744D0FF5E63D007DD","2B1A06E9228B7278227621EF1B879A1D"]}' http://localhost:12345/xxx
404 page not found
```

### Libraries

- for sql queries I usually use `github.com/jmoiron/sqlx`, as it provides marshalling of query results into typed structs, rather than `interface{}`

- the service could've be written using `net/http` and json marshalling, I used `github.com/go-kit` as I wanted to familiarise myself with a new library


