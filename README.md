# ACME webhook using the dns-lexicon python library

The python library [lexicon](https://pypi.org/project/dns-lexicon/) supports many different
DNS providers; by using it we can create a cert-manager DNS01 solver which works for any of those!

As of when this was last updated, the providers supported were:

The current supported providers are:

* Aliyun.com
* AuroraDNS
* AWS Route53
* Azure DNS
* Cloudflare
* ClouDNS
* CloudXNS
* ConoHa
* Constellix
* DigitalOcean
* Dinahosting
* DirectAdmin
* DNSimple v1, v2
* DnsMadeEasy
* DNSPark
* DNSPod
* Dreamhost
* Dynu
* EasyDNS
* Easyname
* EUserv
* ExoScale
* Gandi RPC (old) / LiveAPI
* Gehirn
* Glesys
* GoDaddy
* Google Cloud DNS
* Gransy (sites subreg.cz, regtons.com and regnames.eu)
* Hover
* Hurricane Electric DNS
* Hetzner
* Infoblox
* Infomaniak
* Internet.bs
* INWX
* Joker.com
* Linode
* Linode v4
* LuaDNS
* Memset
* Mythic Beasts (v2 API)
* Njalla
* Namecheap
* Namesilo
* Netcup
* NFSN (NearlyFreeSpeech)
* NS1
* OnApp
* Online
* OVH
* Plesk
* PointHQ
* PowerDNS
* Rackspace
* Rage4
* RcodeZero
* RFC2136
* Sakura Cloud by SAKURA Internet Inc.
* SafeDNS by UKFast
* SoftLayer
* Transip
* UltraDNS
* Value-Domain
* Vercel
* Vultr
* WebGo
* Yandex
* Zilore
* Zonomi

Though all of these should be supported I haven't tested all of them, just the
ones that I use. I have left logging on pretty heavily in the webhook which should
help with any troubleshooting.

Installation
------------

    helm upgrade -i -n cert-manager ./deploy/cert-manager-dns-lexicon-webhook --set groupName 'dns-lexicon.mycompany.com'


Credits
-------

Credit where it is due, this project was based on the [cert-manager-webhook-example](https://github.com/cert-manager/webhook-example)
project and borrowed a lot of ideas and a bit of code from the [dnsmadeeasy-webhook](https://github.com/k8s-at-home/dnsmadeeasy-webhook) webhook.

This was the first golang project I've made, so there are probably things that could be improved -- any assistance with maintenance would be appreciated.

Long term goal
--------------

It is my sincere hope that this project will become unnecessary and (cert-manager will add built-in support)[https://github.com/cert-manager/cert-manager/issues/4979].