# HPE OneView machine tools

[![Build Status](https://travis-ci.org/HewlettPackard/docker-machine-oneview.svg?branch=master)](https://travis-ci.org/HewlettPackard/docker-machine-oneview)

This container provides a way to execute docker-machine and the HPE OneView plugin
without having to install them locally or on a given system.

It provides everything you need to get started using the docker-machine and oneview plugin.

Usefull for hacking on `docker-machine create -d oneview` commands.

## Getting started

1. Clone this repository.  Install docker and GNU make tool.
2. Use this project in a bash shell
3. run `make build` to create the container.
4. run `eval $(make alias)` to setup the local alias in your bash shell for this project.

## Using it for the first time

The container and alias that is created by this project basically wraps all
the requirements you need to run docker-machine commands.   Common commands might look
like:

```
oneview-machine --help
```

### Creating a new machine for HPE OneView

1. Create a creds.env in your local shell
2. Source the creds.env
3. Run `oneview-machine create -d oneview mymachine` and this will use the creds.env 
   vars to connect and build the machine called mymachine using HPE OneView

## Contributing

Want to hack on oneview-machine-tools? Fork this repo and start hacking.

