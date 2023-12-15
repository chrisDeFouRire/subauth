# Subauth

Subauth lets you check JWT tokens signed with RS256, in nginx sub request auth flow.

I built it as an alternative to the native JWT support in Nginx Plus and other API gateways.

## How it works

Nginx can perform a sub request to check JWT authentication. **Subauth** supports this by parsing the `Authorization: Bearer <token>` http header, decrypting it as a JWT signed with RS256, and checking the signature.

## Getting Started with Docker

A Docker image is published for this service (Linux amd64|arm64, size 15MB):

	docker pull ghcr.io/chrisdefourire/subauth:v0.0.4

You must specify a `PUBLIC_KEY` environment variable with the base64 encrypted RSA public key but without the `-----BEGIN PUBLIC KEY-----` and `-----END PUBLIC KEY-----` header and footer.

**Subauth** will listen on port 8080.

Run it with:

	docker run -d -p 127.0.0.1:8080:8080 -e PUBLIC_KEY="MIICnTCCAYUCB...I43gXA7Fg==" ghcr.io/chrisdefourire/subauth:v0.0.4

## Nginx configuration

Once subauth is up and running, you can use it from within your nginx configuration:

	server {
		.../...

		# this location defines the entry point for our sub request
		location /subauth {
			internal;
			proxy_pass http://127.0.0.1:8080; # or use whatever hostname:port points to subauth
			proxy_pass_request_body off;
			proxy_set_header Content-Length "";
			proxy_set_header X-Original-URI $request_uri;
		}

		# let's protect /private and check the JWT token
		location /private {
			auth_request /subauth; # the auth request will run first
			proxy_pass http://host:port; # then the proxy pass if no 401 occurs
		}
	}
