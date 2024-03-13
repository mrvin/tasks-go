## Сервис "фотогалерея"

Реализовать HTTP API фотогалереи

### Возможности API: 
 - Загрузка фото.
 	- Сохранять данные о новых файлах в БД
	- Генерировать preview
 - Просмотр списка фото
 - Удаление фото

### Дополнительные требования:
 - Формат ответа - json
 - В качестве БД использовать sqlite

### Сборка и запуск приложения
```shell script
$ cd cmd/photo-gallery-server/
$ make build
go build -o ../../bin/photo-gallery-server -ldflags '-w -s'
$ cd ../../bin
$ ./photo-gallery-server -config ../configs/photo-gallery.yml
```

### Клиент для http API
```shell script
$ ./photo-gallery-client -name ../testdata/Poster_An_Interview_with_God_2018.jpg
$ ./photo-gallery-client -name ../testdata/Poster_The_Red_Turtle_2016.jpg
$ curl -i -X GET 'http://localhost:8088/api/v1/listphotos'
$ curl -i -X DELETE 'http://localhost:8088/api/v1/photo?name=Poster_An_Interview_with_God_2018.jpg'
$ curl -i -X DELETE 'http://localhost:8088/api/v1/photo?name=Poster_The_Red_Turtle_2016.jpg'
```