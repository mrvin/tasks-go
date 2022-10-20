## Сервис хранилище изображений

Необходимо написать сервис на Golang работающий по gRPC.

Требования:
1. Принимать бинарные файлы (изображения) от клиента и сохранять их на
жесткий диск.
2. Иметь возможность просмотра списка всех загруженных файлов в формате:
Имя файла | Дата создания | Дата обновления
3. Отдавать файлы клиенту.
4. Ограничивать количество одновременных подключений с клиентами:
- на загрузку/скачивание файлов - 10 конкурентных запросов;
- на просмотр списка файлов - 100 конкурентных запросов.

#### Сборка сервера и клиента
```shell script
$ cd server/
$ make build
go build -o ../bin/server-imgstorage
$ cd ../client/
$ make build
go build -o ../bin/client-imgstorage
```

#### Запуск сервера
```shell script
$ cd bin/
$ ./server-imgstorage -config ../configs/imgstorage.yml
2022/10/20 19:41:55 Start gRPC server: localhost:55555
2022/10/20 19:50:14 Image "Cindy.jpg" upload, 78401 bytes
2022/10/20 19:51:15 Image "Claudia.jpg" upload, 112714 bytes
2022/10/20 19:54:52 Image "Claudia.jpg" download, 112714 bytes
2022/10/20 19:55:22 Image "Cindy.jpg" download, 78401 bytes

```

#### Загрузка файлов клиентом
```shell script
$ cd bin/
$ ./client-imgstorage -upload -name ../testdata/Cindy.jpg
2022/10/20 19:50:14 Upload image "Cindy.jpg"
$ ./client-imgstorage -upload -name ../testdata/Claudia.jpg
2022/10/20 19:51:15 Upload image "Claudia.jpg"
```

#### Скачивание файлов клиентом
```shell script
$ ./client-imgstorage -download -name Claudia.jpg
2022/10/20 19:54:52 Image "Claudia.jpg" saved, 112714 bytes
$ ./client-imgstorage -download -name Cindy.jpg
2022/10/20 19:55:22 Image "Cindy.jpg" saved, 78401 bytes
```

#### Получение списка файлов клиентом
```shell script
$ ./client-imgstorage -list
File name    Modified date
---------    -------------
Cindy.jpg    20 Oct 2022 16:50
Claudia.jpg  20 Oct 2022 16:51
```
