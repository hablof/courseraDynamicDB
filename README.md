# Простое CRUD веб-приложение для MySQL-базы
 
## Динамическая структура БД
При старте приложения считывается список таблиц и полей из указанной базы.

## Валидация
Основываясь на считаной схеме происходит валидация данных при создании и редактировании записей.
Проверяется возможность установки `null`, а также соответствие типов `string`, `float`, `int`.
  
## API
+  **GET**  `/` - возвращает список всех таблиц
+  **GET**  `/table?limit=5&offset=7` - возвращает список из `limit` записей начиная с `offset` из таблицы `table`
+  **GET**  `/table/id` - возвращает информацию о самой записи `id` из таблицы `table`
+  **PUT**  `/table` - создаёт новую запись в таблице `table`
+  **POST**  `/table/id` - обновляет запись
+  **DELETE**  `/$table/$id` - удаляет запись
  
Данные для создания и редактирования записей считываются из тела запроса в формате `x-www-form-urlencoded`. Значение `null` кодируется как `%00`.

Выходные данные отсылаются в формате `json`, поля с `null`-значениями не отсылаются

Используется порт `:8082`
  
## Архитектура
Архитектурно код разделён на 4 слоя:
+  ***router*** - разбирает входящие http-запросы и передаёт данные на слой сервиса, а также возвращает ошибки и данные
+  ***repository*** - формирует запросы и обращается к базе данных
+  ***dbexplorer*** - формирует схему базы данных исходя из информации, полученной от репозитория
+  ***service*** - основываясь на схеме базы данных, валидирует входящие данные
  
Слои разделены с помощью механизма интерфейсов.
  
## Тестирование
Код *router*, *repository*, *service* покрыт unit-тестами. Благодаря разделению с помощью интерфейсов, методы зависимых объектов реализованы с помощью Mock-объектов.
  
## Запуск
Для запуска приложения подготовлены:
+  ***docker-compose***, описывающий контейнер с приложением и контейнер с базой
+  ***makefile*** с инструкциями

* * *
P.S. Приложение создавалось как решение задачи из старого курса с платформы *coursera* "Разработка веб сервисов на Golang" от Василия Романова из Mail.ru. Текст задания в файле `DB Explorer _ Coursera.pdf`.