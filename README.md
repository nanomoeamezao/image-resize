# Тестовое задание - сервер, изменяющий размер изображения

## Собрать контейнер
```
make build
docker build -t res .
docker run -p 3300:3300 -d --name res-con res 
```

## Ручка:
```
http://localhost:3300/resize?url=encoded_url&height=xx&width=xx
```
