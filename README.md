# go-mastopush [![Build Status](https://travis-ci.org/buckket/go-mastopush.svg)](https://travis-ci.org/buckket/go-mastopush) [![Go Report Card](https://goreportcard.com/badge/github.com/buckket/go-mastopush)](https://goreportcard.com/report/github.com/buckket/go-mastopush) [![codecov](https://codecov.io/gh/buckket/go-mastopush/branch/master/graph/badge.svg)](https://codecov.io/gh/buckket/go-mastopush) [![GoDoc](https://godoc.org/github.com/buckket/go-mastopush?status.svg)](https://godoc.org/github.com/buckket/go-mastopush)

**go-mastopush** implements the decryption portion of the [Web Push standard](https://developers.google.com/web/fundamentals/push-notifications/) ([RFC8030](https://tools.ietf.org/html/rfc8030), [RFC8291](https://tools.ietf.org/html/rfc8291)),
as well as additional helper functions, which allow for easy decryption and parsing of Push Notifications sent by [Mastodon](https://github.com/tootsuite/mastodon).

Here’s the output of the included example project:
```sh
[buckket@uncloaked go-mastopush]$ ./go-mastopush 
2019/05/19 16:52:33 Added new push subscription (ID: 1, Endpoint: https://example.org/go-mastopush/)
2019/05/19 16:52:33 Mastodon ServerKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX="
2019/05/19 16:52:46 Incoming request from [::1]:39306
2019/05/19 16:52:46 JWT Header: &{Algorithm:ES256 KeyID: Type:JWT ContentType:}
2019/05/19 16:52:46 JWT Payload: &{Issuer: Subject:mailto:no-reply@example.org Audience:[https://example.org] ExpirationTime:1558363966 NotBefore:0 IssuedAt:0 JWTID:}
2019/05/19 16:52:46 New push notification: 
{
	"access_token": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"preferred_locale": "en",
	"notification_id": "701",
	"notification_type": "mention",
	"icon": "https://social.buckket.org/avatars/original/missing.png",
	"title": "You were mentioned by mirror",
	"body": "@buckket Hello. Testing 123"
}
```

## Installation

### From source

    go get -u github.com/buckket/go-mastopush

## Usage

[Here’s](https://github.com/buckket/go-mastopush/tree/master/example) a simple example project. Check [GoDoc](https://godoc.org/github.com/buckket/go-mastopush) for the full documentation.

## Limitations

- Only supports `aesgcm` and not `aes128gcm`. Which is fine, because Mastodon only uses the former.
  But implementing the later should be straight forward as well, as parsing the HTTP request headers is no
  longer necessary.
- Documentation is still lacking.

## Notes

- A remotely similar project which forwards notifications (to APN in this case) instead of decrypting them can be found here:
[https://github.com/DagAgren/toot-relay](https://github.com/DagAgren/toot-relay)

## License

 GNU GPLv3+
 