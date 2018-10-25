# dns-router

simple personal dns router

# features

- default handler
- configurable dns backends per pattern
- pattern based dispatch 
- lua scripting support
- redis support
- static host override 


# Per Application DNS override

Using `LD_PRELOAD` it is possible to load a shared library and change the DNS resolution functions for a single application.

While not part of dns-router it may be useful so I figured i'd include it here. 

See https://github.com/BackSlasher/dns-override for details


# install (from source)

<pre>
   export GOPATH=$HOME/go-apps/dns-router
   go get github.com/sigmonsays/dns-router/...
   go install github.com/sigmonsays/dns-router/...
</pre>

The only remaining step is to add $GOPATH/bin to $PATH

<pre>
   export PATH="$GOPATH/bin:$PATH"
</pre>

# configuration

by default dns-router reads its configuration from /etc/dns-router/config.yaml


`Example configuration`

<pre>
   bindaddr: 127.0.0.1:53
   default:
     servers:
     - 10.11.97.15:53
     - 10.12.64.15:53
     - 10.12.65.15:53
   hosts:
     foobar.example.net.: 192.168.6.4
   backends:
   - pattern: example.net
     servers:
      - 4.4.4.4:53
      - 8.8.8.8:53
   - pattern: localdomain
     servers:
     - 127.0.0.1:53
</pre>

`Notes on above configuration`

- `hosts` allows assigning static addresses for a given host (`A` record with a single IP)

