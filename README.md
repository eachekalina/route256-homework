# Запросы

## Создание

```shell
curl \
    -H "Authorization: Basic dXNlcjp0ZXN0cGFzc3dvcmQ=" \
    -H "Content-Type: application/json" \
    -d '{"id":1,"name":"SomePVZ","address":"Moscow,Russia","contact":"pvz@example.com"}' \
    -k \
    https://localhost:9443/pickup-point
```

## Получение списка

```shell
curl \
    -H "Authorization: Basic dXNlcjp0ZXN0cGFzc3dvcmQ=" \
    -k \
    https://localhost:9443/pickup-point
```

## Получение по идентификатору

```shell
curl
    -H "Authorization: Basic dXNlcjp0ZXN0cGFzc3dvcmQ=" \
    -k \
    https://localhost:9443/pickup-point/1
```

## Изменение

```shell
curl \
    -X PUT \
    -H "Authorization: Basic dXNlcjp0ZXN0cGFzc3dvcmQ=" \
    -H "Content-Type: application/json" \
    -d '{"id":1,"name":"SomePVZ","address":"Moscow,Russia","contact":"pvz_fixed@example.com"}' \
    -k \
    https://localhost:9443/pickup-point/1
```

## Удаление

```shell
curl \
    -X DELETE \
    -H "Authorization: Basic dXNlcjp0ZXN0cGFzc3dvcmQ=" \
    -k \
    https://localhost:9443/pickup-point/1
```