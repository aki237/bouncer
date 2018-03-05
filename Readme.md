# bouncer - simple host based HTTP router

bouncer is simple host header based HTTP router.
I built this from ground up for my own use at aki237.me.

This application routes(proxies) the requests based on the `Host`
HTTP header. Say there is a git instance running in `localhost:3000`,
and there is a file server running at `localhost:8100`, bouncer is used
route the HTTP requests based on the Host header. If `git.example.com`
is requested the request is forwarded to `localhost:3000` whereas
`fs.example.com` is forwarded to `localhost:8100`. Further if TLS certificates
are specified (certificates with wildcard domain names or multiple valid domain names),
bouncer automatically listens at `:80` and `:443`. So there is no need to
configure cerificates for all the apps.
(Also `http` requests sre automatically redirected to `https`)

## Usage
`bouncer` looks for a configuration at `/etc/bouncer.conf`. Configuration can be 
explicitly specified through the commandline option `-c`.
The available options can be seen through `-help` option.

```
$ bouncer -help
Usage of bouncer:
  -c string
    	Configuration file to use (default "/etc/bouncer.conf")
  -cert string
    	Full chain certificate to use for TLS connections
  -key string
    	Private Key to use for TLS connections
```

Provide the public certificate through `-cert` option and private key
through `-key` option.

### Configuration

The format of the configuration file is very simple.
```conf
# This is a comment

#  hostname     local
git.example.com :3000
fs.example.com  :8080
```

All empty lines and commented lines are ignored.
Every other line contains a hostname and local address pair separated
by one or more whitespaces in which hostname comes first.

See [meta/example.conf](meta/example.conf)
