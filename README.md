
apitinyfile
===========

apitinyfile is a tiny api server to allow (some or all of) reading, writing,
and deleting files in a single directory, with support for TLS and basic
authentication.

It is intended as a lightweight method of transferring files in ad-hoc or
narrowly-focussed situations, with less overhead and exposure than using
something broad like ssh.

It should be run as a non-privileged user with only the required permissions
for the directory in question.


Installation
------------

If you have `go` installed, you can do:

    go install github.com/gavincarr/apitinyfile@latest

which installs the latest version of `apitinyfile` in your `$GOPATH/bin`
or `$HOME/go/bin` directory (which you might need to add to your `$PATH`).


api
---

apitestfile supports the following routes (if enabled by runtime options
on the server):

```
- GET /:filename - returns the contents of `filename` in your directory (if it exists)
- PUT /:filename - writes the request body to `filename` in your directory (creates/overwrites)
- DELETE /:filename - deletes `filename` in your directory (if it exists)
```

For example, if $URL is your apitinyfile base endpoint, you can do the following
with `curl` to write, fetch, and delete a file called `foo`:

```
# Copy example.txt to `foo`
curl -X PUT --data-binary @example.txt $URL/foo

# Fetch `foo`
curl $URL/foo

# Delete `foo`
curl -X DELETE $URL/foo
```


Author
------

Copyright 2023 Gavin Carr <gavin@openfusion.net>.


Licence
--------

apitinyfile is available under the terms of the MIT Licence.

