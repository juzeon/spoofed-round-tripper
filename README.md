# Spoof RoundTripper

A spoofed Go's http.RoundTripper that leverages uTLS to support JA3, JA4, and HTTP/2 Akamai fingerprints of mainstream browsers for use in different HTTP client libraries to bypass Cloudflare or other firewalls.