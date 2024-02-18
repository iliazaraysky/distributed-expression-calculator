# Распределенный вычислитель арифметических выражений

## Оглавление

1. [Описание проекта](#описание-проекта)

2. [Установка и настройка](#установка-и-настройка)

 	2.1. [Установка Docker](#установка-docker)
 
 	2.2. [Запуск проекта](#запуск-проекта)
	
 	2.3. [Возможные ошибки](#возможные-ошибки)
 
	2.4. [Выключение проекта](#выключение-проекта)

3. [Пример использования](#пример-использования)
	
 	3.1. [Запросы в curl](#запросы-в-curl)

   	3.2. [Работа в GUI](#работа-в-GUI)

4. [Схема работы Worker](#схема-работы-worker)

5. [Схема работы Backend](#схема-работы-backend)

## 1. Описание проекта
<a name="описание-проекта"></a>
Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качестве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталось ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арифметического выражения можно вычислять параллельно.

### Используемые технологии
1. Backend и Workers написаны на Golang
2. Frontend написан на React + Bootstrap
3. База данных для сохранения состояния Postgresql
4. Брокер сообщений, для организации очереди RabbitMQ
5. Считает арифметические выражения библиотека github.com/maja42/goval

## 2. Установка и настройка
<a name="установка-и-настройка"></a>
Клонируем себе на компьютер репозиторий. Переходим в папку с проектом на компьютере

### 2.1. Установка Docker
<a name="установка-docker"></a>
Для работы с сервисом необходимо установить на компьютер Docker

Версии для Mac, Windows, Linux доступны по следующему адресу:
https://docs.docker.com/get-docker/

1. Скачиваем установочный файл под вашу операционную систему
2. Запускаем / Устанавливаем / Перезагружаем компьютер

## 2.2. Запуск проекта
<a name="запуск-проекта"></a>
Первый раз проект запускается долго. Нужно дождаться, чтобы скачались необходимые образы (я постарался прописать везде alpine, образы с минимальным набором компонентов, но получилось все равно много). Если вы запускаете проект без флага `-d` , сигналом к началу проверки приложения, могут послужить сообщения в терминале, что воркеры ждут заданий из очереди:

```
worker1   | 2024/02/18 10:08:35 No messages in the queue, waiting
worker2   | 2024/02/18 10:08:37 No messages in the queue, waiting
worker1   | 2024/02/18 10:08:40 No messages in the queue, waiting
worker2   | 2024/02/18 10:08:42 No messages in the queue, waiting
worker1   | 2024/02/18 10:08:45 No messages in the queue, waiting
```

1. Клонируем репозиторий
2. Переходим в папку с программой в командной строке
3. Для создания образов и запуска контейнеров набираем:

```bash
docker-compose up --build
```

4. Если хотите продолжить работать в этом окне можно запустить docker в фоновом режиме, так даже удобнее, не нужно будет создавать дополнительного окна для завершения работы проекта. О состоянии сборки в таком случае можно посмотреть непосредственно в приложении Docker - вкладка Containers

```bash
docker-compose up -d --build
```
## 2.3. Возможные ошибки
<a name="возможные-ошибки"></a>
Я попробовал запускать проект на разных устроуствах, на 1 из 3, вылезла ошибка:
```
ERROR: for postgres  Cannot create container for service postgres: Conflict.
The container name "/postgres" is already in use by container "ac58e56714bbf18014776cfa3e47ae31d597db721cf00d0c0b36c0621dffd10a".
You have to remove (or rename) that container to be able to reuse that name.
```

Из текста сообщения понятно, в системе уже есть контейнер с названием `postgres`. Я исправил это тем, что удалил контейнер с этим названием командой:
```bash
docker rm postgres
```

## 2.4. Выключение проекта
<a name="выключение-проекта"></a>
Набрать из консоли, находясь в той же директории откуда запускался проект

```bash
docker-compose down
```

## 3. Пример использования
<a name="пример-использования"></a>
В проекте реализован GUI. Удобнее всего посмотреть работу из него. Но первым пунктом описания сделаны примеры запросов, через `curl`. 

### 3.1. Запросы в curl. Добавление вычисления арифметического выражения.
<a name="запросы-в-curl"></a>

Отправляем POST-запрос по адресу:
```bash
http://localhost:8080/add-expression
```

Тело запроса в формате JSON, структура:
```json
{"text": "2+2+2"}
```

Пример запроса:
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"text\": \"3 + 1\"}" http://127.0.0.1:8080/add-expression
```
	
### Получение списка выражений со статусами

Отправляем GET-запрос по адресу:
```bash
http://localhost:8080/get-operations
```

Ответ получаемый от сервера в формате JSON, структура:
```json
{
    "data": [
        {
            "unique_id": "abb1f26f-2014-43f6-9f88-22fd6bf7aedc",
            "query_text": "2 + 1 + 3 + 1",
            "creation_time": {
                "Time": "2024-02-11T21:17:24.764203Z",
                "Valid": true
            },
            "completion_time": {
                "Time": "2024-02-11T21:17:36.870463Z",
                "Valid": true
            },
            "execution_time": "00:00:12.10626",
            "server_name": {
                "String": "worker3",
                "Valid": true
            },
            "result": {
                "String": "7",
                "Valid": true
            },
            "status": "Done"
        }
    ],
    "total_items": 4,
    "total_pages": 1,
    "current_page": 1,
    "items_per_page": 5
}
```

Ответ содержит информацию о пагинации, которую мы высчитываем на стороне backend и возвращаем для frontend.

Данные берутся из БД Postgresql - таблица requests 

Пример запроса:
```bash
curl -X GET -H "Content-Type: application/json" http://127.0.0.1:8080/get-operations
```

Так как у этого запроса есть пагинация, следует указать как получать следующие страницы:
```bash
http://localhost:8080/get-operations?page=2
```

### Получение значения выражения по его идентификатору.

Информацию о каждой операции также можно получить по GET-запросу по адресу:
```bash
http://localhost:8080/get-request-by-id/{uuid}
```

Ответ получаемый от сервера в формате JSON, структура:
```json
{
    "unique_id":"abb1f26f-2014-43f6-9f88-22fd6bf7aedc",
    "query_text":"2 + 1 + 3 + 1",
    "server_name":"worker3",
    "result":"7",
    "status":"Done"
}
```

Пример запроса:
```bash
curl -X GET -H "Content-Type: application/json" http://localhost:8080/get-request-by-id/abb1f26f-2014-43f6-9f88-22fd6bf7aedc
```

### Получение списка доступных операций со временем их выполнения.

У нас в задании говорится о разделении времени выполнения на операции (это я понял, когда уже было поздно), поэтому тут только запрос на установление таймаута для каждого воркера. В принципе, функционал о том же самом, и можно сделать отдельную таблицу в БД, написать отдельный запрос, и все это будет работать по ТЗ. Но в виде самого наипростейшего решения, сойдет и мой вариант... Ведь он работает =)

Данный эндпойнт умеет обрабатывать как GET так и POST запросы.
Если мы хотим узнать настройки воркеров, делаем GET-запрос по адресу:
```bash
http://localhost:8080/setup-workers
```

Ответ получаемый от сервера в формате JSON, структура:
```json
[
    {
        "worker_name": "worker2",
        "last_task": "70d9df9f-84ec-426a-b408-31fe7fe8380f",
        "status": "ready",
        "last_timeout_setup": "2024-02-17T21:10:21.870463Z",
        "current_timeout": 11
    },
    {
        "worker_name": "worker3",
        "last_task": "abb1f26f-2014-43f6-9f88-22fd6bf7aedc",
        "status": "ready",
        "last_timeout_setup": "2024-02-17T21:11:37.870463Z",
        "current_timeout": 12
    },
    {
        "worker_name": "worker1",
        "last_task": "f3715129-63f5-44a6-ac42-adbfb186ba60",
        "status": "ready",
        "last_timeout_setup": "2024-02-17T21:16:31.870463Z",
        "current_timeout": 10
    }
]
```

Пример GET - запроса
```bash
curl -X GET -H "Content-Type: application/json" http://localhost:8080/setup-workers
```

Если мы хотим поменять настройку у воркера, делаем POST-запрос по адресу:
```bash
http://localhost:8080/setup-workers
```

Тело запроса в формате JSON, структура:
```json
{
 "worker_name":"worker1",
 "timeout_data":10
}
```

Пример POST - запроса
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"worker_name\": \"worker1\", \"timeout_data\": 19}" http://localhost:8080/setup-workers
```

## 3.2. Работа в GUI
<a name="работа-в-GUI"></a>
После запуска проекта он доступен по адресу:
```bash
http://localhost:3000/
```

### Главная страница
На главной странице находится строка отправки выражения в которую мы вбиваем данные для вычисления и нажимаем кнопку "Отправить". Сигналом о том, что выражение отправлено, является всплывающее окно.

Также на главной странице располагается таблица состояния воркеров. Принцип действия этой таблицы заключается в следующем. Frontend делает регулярные запросы на главные страницы воркеров, если страница доступна, в таблице отображается "Online" и время последнего запроса. Иначе отображается Offline. 

Ссылки на работу воркеров, активные, ведут на их главные страницы

### Список выражений
На этой вкладке отображаются все выражения, которые были посчитаны воркерами в упрощенном виде. Здесь мы можем получить уникальный идентификатор и запрос с результатом


### Список операций
На этой вкладке показаны расширенные таблицы всех операций: 
```
UUID запроса, Запрос, Время создания, Время завершения, Время выполнения, Сервер, Результат, Статус
```

Тут же реализована пагинация, и возможность перейти на конкретную операцию (по uuid)

### Вычислительные мощности
Состоит трех блоков. Смена таймаута на воркерах, статус воркеров и получение сведений о текущем состоянии воркеров. 

## 4. Схема работы Worker
<a name="схема-работы-worker"></a>
[Схема Worker](https://github.com/iliazaraysky/distributed-expression-calculator/blob/main/%D0%A1%D1%85%D0%B5%D0%BC%D0%B0%20%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D1%8B%20Workers.png)

## 5. Схема работы Backend
<a name="схема-работы-backend"></a>
Поменять таймаут очень просто, выбираем в выпадающем списке нужного нам воркера, вводим ниже число, нажимаем кнопку "Отправить".

Текущие настройки воркеров, это постоянный GET запрос на backend, `curl -X GET -H "Content-Type: application/json" http://localhost:8080/setup-workers`, увидеть изменения в настройках воркеров можно довольно быстро. Не уверен, что так делают в реальных проектах, делал исключительно для демонстрации возможностей
