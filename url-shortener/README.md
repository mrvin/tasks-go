## Сокращатель URL-адресов

#### 

```shell script
$ curl -i -X POST 'http://localhost:8081/url' -H "Content-Type: application/json" -d '{"url":"https://www.google.com/","alias":"/gg"}'
$ curl -i -X GET 'http://localhost:8081/gg'
$ curl -i -X DELETE 'http://localhost:8081/gg'
```
