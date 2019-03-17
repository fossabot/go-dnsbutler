# go-DNSButler

This tool will update multiple DynDns provider at once.

## Why

Some providers like Strato doesn't allow wildcard subdomains for a DynDNS and some routers just allows one endpoint to update a DynDNS service. So it's not possible to use different endpoints with a reverse proxy on subdomains like (jenkins.example.org, gitea.example.org, ...)

There are workarounds for such a case like bash scripts but I want a simple, stable and easy solution running on one of my pis.

## Rewrite

I have rewritten this tool in golang from JS to get rid of nodejs as a runtime dependency on my old raspberry pi 1. The rewrite was a fun-two-hour-sunday afternoon project. I'm still impressed how productive you can be with go!

## Getting started

This describes how to get up and running on a Raspberry Pi with Rasbian stretch.

```bash
# Create a user for dnsbutler
sudo adduser dnsbutler --disabled-login --disabled-password

cd /home/dnsbutler

# Switch to user dnsbutler
sudo su dnsbutler

wget ...

tar xvfz dnsbutler-arm5-0.1.0.tar.gz

chmod +x dnsbutler

# First start to test if everything is working
# and to generate the dnsbutler.json
./dnsbutler

# Configure your targets
nano dnsbutler.json

# "%s" will be replaced with the new IP

# Example:
{
    "ipProvider": "https://api.ipify.org/",
    "listenAddr": ":5000",
    "targets": [
        "https://dynamicdns.park-your-domain.com/update?host=@&domain=example.org&password=mysecret&ip=%s",
        "https://dynamicdns.park-your-domain.com/update?host=*&domain=example.org&password=secrect&ip=%s"
    ]
}

# Switch back to your normal user
exit

# If you have you firewall active, and you should have, 
# add a rule (my router has the ip 192.168.178.1)
sudo ufw allow from 192.168.178.1 to any port 5000 proto tcp

cd ~

wget url-toservice from github

sudo mv dnsbutler.service /etc/systemd/system/
sudo chmod 755 /etc/systemd/system/dnsbutler.service

sudo systemctl enable dnsbutler.service
sudo service dnsbutler start
sudo service dnsbutler status
```

### Configure the AVM FritzBox

The url for the FritzBox will look like

```bash
http://ipforyourserver:5000/?ip=<ipaddr>
```

The domain, username and password fields are ignored and can be filled with random strings.

### Todo

* add binaries
* add automatically releases per travis
* update readme with correct path for service and release
