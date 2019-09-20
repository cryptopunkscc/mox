# Money over XMPP

## <span style="color:red">EXPERIMENTAL CODE - EXPECT SERIOUS ISSUES</span>

## What?

Mox let's you connect your LND node to your XMPP account, so that you can send satoshis to an XMPP address.

## Why?

It's convenient. You can exchange money using a very familiar email-like address. You can have a setup where you can use a single address for email, chat and money!

It's more private. You can post your XMPP address to receive donations without publicly disclosing any bitcoin addresses.

## How?

### Build the code

```sh
$ git clone https://github.com/cryptopunkscc/mox
$ cd mox
$ go build ./cmd/moxd
$ go build ./cmd/mox-cli
```

### Config file

```json
{
  "xmpp": {
    "jid": "<jid>",
    "password": "<password>"
  },
  "chatbot": {
    "admin_jid": "<jid>"
  },
  "lnd": {
    "host": "localhost",
    "port": 10009,
    "macaroon": "/path/to/admin.macaroon",
    "cert": "/path/to/tls.cert"
  }
}

```

### Run

```sh
$ moxd -c config.json
```