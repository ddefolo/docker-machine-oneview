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
