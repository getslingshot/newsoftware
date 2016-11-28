## Slingshot Server

## Development Environment

1.  Pull cource code
2.  Make sure you have a fully operational Go development box either 32 or 64 bits.  Can be Windows or Linux, but preferred to run on Linux. With proper folder structure [ref](https://golang.org/doc/code.html)
3.  Make sure you have MongoDB installed or you have credentials for a remove MongoDB instance
4.  Code uses authenticated mode for MongoDB.
5.  Feel free to change logon parameters to your liking to a new file called `.env`:

```bash
  MONGO_HOST=localhost
  MONGO_NAME="Slingshot"
  MONGO_USER=""
  MONGO_PASS=""
```

6.  go run server.go
7.  Open web browser http://localhost:3000

## Import roles

```
mongoimport --db DB_NAME --collection roles --file roles.json
```
## About Middlewares used

There's a type `Adapter` which implements an `http.Handler`, so, our middlewares will implement this Adapter. Our middleware will have following signature:

```go
func myMiddleware() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do something clever here
			h.ServeHTTP(w, r)
		})
	}
}
```

Then use it on `router.go`

```go
http.Handle("/myroute", Chain(
    handler(ApprovedSoftwareHandler),
    mongo(db), // mongo session middleware
    authenticate(), // authentication middleware, user must be logged in
    myMiddleware(), // my clever stuff
))
```

## Defining a new route handler

1. Add a test (no kidding,  you should start with TDD ASAP!)
2. Add a new line to router.go

```go
http.Handle("/myroute", Chain(
    handler(ApprovedSoftwareHandler),
    mongo(db), // only needed if we need a session to our database
    authenticate(), // only needed if user must be logged in
))
```

3. Create a new file with your regular handler

```go
fun MyAwesomeHandler(h http.ResponseWriter, r *http.Request){
  	//...
}
```

## Accessing Database from a handler

```go
db := context.Get(r, "db").(*mgo.Database)
c := db.C("users")
```

There's a middleware called `mongo` that will take care of errors and or close session. If we are unable to connect to DB from the beggining, it will not work, neigther start.

## Development environment ( bonus )

Using CompileDaemon:

```bash
go get github.com/githubnemo/CompileDaemon
```

run:

```bash
CompileDaemon -build="go build -o server" -directory=. -command="./server" -include="*.tmpl"
```
