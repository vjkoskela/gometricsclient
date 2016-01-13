go-metrics-client
=================

<a href="https://raw.githubusercontent.com/vjkoskela/gometricsclient/master/LICENSE">
    <img src="https://img.shields.io/hexpm/l/plug.svg"
         alt="License: Apache 2">
</a>
<a href="https://travis-ci.org/vjkoskela/gometricsclient/">
    <img src="https://travis-ci.org/vjkoskela/gometricsclient.png"
         alt="Travis Build">
</a>

Implementation of [ArpNetworking's Metrics Client Java](https://github.com/ArpNetworking/metrics-client-java) for [Go](https://golang.org).

Dependency
----------

First, retrieve the library into your workspace:

    go> go get github.com/vjkoskela/gometricsclient

To use the library in your project(s) simply import it:

```go
import "github.com/vjkoskela/gometricsclient"
```

''TODO''

Development
-----------

To build the library locally you must satisfy these prerequisites:
* [Go](https://golang.org/)

Next, fork the repository, get and build:

Getting and Building:

```bash
go> go get github.com/$USER/gometricsclient
go> go install github.com/$USER/gometricsclient
```

Testing:

```bash
go> go test -coverprofile=coverage.out github.com/$USER/gometricsclient
go> go tool cover -html=coverage.out
```

To use the local forked version in your project simply import it:

```go
import "github.com/$USER/gometricsclient"
```

_Note:_ The above assumes $USER is the name of your Github organization containing the fork.

License
-------

Published under Apache Software License 2.0, see LICENSE


&copy; Ville Koskela, 2016
