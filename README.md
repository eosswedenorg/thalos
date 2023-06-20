# Thalos

Thalos is a application that makes it easy for users to stream blockchain data from an Antelope SHIP node.

It handles all the technical stuff for you:

 * Decoding of antelope's binary format.
 * Websocket connection (with reconnection)
 * Decoding of action data according to contract ABI

And then sends the data over redis in plain json (or other popular formats if you want!)

## Compiling

You will need golang version `1.18` or later to compile the source.

Compile using make:

```shell
$ make
```

or using go directly if you dont have make installed.

```shell
$ go build -o build/thalos-server cmd/main/main.go
```

## Install

There are several ways to install thalos, via package manager, downloading a pre-built binary or building directly from source.

### Package Managers

* [Debian/Ubuntu based (apt)](docs/install/debian.md)

## Manually using prebuilt binaries

You can get the latest archive package [here](https://github.com/eosswedenorg/thalos/releases/latest)

Simply download using your webbrowser or via curl:

```sh
curl -Ls https://github.com/eosswedenorg/thalos/releases/download/<version>/thalos-server-<version>-linux-amd64.tar.gz | tar -z --one-top-level=thalos -xvf -
```

**NOTE**: Using curl command above, the files are extracted into the `thalos` subdirectory of the current directory where the command is run.

## From source

Follow the instructions from the [Compiling](#compiling) section.

After building the binary you can install it along with basic config file and start/stop scripts using `install.sh`

```shell
./install.sh /path/to/your/directory/of/choice
```

## Runtime dependencies

Make sure `redis` is installed as thalos uses it for both cache and message broker.

## Author

Henrik Hautakoski - [henrik@eossweden.org](mailto:henrik@eossweden.org)