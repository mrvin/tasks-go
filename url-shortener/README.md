## Сервис для сокращения URL-адресов

Необходимо разработать сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/.
Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовывать
5 методов и их логику:

#### Регистрация пользователя
- Эндпоинт - POST /users
- Параметры запроса:
   - JSON-объект в теле запроса с параметрами:
        - user_name – имя пользователя
        - password – пароль
 - Статус ответа 201 если пользователь создан успешно
#### Создание нового сокращенного URL-адреса
 - Эндпоинт: POST /data/shorten
 - Параметры запроса:
    - JSON-объект в теле запроса с параметрами:
        - url – исходный, полный URL-адрес
        - alias - сокращенный путь (необязательный параметр)
 - Статус ответа 201 если новый URL-адреса создан успешно. Ответ должен содержать в теле JSON-объект:
    - alias – сокращенный путь
#### Перенаправление URL-адреса
 - Эндпоинт: GET /{alias}
 - Статус ответа 302 (Перенаправление) если alias существует
 - Статус ответа 404 если alias не найден
#### Удаление сокращенного URL-адреса
- Эндпоинт: DELETE /{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно
#### Получение количества переходов по сокращенному URL-адресу
 - Эндпоинт: GET /statistics/{alias}
 - Статус ответа 200 если количества переходов получено успешно. Ответ должен содержать в теле JSON-объект. Объект содержит параметры:
    - count – количества переходов по сокращенному URL-адресу

### Сборка и запуск приложения в Docker Compose
```shell script
$ make build
...............
$ make up
...............
```

### Пример использования http API
```shell script
$ curl -i -X GET 'http://localhost:8081/health'
$ curl -i -X POST 'http://localhost:8081/users' -H "Content-Type: application/json" -d '{"user_name":"Bob","password":"qwerty"}'
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8081/data/shorten' -H "Content-Type: application/json" -d '{"url":"https://en.wikipedia.org/wiki/Systems_design","alias":"zn9edcu"}'
$ curl -i -X GET 'http://localhost:8081/zn9edcu'
$ curl --user Bob:qwerty -i -X GET 'http://localhost:8081/statistics/zn9edcu'
$ curl --user Bob:qwerty -i -X DELETE 'http://localhost:8081/zn9edcu'
```

### Полезные ссылки
- [Пишем REST API сервис на Go - УЛЬТИМАТИВНЫЙ гайд](https://www.youtube.com/watch?v=rCJvW2xgnk0)
- [LRU cache](https://github.com/hashicorp/golang-lru)