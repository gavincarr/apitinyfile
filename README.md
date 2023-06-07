
apitinyfile
===========

`apitinyfile` is a tiny api server to allow (some or all of) reading, writing,
and deleting files in a single directory, with support for TLS and basic
authentication.

It is intended as a lightweight method of transferring files in ad-hoc or
narrowly-focussed situations, with less overhead and exposure than using
something broad like `ssh`.

It should be run as a non-privileged user with only the required permissions
for the directory in question. It should NOT be run as `root`.


Installation
------------

If you have `go` installed, you can do:

    go install github.com/gavincarr/apitinyfile@latest

which installs the latest version of `apitinyfile` in your `$GOPATH/bin`
or `$HOME/go/bin` directory (which you might need to add to your `$PATH`).


Usage
-----

On your server:

```
# By default, binds to *:3137, with no TLS or authentication.
# Requires that you specify what operations to support (`-r/-w/-d` for
# read/write/delete, or `-a` for all), and the directory to use for all
# files (absolute or relative).
apitinyfile -rwd /path/to/directory

# To use TLS, you must also supply the TLS certificate and key files to use:
apitinyfile -a -c /path/to/tls/cert -k /path/to/tls/key /path/to/directory

# And to use basic authentication, you must also supply a valid `htpasswd` file
# with users and encrypted passwords (note that you should ALWAYS use TLS with
# basic authentication, since otherwise your credentials will be travelling in
# the clear on every request):
apitinyfile -a -c /path/to/tls/cert -k /path/to/tls/key -p /path/to/htpasswd /path/to/directory

# To bind to a different port and/or hostname, use the `-l/--listen` option:
apitinyfile -a -l 192.168.10.1:3000 /path/to/directory

# See all command-line options:
apitinyfile -h
```


api
---

`apitestfile` supports the following routes (if enabled by the operations
options above on the server):

```
- GET /:filename - returns the contents of `filename` in your directory (if it exists)
- PUT /:filename - writes the request body to `filename` in your directory (creates/overwrites)
- DELETE /:filename - deletes `filename` in your directory (if it exists)
```

For example, if $URL is your `apitinyfile` base endpoint, you can do the following
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

`apitinyfile` is available under the terms of the MIT Licence.

