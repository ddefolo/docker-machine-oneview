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

### Setup ICsp boot image

Provisioning an operating system onto the allocated hardware that this driver will create requires us to have a working Insight Control server provisioning (ICsp) OS build plan (OSbp) created.

1. Use one of the standard RedHat Linux 7.1 boot images.
2. The boot image can be named anything, but this driver will use RHEL71_DOCKER_1.8 for the default.  If you want an alternate name, please make sure to pass the --oneview-os-plan option with the alternate name.
3. On the build plan (step 25) add a bash script step to the next to the last step that is waiting on the server to boot.  The script contents for this step should appear as the following:
   get script from : drivers/oneview/scripts/docker_os_build_plan.sh
4. For the OS build plan, setup two custom attributes that will be popluated and passed to the script on step 3.   The first attribute is `docker_user`, and the second attribute is `public_key`.   Note, If a different value for docker_user is desired, then you should set the user and specify the user through the --oneview-ssh-user argument.  This will configure the provisioning to occur with that user.   By default the public key will be generated for each environment you create, this is not configurable.  The script added should have arguments :
@docker_user@ "@public_key@" @docker_hostname@ "@proxy_config@" "@proxy_enable@" "@interface@"

### Extra setup on OS Build Plan

There may be additional network requirements that you should specify so that specific local network requirements are met.  For example, proxy settings.
To work around this, please setup an additional build script that you customize the target os configuration with what will be required for your network.

This step is only needed for systems that need additional setup. Some examples

#### update /etc/yum.repos.d/ with special repos
```bash
cat > /etc/yum.repos.d/RedHat.repo << REPO
# base
[RedHat-7.1Server-x86_64-Server]
name=RedHat Linux 7.1Server - os - x86_64 - Server
baseurl=http://linuxcoe.corp.hp.com/LinuxCOE/RedHat-yum/7.1Server/en/os/x86_64
gpgcheck=0
# updates x86_64
[RedHat-7.1Server-x86_64-errata]
name=RedHat-7.1Server-x86_64-errata
baseurl=http://linuxcoe.corp.hp.com/LinuxCOE/RedHat-updates-yum/7.1Server/en/os/x86_64
gpgcheck=0
REPO
```

#### update /etc/yum.conf to use network appropriate settings when you are in a LR1 for example
```bash
$proxy=$1
cat > /etc/yum.conf << YUMCONF
[main]
cachedir=/var/cache/yum/$basearch/$releasever
keepcache=0
debuglevel=2
logfile=/var/log/yum.log
exactarch=1
obsoletes=1
gpgcheck=1
plugins=1
installonly_limit=3
proxy=$proxy
YUMCONF
```

## Version supported

The `--oneview-apiversion` option or ONEVIEW_APIVERSION environment variable should be
set to one of the following version integers.

| api version (see --oneview-apiversion) | HP OneView Version |   HP ICsp Version     |
|----------------------------------------|--------------------|-----------------------|
| 120                                    | 120                | 108                   |
| 200                                    | 200                | 108                   |

TODO: we need to re-work this logic to make it automatic.

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
|                            |
| `--oneview-ilo-user`       | ILO user id that is used during ICSP server creation
| `--oneview-ilo-password`   | ILO password that is used durring ICSP server creation
| `--oneview-ilo-port`       | Optional ILO port to use, defaults to 443


## OneView Server Template

* TODO describe the Server Template configuration steps and references.

## OneView ICSP OS Build Plan

OS build plans will help configure and setup docker OS with OneView so that
docker is ready to install.  This section will describe how to customize and
configure docker specific OS server build plans.

* TODO need steps / reference to OS Server build plan.
