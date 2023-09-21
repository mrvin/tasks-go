## Сервис Фибоначчи

Реализовать сервис, возвращающий срез последовательности чисел из ряда
Фибоначчи.

Сервис должен отвечать на запросы и возвращать ответ. В ответе должны
быть перечислены все числа, последовательности Фибоначчи с порядковыми
номерами от x до y.

Требования:
1. Требуется реализовать два протокола: HTTP REST и GRPC.
2. Кэширование. Сервис не должен повторно вычислять числа
из ряда Фибоначчи. Значения необходимо сохранить в Redis или Memcache.
3. Код должен быть покрыт тестами.

#### Сборка и запуск сервиса в Docker

```shell script
git clone https://github.com/mrvin/tasks-go.git
cd tasks-go/004-fibonacci/
docker build -t fib-server .
```

```shell script
docker pull redis
docker run --name fib-redis-db -d redis
docker run --name fib-server --link fib-redis-db:redis --publish 8080:8080 --publish 55555:55555 fib-server

```

```shell script
curl "http://localhost:8080/fibonacci?from=3&to=10"
```
