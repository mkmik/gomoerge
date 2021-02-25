# gomoerge
run go mod tidy over a conflicted go.mod


## Example

```console
$ cat go.mod
module github.com/mkmik/gomoerge

go 1.16

<<<<<<< HEAD
require github.com/google/renameio v0.1.0
=======
require github.com/google/renameio v1.0.0
>>>>>>> foobar
```

```console
$ gomoerge
```

```console
$ cat go.mod
module github.com/mkmik/gomoerge

go 1.16

require github.com/google/renameio v1.0.0
```

