# go-flowdock #

go-flowdock is Go client library for accessing the [Flowdock API][].

[![GoDoc](https://godoc.org/github.com/jtdoepke/go-flowdock/flowdock?status.png)](https://godoc.org/github.com/jtdoepke/go-flowdock/flowdock)
[![Build Status](https://travis-ci.org/jtdoepke/go-flowdock.png?branch=master)](https://travis-ci.org/jtdoepke/go-flowdock)
[![Coverage Status](https://coveralls.io/repos/jtdoepke/go-flowdock/badge.png)](https://coveralls.io/r/jtdoepke/go-flowdock)

go-flowdock requires Go version 1.1 or greater.

## Usage ##

```go
import "github.com/jtdoepke/go-flowdock/flowdock"
```

### Authentication ###

The go-flowdock library does not directly handle authentication.  Instead, when
creating a new client, pass an `http.Client` that can handle authentication for
you.  The easiest and recommended way to do this is using the [goauth2][]
library, but you can always use any other library that provides an
`http.Client`.  If you have an OAuth2 access token (for example, a [personal
API token][]), you can use it with the goauth2 using:

```go
t := &oauth.Transport{
  Token: &oauth.Token{AccessToken: "... your access token ..."},
}

client := flowdock.NewClient(t.Client())

// list all flows the authenticated user is a member of or can join
flows, _, err := client.Flows.List(true, nil)
```

See the [goauth2 docs][] for complete instructions on using that library.

Some API methods have optional parameters that can be passed. For example,
To not return users when listing Flows you can pass in options:

```go
client := flowdock.NewClient(t.Client())
opt := flowdock.FlowsListOptions{User: false}
flows, _, err := client.Flows.List(true, &opt)
```

For complete usage of go-flowdock, see the full [package docs][].

## Contributing ##

This is very early in the implementation and I am basing the client heavily on
the [go-github][] implementation. Feel free to open a pull request and use this
lib or go-github as a guide.

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.

[Flowdock API]: https://www.flowdock.com/api
[goauth2]: https://code.google.com/p/goauth2/
[goauth2 docs]: http://godoc.org/code.google.com/p/goauth2/oauth
[personal API token]: https://flowdock.com/account/authorized_applications
[package docs]: http://godoc.org/github.com/jtdoepke/go-flowdock/flowdock
[go-github]: https://github.com/google/go-github
