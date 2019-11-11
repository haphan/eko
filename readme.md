# haphan/eko

```bash
▶ curl -Ss -X POST \
 -d "foo=bar&a=b" \
 -H "Accept: application/json" \
 -H "Foo: bar" -H "foo: two" \
 "http://localhost.127.0.0.1.nip.io:8080/nz/"  | jq .
{
  "method": "POST",
  "path": "/nz/",
  "headers": {
    "accept": "application/json",
    "content-length": "11",
    "content-type": "application/x-www-form-urlencoded",
    "foo": "bar, two",
    "user-agent": "curl/7.65.1"
  },
  "body": "foo=bar&a=b",
  "hostname": "localhost.127.0.0.1.nip.io:8080",
  "subdomains": [
    "localhost",
    "127",
    "0",
    "0",
    "1",
    "nip",
    "io"
  ],
  "query": {},
  "protocol": "http/1.1",
  "remoteaddr": "127.0.0.1:63984",
  "os": {
    "hostname": "MR-SG-L032.local"
  }
}
```