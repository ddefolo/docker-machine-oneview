<!--[metadata]>
+++
title = "OneView"
description = "HP OneView driver for machine"
keywords = ["machine, OneView, driver"]
[menu.main]
parent="smn_machine_drivers"
+++
<![end-metadata]-->

# HP OneView
Create machines using HP OneView API 1.20

## Pre-Req:

* setup enclosure and server profile
* setup icsp boot image
* setup service accounts

### Setup enclosure and server profile

### Setup icsp boot image

Provisioning an operating system onto the allocated hardware that this driver will create requires us to have a working Insight Controller Server Provisioning (ICSP) build plan created.

1. Use one of the standard RedHat Linux 7.1 boot images.
2. The boot image can be named anything, but this driver will use RHEL71_DOCKER_1.8 for the default.  If you want an alternate name, please make sure to pass the --oneview-os-plan option with the alternate name.
3. On the build plan add a bash script step to the next to the last step that is waiting on the server to boot.  The script contents for this step should appear as the following:

```bash
#!/bin/bash
echo "This script will pre-configure the server to run docker"
DOCKER_USER_INPUT=$1
DOCKER_PUBKEY_INPUT=$2
if [ -z "${DOCKER_PUBKEY_INPUT}" ]; then
  echo "ERROR : this script requires a public key for docker user!"
  echo "USAGE: $0 <docker user> '<public key>'"
  exit 1
fi

DOCKER_USER=${DOCKER_USER_INPUT:-"docker"}
DOCKER_PUBKEY=${DOCKER_PUBKEY_INPUT}
# boot the external interface, replace this to another interface dependening on your hardware
ifup eno50

# optionally set some proxy server configuration options up
# export http_proxy="http://proxy:8080"
# export https_proxy="http://proxy:8080"
# export HTTP_PROXY="http://proxy:8080"
# export HTTPS_PROXY="http://proxy:8080"

# create a service account
useradd "${DOCKER_USER}" -d "/home/${DOCKER_USER}"

# setup .ssh folder
if [ ! -d "/home/${DOCKER_USER}/.ssh" ]; then
  mkdir -p "/home/${DOCKER_USER}/.ssh"
  chmod 700 "/home/${DOCKER_USER}/.ssh"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh"
fi
if [ ! -f "/home/${DOCKER_USER}/.ssh/authorized_keys" ] ; then
  touch "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chmod 600 "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh/authorized_keys"
fi
cat >> "/home/${DOCKER_USER}/.ssh/authorized_keys" << EOF
${DOCKER_PUBKEY}
EOF

# give sudoers access
cat >> "/etc/sudoers.d/90-${DOCKER_USER}" << SUDOERS_EOF
# User rules for icsp docker user
${DOCKER_USER} ALL=(ALL) NOPASSWD:ALL
SUDOERS_EOF

```

<!-- list-start: 4 -->4. For the OS build plan, setup two custom attributes that will be popluated and passed to the script on step 3.   The first attribute is `docker_user`, and the second attribute is `public_key`.   Note, If a different value for docker_user is desired, then you should set the user and specify the user through the --oneview-ssh-user argument.  This will configure the provisioning to occur with that user.   By default the public key will be generated for each environment you create, this is not configurable.  The script added should have arguments : @docker_user@ "@public_key@"


### Setup service accounts

## Options:

> **Note**: You must use a base operating system supported by Machine.

Environment variables and default values:

| CLI option                 | Description
|----------------------------|--------------------------------------------|
| `--oneview-ov-user`        | String User to OneView
| `--oneview-ov-password`    | String Password to OneView
| `--oneview-ov-domain`      | String Domain to OneView
| `--oneview-ov-endpoint`    | String url end point, base path
|                            |
| `--oneview-icsp-user`      | String User to ICSP
| `--oneview-icsp-password`  | String Password to ICSP
| `--oneview-icsp-domain`    | String Domain to ICSP
| `--oneview-icsp-endpoint`  | String url end point, base path
|                            |
| `--oneview-sslverify`      | Bool false means no https verification
| `--oneview-apiversion`     | Int version of api to use 120 is default
|                            |
| `--oneview-ssh-user`       | OneView build plan ssh user account
| `--oneview-ssh-port`       | OneView build plan ssh host port
|                            |
| `--oneview-server-template`| OneView server template to use for blade provisioning, see OneView Server Template for setup.
| `--oneview-os-plan`        | OneView ICSP OS Build plan to use for OS provisioning, see ICS OS Plan for setup.


## OneView Server Template

* TODO describe the Server Template configuration steps and references.

## OneView ICSP OS Build Plan

OS build plans will help configure and setup docker OS with OneView so that
docker is ready to install.  This section will describe how to customize and
configure docker specific OS server build plans.

* TODO need steps / reference to OS Server build plan.
