# Thalos

Thalos is a application that makes it easy for users to stream blockchain data from an Antelope SHIP node.

Consult the [documentation](https://thalos.waxsweden.org/docs) for more information.

## Compiling

You will need golang version `1.20` or later to compile the source.

Compile using make:

```shell
$ make
```

or using go directly if you dont have make installed.

```shell
$ go build -o build/thalos-server cmd/thalos/*.go
```
## Author

Henrik Hautakoski - [henrik@eossweden.org](mailto:henrik@eossweden.org)