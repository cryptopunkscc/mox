# Money over XMPP

## <span style="color:red">EXPERIMENTAL CODE - EXPECT SERIOUS ISSUES</span>

## Why?

Exchanging bitcoin addresses and lightning invoices is cumbersome and a user experience problem. It would be
much better to have a reusable, decentralized bitcoin handle that you can publicly disclose without harm
to your privacy.

## What?

Mox let's you connect your LND node to your XMPP address, so that you can send satoshis to an address that looks
like a regular email address!

## How?

### Build the code

```sh
$ git clone https://github.com/cryptopunkscc/mox
$ cd mox
$ make
```

### Config file

The minimal config file:

```yaml
xmpp:
  jid: "yourjid@somehost.example"
  password: "xmpppassword"
wallet:
  backend: "lnd"
```

You can use a non-standard LND directory:

```yaml
wallet:
  lnddir: "/home/user/.lnd-alt"
```

If you need to connect to a remote LND node:

```yaml
wallet:
  host: "lnd.myserver.example"
  port: 10009,
  cert: "/path/to/tls.cert"
  macaroon: "/path/to/admin.macaroon"
```

### Run

```sh
$ moxd -c config.json
```