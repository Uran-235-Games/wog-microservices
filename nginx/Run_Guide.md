# Запуск nginx сервера на windows:

## Установка и настройка
1. Установка nginx
    1. Скачать zip архив с сайта nginx
    2. Распаковать например в `C:/nginx-1.28.0`
2. Taskfile
    1. Настроить переменную `NGINX_FOLDER_PATH`  
        - Пример: `NGINX_FOLDER_PATH: "C:/nginx-1.28.0"`
    2. Настроить переменную `PROJECT_NGINX_CONF_PATH` в команде `nginx-reload`  
        - Пример: `PROJECT_NGINX_CONF_PATH: "D:/sklad/txt/World-Of-Go/wog-microservices/nginx/nginx.conf"`

## Команды в Taskfile для работы с nginx:
- `task nginx-stop` (сокращение: `task rs`)  
    Завершает все процессы nginx
- `task nginx-reload` (сокращение: `task nr`)  
    Запускает nginx с обновленной конфигурацией
    - копирует nginx.conf из директории проекта в папку nginx
    - валидирует синтаксис nginx.conf
    - запускает nginx
- `task nginx -- {args}`  
    Команда вызывает nginx из нужной директории и передает в нее все аргументы команды указанные после `--`  
    Пример: `task nginx -- -s reload` -> выполнит команду `.\nginx.exe -s reload` из нужной директории