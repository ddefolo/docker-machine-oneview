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
![](/docs/img/DockerMachineOneView.png)

## Pre-Req:

* setup enclosure and server profile
* setup icsp boot image
* setup service accounts

### Setup enclosure and server profile

### Setup ICsp boot image

Provisioning an operating system onto the allocated hardware that this driver will create requires us to have a working Insight Control server provisioning (ICsp) OS build plan (OSbp) created.  In addition the ICsp server should have dhcpv4 setup so that public interfaces for the server receive a routable ip address at startup.

1. Use one of the standard RedHat Linux 7.1 boot images located under OS Build Plans (ProLiant OS - RHEL 7.1 x64 Scripted Install).
2. Choose the action to save a new OS build plan.  The boot image can be named anything, but this driver will use RHEL71_DOCKER_1.8 for the default.  If you want an alternate name, please make sure to pass the --oneview-os-plan option with the alternate name.
3. On the build plan (step 25) add a bash script step to the next to the last step that is waiting on the server to boot.  The script contents for this step should appear as the following:
   get script from : ```drivers/oneview/scripts/docker_os_build_plan.sh```
You can choose to name the build step docker_os_build_prereq or anything that applies for your setup.  The purpose for this script is to prepare the environment with basic user configuration and networking startup.  The script should avoid fully provisioning docker, as this is managed by upstream docker contributions to the docker-machine project.
4. Configure the parameters for the build step that was added in step 3 to have the following arguments :
@docker_user@ "@public_key@" @docker_hostname@ "@proxy_config@" "@proxy_enable@" "@interface@"

### Build Step Arguments
Build step arguments can be controlled by options passed to the docker-machine-oneview driver.  Update these options as needed.

* @docker_user@ - used as the admin account with sudo privilidges to install and run docker commands.
* @public_key@ - durring docker-machine create a private public key stored in ~/.docker/machine folder will be generated.  This will be the public key configured for @docker_user@
* @proxy_config@ - these are host machine proxy configuration settings that will be set on the host machine.   Example:
```
export proxy_config='http_proxy=https://proxy.company.com:8080/
https_proxy=https://proxy.company.com:8080/
no_proxy=/var/run/docker.sock,company.com,localhost,127.0.0.1'
```
The string will be stored in /etc/environment for the host machine.
* @proxy_enable@ when set to true, @proxy_config@ will be saved.
* @interface@ the name of the network interface that should be used for docker and machine traffic.

### Extra setup on OS Build Plan

There may be additional network requirements that you should specify so that specific local network requirements are met.  For example, internal yum repository settings.
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

This driver will work with specific combinations of HP ICsp and HP OneView.  You can check the version by navigating to http(s)://host/rest/version endpoint.

| Supported | HP OneView Version |   HP ICsp Version     |
|----------------------------------------|--------------------|-----------------------|
| Yes                                    | 120                | 108                   |
| Yes                                    | 200                | 108                   |

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

* HP OneView 1.2 users.  Server templates are identified as server profiles that have no hardware assignment.  All settings on the server template will be used.
* HP OneView 2.0+, use server templates under HP OneView Server Templates navigation.

## OneView ICsp OS Build Plan

* HP ICsp should be configured for OS provisioning with RedHat 7.1.
* HP ICsp should have DHCP enabled for ip assigments on public and private interfaces.
