### Get Sarama

```bash
$ go get github.com/IBM/sarama
```

### Build Docker

```bash
$ docker compose up -d
```

### Start Consumer

```bash
$ go run consumer/main.go -brokers="127.0.0.1:29091" -topics="sarama" -group="example"
```

### Start Producer

```bash
$ go run producer/main.go -brokers="127.0.0.1:29091" -topic "sarama" -producers 2 -records-number 5
```

### Example Context

```
Topic `sarama` has 2 partitions produce 4 message with key 1,2,3,4
Start 2 consumer belong to same group `example`. 
Result:
- Client 1 consume message with key 1,3
- Client 2 consume message with key 2,4
...
```