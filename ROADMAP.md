# Docker Machine OneView Roadmap

Bare metal provisioning integration for docker-machine.  Works in lab environments
for HP C7000 systems.  Goals are simillar to docker/machine project and will
require https://github.com/docker/machine.

Machine currently works really well for development and test environments. The
goal is to make it work better for provisioning and managing production
environments.

This is not a simple task -- production is inherently far more complex than
development -- but there are three things which are big steps towards that goal:
**client/server architecture**, **swarm integration** and **flexible
provisioning**.

### Docker Engine / Swarm Configuration
Currently there are only a few things that can be configured in the Docker Engine and Swarm.  This will enable more operations such as Engine labels and Swarm strategies.

Project Planning
================

An [Open-Source Planning Process](https://github.com/docker/machine/wiki/Open-Source-Planning-Process) is used to define the Roadmap. [Project Pages](https://github.com/docker/machine/wiki) define the goals for each Milestone and identify current progress.
