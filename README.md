[![codecov](https://codecov.io/gh/RomanAgaltsev/ya_gophermart/graph/badge.svg?token=QCM8F0QPAZ)](https://codecov.io/gh/RomanAgaltsev/ya_gophermart)
![golangci](https://github.com/RomanAgaltsev/ya_gophermart/actions/workflows/golangci-lint.yml/badge.svg)

# Накопительная система лояльности «Гофермарт»

Первый дипломный проект курса "Продвинутый Go-разработчик" в Яндекс Практикуме.

## О проекте

* Автор: Роман Агальцев
* Когорта: 35
* Почта: roman-agalcev@yandex.ru

## Переменные окружения

Для запуска приложения с использованием Docker Compose, необходимо подготовить .env файл со следующими переменными
окружения:

* RUN_ADDRESS - адрес и порт запуска сервиса
* DATABASE_URI - строка подключения к базе данных Postgres
* ACCRUAL_SYSTEM_ADDRESS - адрес системы расчет начислений
* SECRET_KEY - секретный ключ, используемый при аутентификации пользователей

Кроме этого, для инициализации базы данных приложения на Postgres, в файле переменных окружения необходимо дополнительно
определить переменные:

* POSTGRES_USER - пользователь postgres
* POSTGRES_PASSWORD - пароль postgres
* POSTGRES_DB - база данных postgres
* POSTGRES_APP_USER - пользователь базы данных приложения
* POSTGRES_APP_PASS - пароль пользователя базы данных приложения
* POSTGRES_APP_DB - база данных приложения

Запускать docker-compose run следует с указанием подготовленного .env файла.