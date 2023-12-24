# Spoofed RoundTripper

A Go's http.RoundTripper implementation that provides a wrapper for [tls-client](https://github.com/bogdanfinn/tls-client/) and leverages [uTLS](https://github.com/refraction-networking/utls) to spoof TLS fingerprints (JA3, JA4, HTTP/2 Akamai, etc) of mainstream browsers for use in different HTTP client libraries (like [resty](https://github.com/go-resty/resty)) to bypass Cloudflare or other firewalls.

## Features

- Customized TLS Cipher Suites.
- Customized TLS Extensions.
- Built-in fingerprint profiles of mainstream browsers.
- Implements Go's http.RoundTripper so can be used in different 3rd-party HTTP client libraries.

## Install

```bash
go get -u github.com/juzeon/spoofed-round-tripper
```

## Usage

```go
package main

import (
	"fmt"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/go-resty/resty/v2"
	srt "github.com/juzeon/spoofed-round-tripper"
)

func main() {
	// Create a SpoofedRoundTripper that implements the http.RoundTripper interface
	tr, err := srt.NewSpoofedRoundTripper(
        // Reference for more: https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/client-options
		tlsclient.WithRandomTLSExtensionOrder(),// needed for Chrome 107+
		tlsclient.WithClientProfile(profiles.Chrome_120),
	)
	if err != nil {
		panic(err)
	}

	// Set as transport. Don't forget to set the UA!
	client := resty.New().SetTransport(tr).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Use
	resp, err := client.R().Get("https://tls.peet.ws/api/all")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Body()))
}
```

As SpoofedRoundTripper only implements the http.RoundTripper interface and is not http.Transport internally, some usages may need to be modified. For example:

```go
// Don't:
tr, err := srt.NewSpoofedRoundTripper(
    tlsclient.WithRandomTLSExtensionOrder(),
	tlsclient.WithClientProfile(profiles.Chrome_120),
)
if err != nil {
	panic(err)
}
client := resty.New().SetTransport(tr).SetProxy("socks5://127.0.0.1:7890")
// ERROR RESTY current transport is not an *http.Transport instance


// Do:
tr, err := srt.NewSpoofedRoundTripper(
    tlsclient.WithRandomTLSExtensionOrder(),
	tlsclient.WithClientProfile(profiles.Chrome_120),
	tlsclient.WithProxyUrl("socks5://127.0.0.1:7890"),
)
if err != nil {
	panic(err)
}
client := resty.New().SetTransport(tr)

```

## Acknowledgement

<https://github.com/bogdanfinn/tls-client/>

<https://github.com/refraction-networking/utls>

## Useful Resources

<https://tls.peet.ws/>

<https://engineering.salesforce.com/tls-fingerprinting-with-ja3-and-ja3s-247362855967/>

<https://blog.foxio.io/ja4-network-fingerprinting-9376fe9ca637>

<https://www.blackhat.com/docs/eu-17/materials/eu-17-Shuster-Passive-Fingerprinting-Of-HTTP2-Clients-wp.pdf>

