# ACME webhook using the dns-lexicon python library

The python library [lexicon](https://pypi.org/project/dns-lexicon/) supports many different
DNS providers; by using it we can create a cert-manager DNS01 solver which works for any of those!

As of when this was last updated, the providers supported were:

The current supported providers are:

- Aliyun.com
- AuroraDNS
- AWS Route53
- Azure DNS
- Cloudflare
- ClouDNS
- CloudXNS
- ConoHa
- Constellix
- DigitalOcean
- Dinahosting
- DirectAdmin
- DNSimple v1, v2
- DnsMadeEasy
- DNSPark
- DNSPod
- Dreamhost
- Dynu
- EasyDNS
- Easyname
- EUserv
- ExoScale
- Gandi RPC (old) / LiveAPI
- Gehirn
- Glesys
- GoDaddy
- Google Cloud DNS
- Gransy (sites subreg.cz, regtons.com and regnames.eu)
- Hover
- Hurricane Electric DNS
- Hetzner
- Infoblox
- Infomaniak
- Internet.bs
- INWX
- Joker.com
- Linode
- Linode v4
- LuaDNS
- Memset
- Mythic Beasts (v2 API)
- Njalla
- Namecheap
- Namesilo
- Netcup
- NFSN (NearlyFreeSpeech)
- NS1
- OnApp
- Online
- OVH
- Plesk
- PointHQ
- PowerDNS
- Rackspace
- Rage4
- RcodeZero
- RFC2136
- Sakura Cloud by SAKURA Internet Inc.
- SafeDNS by UKFast
- SoftLayer
- Transip
- UltraDNS
- Value-Domain
- Vercel
- Vultr
- WebGo
- Yandex
- Zilore
- Zonomi

Though all of these should be supported I haven't tested all of them, just the
ones that I use. I have left logging on pretty heavily in the webhook which should
help with any troubleshooting.

## Installation

### local install from repo

```bash
helm -n cert-manager upgrade -i dns-lexicon-webhook ./deploy/cert-manager-webhook-dns-lexicon --set groupName='dns-lexicon.mycompany.com'
```

### Using public helm chart

```bash
helm repo add cert-manager-webhook-dns-lexicon <https://lanefu.github.io/cert-manager-webhook-dns-lexicon/>
# Replace the groupName value with your desired domain
helm install --namespace cert-manager dns-lexicon-webhook cert-manager-webhook-dns-lexicon/cert-manager-webhook-dns-lexicon --set groupName=acme.bunny.net
```

And then create a ClusterIssuer, something like this:

```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: namecheap-api-key
      namespace: cert-manager
    type: Opaque
    stringData:
      key: myusername
      secret: myapikey
    ---
    apiVersion: cert-manager.io/v1
    kind: ClusterIssuer
    metadata:
      name: namecheap-lexicon
    spec:
      acme:
        email: me@company.tld
        privateKeySecretRef:
          name: mySecretKeySecret
        server: https://acme-v02.api.letsencrypt.org/directory
        solvers:
        - dns01:
            cnameStrategy: Follow
            webhook:
              config:
                apiKeyRef:
                  name: namecheap-api-key
                  key: key
                apiSecretRef:
                  name: namecheap-api-key
                  key: secret
                production: true
                provider: namecheap
                usePassword: false
                ttl: 600
              groupName: dns-lexicon.company.com
              solverName: lexicon
```

You should be able to create additional ones for each DNS provider you need using this basic template. Note that some providers
use the `--auth-password` parameter instead of `--auth-token`; in that case you need to set `usePassword: true` in the webhook
configuration to make it work. The only way I know to check that easily is to run `lexicon <provider> --help` and check the
available arguments, but this project does not do that for you at this time.

## Original Author's Note on Credits

Credit where it is due, this project was based on the [cert-manager-webhook-example](https://github.com/cert-manager/webhook-example)
project and borrowed a lot of ideas and a bit of code from the [dnsmadeeasy-webhook](https://github.com/k8s-at-home/dnsmadeeasy-webhook) webhook.

## Fork Info

This is a fork of a fork.. I basically needed multi-arch builds..

I forked from <https://github.com/SomeBlackMagic/cert-manager-webhook-dns-lexicon>

Who Forked form <https://github.com/gradecam/cert-manager-dns-lexicon-webhook>
