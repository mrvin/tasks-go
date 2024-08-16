## hh-client-go

Автоматическое обновление резюме на сайте https://hh.ru/ по средствам HeadHunter API: https://dev.hh.ru/.
Обновление происходит с интервалом в 4 часа. В 06:00, 10:00, 14:00, 18:00, 22:00.

#### Сборка
```shell script
$ cd cmd/hh-client-go/
$ make build
go build -ldflags '-w -s' -o bin/hh-client-go
```

### Ссылки:
- [API hh.ru. Быстрый старт](https://habr.com/ru/companies/hh/articles/303168/)
