[![Build Status](https://travis-ci.org/stahlstift/go-dnsbutler.svg?branch=master)](https://travis-ci.org/stahlstift/go-dnsbutler) [![Go Report Card](https://goreportcard.com/badge/github.com/stahlstift/go-dnsbutler)](https://goreportcard.com/report/github.com/stahlstift/go-dnsbutler)

# go-dnsbutler

![logo](https://raw.githubusercontent.com/stahlstift/go-dnsbutler/master/_assets/butler.png)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdemaggus83%2Fgo-dnsbutler.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdemaggus83%2Fgo-dnsbutler?ref=badge_shield)

A tool to update multiple DynDNS providers at once.

## Why

Some providers doesn't allow wildcard subdomains for a DynDNS. My router isn't able to update more then one endpoint and so it's not possible to use differnt domains or sub-domains like jenkins.example.org, gitea.example.org, www.myotherdomain.com ...

There are workarounds for such a case like bash scripts but I want a small, reliable and easy solution running on one of my pis.

## Rewrite

The first version of dnsbutler was done in nodejs. I have rewritten this tool in golang to get rid of the nodejs runtime dependency on my old raspberry pi 1. The initial rewrite was a fun-two-hour-sunday-afternoon-project. I'm still impressed how productive you can be with go! But to be honestly - I also run into some problems like a deadlock ;)

## Getting started

This describes how to get up and running on a Raspberry Pi 1 with Rasbian stretch.

```bash
# Create a user for dnsbutler
user@jessica:~ $ sudo adduser dnsbutler --disabled-login --disabled-password
user@jessica:~ $ cd /home/dnsbutler
# Switch to user dnsbutler
user@jessica:~ $ sudo su dnsbutler

#arm6 for raspberry A, A+, B, B+, Zero
#arm7 for 2, 3, ...
dnsbutler@jessica:~ $ wget https://github.com/stahlstift/go-dnsbutler/releases/download/v0.1.4/dnsbutler-arm6-linux
dnsbutler@jessica:~ $ mv dnsbutler-arm6-linux dnsbutler
dnsbutler@jessica:~ $ chmod +x dnsbutler

# First start dnsbutler to test if everything is working and to generate the dnsbutler.json
dnsbutler@jessica:~ $ ./dnsbutler

# Configure now your targets
dnsbutler@jessica:~ $ nano dnsbutler.json

# "%s" will be then be replaced with the new IP address

# Example:
{
    "ipProvider": "https://api.ipify.org/",
    "listenAddr": ":5000",
    # waitInSec is optional - if not present or 0 - the updates will be fired instantly
    "waitInSec" : 5,
    #optional secret - if not present it will be ignored
    "secret": "YOUR_SECRET",
    "targets": [
        "https://dynamicdns.park-your-domain.com/update?host=@&domain=example.org&password=mysecret&ip=%s",
        "https://dnsentry.example.org:secret@example.org/update?hostname=build.example.org&myip=%s"
    ]
}

# Switch back to your normal user
dnsbutler@jessica:~ $ exit

# If you have you firewall active, and you should have, add a rule like this (my router has the ip 192.168.178.1)
user@jessica:~ $ sudo ufw allow from 192.168.178.1 to any port 5000 proto tcp

# Back to user home
user@jessica:~ $ cd ~

# Get the systemd script and setup
user@jessica:~ $ wget https://raw.githubusercontent.com/stahlstift/go-dnsbutler/master/scripts/systemd/dnsbutler.service
user@jessica:~ $ sudo mv dnsbutler.service /etc/systemd/system/
user@jessica:~ $ sudo chmod 755 /etc/systemd/system/dnsbutler.service

# systemctl enable will ensure that dnsbutler is started after a reboot
user@jessica:~ $ sudo systemctl enable dnsbutler.service
user@jessica:~ $ sudo service dnsbutler start
user@jessica:~ $ sudo service dnsbutler status
```

## Hints

### Update URLs

#### namecheap

```bash
# host=@ means the base url like example.org
# you can add an wildcard (*) subdomain - then the host would be *
# if you add a specific subdomain like wordpress.example.org - then the host would be wordpress
https://dynamicdns.park-your-domain.com/update?host=<YOUR_HOST>&domain=<YOUR_DOMAIN>&password=<YOUR_PASSWORD>&ip=%s
```

#### Strato

```bash
https://<YOUR_USERNAME>:<YOUR_PASSWORD>@<YOUR_HOST>/update?hostname=<FULL_DOMAIN>&myip=%s
```

### Routers

#### AVM FritzBox

The url to insert into the FritzBox (Internet/Freigaben/DynDNS) will look like

```bash
http://ipforyourserver:5000/?ip=<ipaddr> 

#with a secret
http://ipforyourserver:5000/?ip=<ipaddr>&secret=YOUR_SECRET
```

The domain, username and password fields from the fritzbox are ignored and can be filled with some random strings.

## Changelog

### 0.1.4 (2019-03-22)

#### Project

* readme updated
* changelog added

#### Bug fixes

* remove a deadlock causing #2 ([#2](https://github.com/stahlstift/go-dnsbutler/issues/2))

#### Features

* wait n-seconds before calling DynDNS providers ([#7](https://github.com/stahlstift/go-dnsbutler/issues/7))
* update only if a correct secret is provided in qry param ([#6](https://github.com/stahlstift/go-dnsbutler/issues/6))

### 0.1.0 (2019-03-18)

Initial release

#

Made with 🍺 and ❤️


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdemaggus83%2Fgo-dnsbutler.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdemaggus83%2Fgo-dnsbutler?ref=badge_large)