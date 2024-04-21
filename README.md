# Распределенный вычислитель арифметических выражений

## Оглавление

1. [Описание проекта](#описание-проекта)

2. [Установка и настройка](#установка-и-настройка)

 	2.1. [Установка Docker](#установка-docker)

	2.2 [Клонирование репозитория](#клонирование-репозитория)
 
 	2.3. [Запуск проекта](#запуск-проекта)
	
 	2.4. [Возможные ошибки](#возможные-ошибки)
	
 	2.5. [Два воркера не ошибка](#два-воркера-не-ошибка)
	
 	2.6. [Выключение проекта](#выключение-проекта)

3. [Пример использования](#пример-использования)
	
 	3.1. [Запросы в curl](#запросы-в-curl)

   	3.2. [Работа в GUI](#работа-в-GUI)

4. [Схема работы Worker](#схема-работы-worker)

5. [Схема работы Backend](#схема-работы-backend)

6. [Вторая часть. Инструкция](#вторая-часть-инструкция)

7. [Схема работы Регистрации/Авторизации](#схема-работы-reg-auth)

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

Проект ставится в несколько шагов:

0. Устанавливаем Docker (если нет)
1. Клонируем репозиторий
2. Переходим в папку с программой в командной строке
3. Запускаем Docker для создания образов и запуска контейнеров

### 2.1. Установка Docker
<a name="установка-docker"></a>
Для работы с сервисом необходимо установить на компьютер Docker

Версии для Mac, Windows, Linux доступны по следующему адресу:
https://docs.docker.com/get-docker/

1. Скачиваем установочный файл под вашу операционную систему
2. Запускаем / Устанавливаем / Перезагружаем компьютер

### 2.2. Клонирование репозитория
<a name="клонирование-репозитория"></a>

Клонируем себе на компьютер репозиторий. Набираем в терминале:
```bash
git clone https://github.com/iliazaraysky/distributed-expression-calculator.git
```

Переходим в папку с проектом на компьютере
```bash
cd distributed-expression-calculator
```

## 2.3. Запуск проекта
<a name="запуск-проекта"></a>
Для создания образов и запуска контейнеров набираем:
```bash
docker-compose up --build
```

Первый раз проект запускается долго. Нужно дождаться, чтобы скачались необходимые образы (я постарался прописать везде alpine образы с минимальным набором компонентов, но получилось все равно много). Если вы запускаете проект без флага `-d` , сигналом к началу проверки приложения, могут послужить сообщения в терминале, что воркеры ждут заданий из очереди:

```
worker1   | 2024/02/18 10:08:35 No messages in the queue, waiting
worker2   | 2024/02/18 10:08:37 No messages in the queue, waiting
worker1   | 2024/02/18 10:08:40 No messages in the queue, waiting
worker2   | 2024/02/18 10:08:42 No messages in the queue, waiting
worker1   | 2024/02/18 10:08:45 No messages in the queue, waiting
```

## 2.4. Возможные ошибки
<a name="возможные-ошибки"></a>
Я попробовал запускать проект на разных стендах, на 1 из 3 вылезла ошибка:
```
ERROR: for postgres  Cannot create container for service postgres: Conflict.
The container name "/postgres" is already in use by container "ac58e56714bbf18014776cfa3e47ae31d597db721cf00d0c0b36c0621dffd10a".
You have to remove (or rename) that container to be able to reuse that name.
```

Из текста сообщения понятно, в системе уже есть контейнер с названием `postgres`. Я исправил это тем, что удалил контейнер с этим названием командой:
```bash
docker rm postgres
```

Также один раз вылезла ошибка сборки frontend, только у меня и только после того как я зачистил буквально весь Docker от всех Containers, Images, Volumes.

Если встретите подобное, огромная просьба попробовать еще раз пересобрать проект.

Тестирование проводилось на следующих стендах:
1. Десктоп Intel(R) Core(TM) i5 CPU 760 2.80GHz 16,0 ГБ RAM ОС Windows 10 x64 22H2 Сборка 19045.4046
2. Десктоп Intel(R) Core(TM) i5 CPU 4440 3.10GHz 16,0 ГБ RAM ОС Windows 10 x64 22H2 Сборка 19045.4046
3. Ноутбук Intel(R) Core(TM) i3 CPU 7100U 2.40GHz 8,0 ГБ RAM ОС Ubuntu x64 22.04.3 LTS

## 2.5. Два воркера не ошибка
<a name="два-воркера-не-ошибка"></a>
Знаю, в проекте заявлено три воркера. В моей базе данных также предварительно записываются три воркера и в списке на таймауты можно выбрать третий.

Однако я решил его не добавлять, так как код в первых двух максимально схож, а вот ресурсы расходуются знатно. Я очень хочу, чтобы проект запустился и каждый кто будет проверять, не испытывал с этим трудности.

Заодно можно увидеть статус "Offline" в таблице монитонга воркеров. Это не заглушка, frontend честно пытается его пинговать

## 2.6. Выключение проекта
<a name="выключение-проекта"></a>
Набрать из консоли, находясь в той же директории откуда запускался проект

```bash
docker-compose down
```

## 3. Пример использования
<a name="пример-использования"></a>
В проекте реализован GUI. Удобнее всего посмотреть работу из него. Но первым пунктом описания сделаны примеры запросов, через `curl`. 
Во второй части проекты мы добавляли регистрацию и авторизацию. Поэтому в качестве примера, одна из ранее открытых областей, стала закрытой. Теперь нельзя получить доступ к "Списку операций" без токена в заголовке.

### 3.1 Регистрация пользователя
<a name="регистрация-пользователя"></a>

Страница отвечает на GET-запрос, приветствием. Пример запроса:
```bash
curl -X GET -H "Content-Type: application/json" http://127.0.0.1:8080/registration
```

Отправляем POST-запрос по адресу:
```bash
http://localhost:8080/registration
```
Тело запроса в формате JSON, структура:
```bash
{
  "login": "user",
  "password": "password123"
}
```

Пример запроса:
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"login\": \"user1\", \"password\":\"password123\"}" http://127.0.0.1:8080/registration
```

### 3.2. Авторизация. Получение токена
<a name="регистрация-пользователя"></a>

Страница отвечает на GET-запрос, приветствием. Пример запроса:
```bash
curl -X GET -H "Content-Type: application/json" http://127.0.0.1:8080/login
```

Отправляем POST-запрос по адресу:
```bash
http://localhost:8080/login
```

Пример запроса:
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"login\": \"user1\", \"password\":\"password123\"}" http://127.0.0.1:8080/login
```

В ответ сервер присылает Token. Пример токена:
```bash
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMyNzYwNzcsImlhdCI6MTcxMzI3NTQ3NywibG9naW4iOiJ1c2VyMSIsIm5iZiI6MTcxMzI3NTQ3N30.zseZouXKkOlNL7X7mmOrVIFtu-Ekp2nb0lZ7Wnw_SVE"}
```
### 3.3. Запросы в curl. Добавление вычисления арифметического выражения.
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

Так как во второй части задания у нас появилась авторизация. Простой GET-запрос, будет возвращать `Unauthorized`

Пример запроса без авторизации:
```bash
curl -X GET -H "Content-Type: application/json" http://127.0.0.1:8080/get-operations
```

Пример запроса с авторизацией
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMyNzY1MTEsImlhdCI6MTcxMzI3NTkxMSwibG9naW4iOiJ1c2VyMSIsIm5iZiI6MTcxMzI3NTkxMX0.GL1bbHRijfDo1BortkJMCMrxQwRxexHP8sHEX98bDaA" http://localhost:8080/get-operations
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

Поменять таймаут очень просто, выбираем в выпадающем списке нужного нам воркера, вводим ниже число, нажимаем кнопку "Отправить".

Текущие настройки воркеров, это постоянный GET запрос на backend, `curl -X GET -H "Content-Type: application/json" http://localhost:8080/setup-workers`, увидеть изменения в настройках воркеров можно довольно быстро. Не уверен, что так делают в реальных проектах, делал исключительно для демонстрации возможностей


## 4. Схема работы Worker
<a name="схема-работы-worker"></a>
![Схема Worker](https://github.com/iliazaraysky/distributed-expression-calculator/blob/main/%D0%A1%D1%85%D0%B5%D0%BC%D0%B0%20%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D1%8B%20Workers.png)

## 5. Схема работы Backend
<a name="схема-работы-backend"></a>
![Схема Backend](https://github.com/iliazaraysky/distributed-expression-calculator/blob/main/%D0%A1%D1%85%D0%B5%D0%BC%D0%B0%20%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D1%8B%20backend.png)

<a name="вторая-часть-инструкция"></a>
## 6. Вторая часть. Изминения. Инструкция
Для начала стоит указать полный список изменений внесенных во второй части проекта:

### На стороне Backend
1. Добавлен эндпойнт `/registration`
2. Добавлен эндпойнт `/login`
3. Добавлен эндпойнт `/get-operation-by-user-id`. Нужен для того, чтобы сортировать данные по конкретному пользователю
4. Добавлен **Middleware** `authMiddleware`, который проверяет наличие токена в заголовке `Authorization`
5. Изменен **Middleware** `corsHandler`, так как теперь у нас есть дополнительный заголовок `Authorization`
6. База данных, таблица `requests`, добавлен столбец `username`, в котором теперь записывается пользователь отправивший задание
7. База данных, добавлена таблица `users`, в которой хранятся `login` и `password`, необходимые для получения токена
8. Изменена `database_schema` для docker, с учетом появившейся таблицы `users`, и дополнительной колонки `username` в `requests`
9. Изменена схема работы **RabbitMQ**, так как теперь в очередь попадают не только сведения об операции, но и о пользователе, который инициировал операцию

### На стороне Frontend
1. Изменено меню. Добавлены кнопка `Регистрация`
2. Изменено меню. Добавлены кнопка `Логин`
3. Изменено меню. Когда пользователь авторизирован `Регистрация`, меняется на его `Логин`, который в свою очередь ведет на все операции пользователя
4. Создан компонент в **React**, отвечающий за страницу операций пользователя `UserRequestDetails` 
5. Создан компонент в **React**, отвечающий за регистрацию пользователя `Register`
6. Создан компонент в **React**, отвечающий за авторизацию пользователя `Login`
7. Страница `Список операций`, теперь доступна только авторизированным пользователям
8. Отправлять задания с Frontend могут только авторизированные пользователи
9. В `App.js` ведется проверка токена, контролируется время когда он истекает
10. При инициализации `App.js` из токена сохраняется `username` и передается на главную страницу, для последующего добавление при инициализации операции 

### Тесты
Файл с тестами лежит в директории `Backend`. Запускать лучше всего после запуска сервиса.
Полное описание тестов лежит в отдельном файле. [TESTS_DESCRIPTION.md](TESTS_DESCRIPTION.md)

### Инструкция v2.0
1. Скачиваем репозиторий 
	```bash
	git clone https://github.com/iliazaraysky/distributed-expression-calculator.git
	```
2. Переходим в скаченную директорию `distributed-expression-calculator`

3. Запускаем проект
	```bash
	docker-compose up -d --build
	```
4. После запуска проект доступен по адресу:
	```bash
	http://localhost:3000 
	```
5. Перейдя на главную страницу видим информацию о том, что отправка выражений доступна только после авторизации, также нет доступа к списку операций. Тут следует пояснить:
   5.1. На главной странице блок отправки выражений не показывается пока нет токена, это прописано на Frontend
   5.2. На странице `Список операций`, блокировка идет на уровне Backend, средствами **Middleware** `authMiddleware`
6. Обратим внимание на правый верхний угол **Регистрация** -> Вводим `Логин` и `Пароль` -> Создать. Так вы создадите пользователя, от имени которого можно будет совершать операции
   6.1. **Примечание**. Когда заполнены поля и нажата кнопка `Создать`, **React** делает редирект на страницу авторизации, но почему-то не всегда =). Frontend - это сложно, я с ним намучался
7. Переходим на страницу `Вход в кабинет` -> вводим регистрационные данные -> Вход
8. После авторизации, в правом верхнем углу будет оторажаться `Login` и `Выход`.
   8.1. `Выход` удаляет токен из локального хранилища, и обновляет страницу
   8.2. `Login` является ссылкой к списку операций пользователя
9. На странице `Список операций`, теперь доступны к просмотру сделанные ранее операции
10. На `Главной` странице теперь доступно поле ввода выражений
11. Вводим -> Отправляем -> нажимаем на `Логин` -> Видим появившуюся операцию
12. Нажимаем выход -> идем на главную страницу -> строка отправки недоступна -> идем в `Список операций` -> снова недоступно
13. Набираем в терминале
	```bash
	docker-compose down 
	```
14. Снова заходим под своей учетной записью, пользователь на месте, данные сохранены
15. Можно завершить работу `docker-compose` с флагом `-v`, тогда данные не сохранятся, а новый запуск будет с "чистым" сервисом

<a name="схема-работы-reg-auth"></a>
## 7. Схема работы Регистрации/Авторизации

![Схема Регистрации-Авторизации](https://github.com/iliazaraysky/distributed-expression-calculator/blob/main/%D0%A1%D1%85%D0%B5%D0%BC%D0%B0%20%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D1%8B%20%D0%A0%D0%B5%D0%B3%D0%B8%D1%81%D1%82%D1%80%D0%B0%D1%86%D0%B8%D0%B8-%D0%90%D0%B2%D1%82%D0%BE%D1%80%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D0%B8.png?raw=true)

### На этом все! Спасибо, что потратили свое время на проверку!
