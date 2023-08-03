# ActiveMQ
This is an ActiveMq indicator mobile testing project


## 启动ActiveMQ
```shell
cd deployment/docker-compose
docker-compose up -d
```

## 启动queue消费

```shell
go run cmd/entry.go
```

## 启动queue topic 生产

```shell
cd tests

go test -run TestPublishMessage
```

## 启动tipoc消费
```shell
cd tests
go test -run TestActiveMQ
``` 