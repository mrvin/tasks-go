## Сервис для сокращения URL-адресов

Необходимо разработать сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/.
Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовать
3 метода и их логику:

#### Создание нового сокращенного URL-адреса
 - Эндпоинт: POST /api/v1/data/shorten
 - Параметры запроса:
    - JSON-объект в теле запроса с параметрами:
        - url – исходный, полный URL-адрес
        - alias - сокращенный путь (необязательный параметр)
 - Статус ответа 201 если новый URL-адреса создан успешно. Ответ содержит JSON-объект.
#### Перенаправление URL-адреса
 - Эндпоинт: GET /api/v1/{alias}
 - Статус ответа 302 (Перенаправление) если alias существует
 - Статус ответа 404 если alias не найден
#### Удаление сокращенного URL-адреса
- Эндпоинт: DELETE /api/v1/{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно

```shell script
$ curl -i -X GET 'http://localhost:8081/health'
$ curl -i -X POST 'http://localhost:8081/api/v1/data/shorten' -H "Content-Type: application/json" -d '{"url":"https://en.wikipedia.org/wiki/Systems_design","alias":"zn9edcu"}'
$ curl -i -X GET 'http://localhost:8081/api/v1/zn9edcu'
$ curl -i -X DELETE 'http://localhost:8081/api/v1/zn9edcu'
```


https://github.com/xali1ove/ozon-fintech
