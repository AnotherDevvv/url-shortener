
# URL shortener

Quick golang warmup task which allows to shorten long and sometimes ugly urls. For example

* https://translate.google.com/?sl=auto&tl=ru&text=lets%20translate%20something%20as%20example%20to%20produce%20long%20url&op=translate -> https://hostname/MJP5X7LU

Passed URLs are hashed with md5 and encoded with base32 which allows to produce $$2^{40}$$ different short links. In case collision happens ```hash(url1) == hash(url2)``` then new key with double hash is used ```hash(hash(url1))```


### User guide

```sh
curl -v -X POST "http://localhost:1323/shorten?url=https://translate.google.com/?sl=auto%26tl=ru%26text=lets%20translate%20something%20as%20example%20to%20produce%20long%20url%26op=translate"
```

```sh
curl -v http://localhost:1323/MJP5X7LU
```