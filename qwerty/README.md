## qwerty

Написать программу, выводящую все слова, которые можно набрать на
клавиатуре (QWERTY) двигаясь по соседним клавишам. Соседними считаются
клавиши, имеющие пересечения с вертикальными и горизонтальными линиями,
проведёнными через рассматриваемую клавишу. Например, для клавиши D это
E, R, S, F, X, C (но не W), а для U это Y, I, H, J (но не K). Слово
начинается с любой из клавиш и далее может состоять только из тех букв,
которые находятся рядом, например, "DESERT". Слова из словаря. Слова
формируются движением по клавиатуре. То есть, если начало идёт от D, то
D граничит с E и это корректный переход, E граничит с S, и так далее до
перехода от R к T, они тоже соседи, поэтому всё сходится.

#### Сборка
```shell script
$ make build
go build -o qwerty -ldflags '-w -s'
```
#### Установка
```shell script
$ go install -ldflags '-w -s' github.com/mrvin/tasks-go/qwerty@latest
```

#### Пример работы программы:
```shell script
$ ./qwerty -f /usr/share/dict/words
.....
Number of searched word: 162
Max length of the searched word: 6
Max length searched word: qwerty
Execution time: 0.386591485 s
```

