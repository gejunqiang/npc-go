# npc-go
npc OpenApi command line utility

# Build

```bash
CGO_ENABLED=0 GOOS=linux go build -o bin/npc -mod=vendor -ldflags '-s -w' ./cmd/main.go
```

# Usage

## Usage with env variable

```bash
NPC_API_KEY=<key> NPC_API_SECRET=<secret> npc GET '/keypair?Action=ListKeyPair&Version=2018-02-08&Limit=9999&Offset=0'

NPC_API_KEY=<key> NPC_API_SECRET=<secret> npc POST '/keypair?Action=UploadKeyPair&Version=2018-02-08&Limit=9999&Offset=0' '{"KeyContent":"ssh-rsa xxxxxxxxxx", "KeyName": "test"}'
```

## Usage with config file

```bash
$ cat ~/.npc/api.key 
{
  "api_endpoint": "open.c.163.com",
  "api_key": "d7cde2ba0cbf2xxxxxxxxxxxxxxxxxxxx",
  "api_secret": "e36ab5cxxxxxxxxxxxxxxxxxxxxxxx",
  "region": "cn-east-1"
}
```

```bash
npc GET '/keypair?Action=ListKeyPair&Version=2018-02-08&Limit=9999&Offset=0'
```

# Usage with docker

```bash
docker run -i --rm \
    -e NPC_API_KEY=d7cde2ba0cbf2xxxxxxxxxxxxxxxxxxxx \
    -e NPC_API_SECRET=e36ab5cxxxxxxxxxxxxxxxxxxxxxxx \
    gejunqiang/npc-go \
    npc GET '/keypair?Action=ListKeyPair&Version=2018-02-08&Limit=9999&Offset=0'
```
