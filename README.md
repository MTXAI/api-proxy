# API Proxy



## Getting started

For example, OpenAI API.

### 1. build proxy

Run build.sh

```shell
bash build.sh
```



Go build

```shell
go build -v -o ./output/ ./cmd/...
```



### 2. Start proxy

Start proxy http server

```shell
./output/openai-proxy --addr=0.0.0.0 --port=6789 --remote-addr=api.openai.com
```



Start proxy https server

```shell
./output/openai-proxy --addr=0.0.0.0 --port=6789 --remote-addr=api.openai.com --client-cert-file=./client/ca-cert.pem --client-key-file=./client/ca-key.pem --server-cert-file=./server/ca-cert.pem --server-key-file=./server/ca-key.pem
```



Usage detail

```shell
./output/openai-proxy -h

  -addr string
        listen addr, --addr=0.0.0.0 (default "0.0.0.0")
  -client-cert-file string
        client cert file, --client-cert-file=./client/ca-cert.pem
  -client-key-file string
        client key file, --client-key-file=./client/ca-key.pem
  -debug
        debug logs level, --debug (default false)
  -logs-dir string
        output logs dir, --logs-dir=./logs (default "./logs")
  -port int
        listen port, --port=6789 (default 6789)
  -remote-addr string
        remote api addr, --remote-addr=api.openai.com (default "api.openai.com")
  -server-cert-file string
        server cert file, --server-cert-file=./server/ca-cert.pem
  -server-key-file string
        server key file, --server-key-file=./server/ca-key.pem
  -stdout
        output logs stdout, --stdout (default false)
```



### 3. Using proxy

Replace api addr with proxy addr

```shell
curl http://127.0.0.1:6789/v1/chat/completions -H "Content-Type: application/json" -H "Authorization: Bearer <your openai api key>"   -d '{"model": "gpt-3.5-turbo","messages": [{"role": "system","content": "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair."},{"role": "user","content": "Compose a poem that explains the concept of recursion in programming."}]}'
```

