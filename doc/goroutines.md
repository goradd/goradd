# Go Routines

You must take special care if you are launching a Go routine from within 
a Goradd form or web response. If something panics within the Go routine,
the server will stop.

Normally, you do not need to worry about this, because every web request 
is handled within a go routine by the standard Go web server, and Goradd
automatically catches panics before they exit the go routine.

However, if you launch another go routine within your handler, and something
in that go routine panics, Goradd will not be able to catch it and your
web server will stop.

The solution is to catch panics within the go routine using a defer-recover
pattern and prevent them from exiting. If you need to access databases within
your go routine, you can call db.PutContext to get a database context that
you can use in your database calls.

For example:

```go
go func() {
    ctx := db.PutContext(nil) // get a database context
    // catch all panics, since if a panic reaches a go routine, it will crash the server
    defer func() {
        r := recover()
        log.Error(r) // log the error, or do some other useful thing with it
    }()
	
    person := model.LoadPerson(ctx, id)
    etc...
}

```