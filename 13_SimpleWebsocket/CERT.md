# command to generate a Cert for wss or https.
```bash
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
```
## What this does:
- x509 → Create a self-signed certificate.
- newkey rsa:4096 → Generate a new RSA key (4096 bits).
- keyout server.key → Save the private key.
- out server.crt → Save the certificate.
- days 365 → Valid for 1 year.
- nodes → No password required for the key (useful for servers).