# Checkout.com API Challenge

This is a simple API service that performs some Gateway actions on user Transactions. Those actions need to first be initiated by an authorization that provides a unique key.
That key will then be used for Capture, Refund & Void.

We assume that once a capture is made without a respective refund - meaning that there is a captured amount - void will not succeed.

# Docker Run

We will need docker installed in our system and after navigating to the folder containing the project, we
can fire up our service with:
```
$ docker build -t checkout-api-run .
$ docker run --rm -e ACCESS_SECRET=supersecret -p 2012:2012 checkout-api-run
```

This will fire up the server listening on port 2012. We can then access http://localhost:2012/login and using `username:password` we can get back an authentication token.
We use that token as a `Token` Header in all subsequent requests. More info can be found in [docs] once the server is up and running.

# Build & Testing

If we wish to build the app from scratch as well as testing our code, we should access our project folder using docker's default golang image. 

Otherwise we can install golang from scratch in our machine and do the same actions

```
$ export WORKDIR=<extracted path of app>
$ cd $WORKDIR
$ docker run -it --name golang-dev -p 2012:2012 --rm -v $(pwd):/go/src/github.com/nktsitas/checkout-techlab golang:latest
```

This will mount our current folder into a golang environment where we can build and run the service from scratch from inside the container:

Then from inside the docker container - or our installed go environment in our local machine:
```
$ cd /go/src/github.com/nktsitas/checkout-techlab
$ ACCESS_SECRET=supersecret go run main.go
```

To run golang tests
```
$ cd /go/src/github.com/nktsitas/checkout-techlab
$ go test -v ./...
```

[docs]: <http://localhost:2012/swagger/index.html>