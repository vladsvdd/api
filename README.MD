1) Чтобы поднять docker контейнера, запускаем команду
```
doker compose up --build
```
2) Чтобы запустились миграции нужно запустить команды из файла entrypoint.sh
либо раскоментировать строку ниже в файле docker-compose.yml и запустить docker заного. Потом нужно закоментировать,
чтобы контейнер нормально запустился(без комментирования этой строки, контейнер останавливает работу)
```
command: ["./entrypoint.sh"]
```

Проверяем работу API
```
GET
http://172.101.0.6:8081/api/v1/good/list/?limit=5&offset=1

POST
http://172.101.0.6:8081/api/v1/create/1
добавляем в тело запроса
{
    "name": "кек",
    "description": "чебурек"
}

DELETE
http://172.101.0.6:8081/api/v1/good/remove/2?projectId=1

PATCH
http://172.101.0.6:8081/api/v1/good/reprioritize/?id=31&projectId=1
body:
{
    "newPriority": 27
}

PATCH
http://172.101.0.6:8081/api/v1/good/update?id=1&projectId=1
{
    "name": "кек",
    "description": "чебурек32"
}
```