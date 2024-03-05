## Сервис "книги"

Спроектировать базу данных, в которой содержится авторы книг и сами книги. Необходимо
написать сервис который будет по автору искать книги, а по книге искать её авторов.

### Требования к сервису: 
 - Сервис должен принимать запрос по GRPC.
 - Должна быть использована база данных MySQL.
 - Код сервиса должен быть хорошо откомментирован.
 - Код должен быть покрыт unit тестами.
 - В сервисе должен лежать Dockerfile, для запуска базы данных с тестовыми данными.
 - Должна быть написана документация, как запустить сервис.
 - Плюсом будет если в документации будут указания на команды, для запуска сервиса и его окружения, через Makefile.

### Сборка и запуск приложения в Docker Compose

```shell script
$ make build
...............
$ make up
...............
```

### Пример использования сервиса "книги" при помощи клиента
```bash
$ ./books-client 
0 - Exit
1 - Save book
2 - Search by author
3 - Search by title
1
Title:The Go Programming Language
Authors:Alan A. A. Donovan, Brian W. Kernighan
0 - Exit
1 - Save book
2 - Search by author
3 - Search by title
2
Author:Brian W. Kernighan
Titles: [The Go Programming Language]
0 - Exit
1 - Save book
2 - Search by author
3 - Search by title
3
Title:The Go Programming Language
Authors: [Alan A. A. Donovan Brian W. Kernighan]
0 - Exit
1 - Save book
2 - Search by author
3 - Search by title
0
Exit
```