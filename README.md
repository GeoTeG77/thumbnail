# Thumbnail (Тестовое задание)

Программа для получения preview youtube video по ссылке на него.

## Оглавление

- [Описание](#описание)
- [Функциональные возможности](#функциональные-возможности)
- [Установка](#установка)
- [Использование](#использование)
---

## Описание

Приложение получает preview youtube video по ссылке формата https://www.youtube.com/watch?v=<id_video> через GRPC-proxy.

## Функциональные возможности

- Запуск grpc клиента для получения thumbnail требуемого видео (1);
- Клиент поддерживает два режима получения данных синхронный и асинхронный (флаг --async)(2);
- Кеширование полученных изоображений в Redis сроком на 5min по ключу url-video (3);
- Возможность работать со списком формата (4);
---

## Установка

Для установки приложения необходимо выполнить следующие шаги:
1. Установить Docker и Docker-compose.
2. Установить git.
2. Открыть терминал в требуемой папке и скопировать проект с gitHub
git clone https://github.com/GeoTeG77/thumbnail.git
4. Запустить в этом же терминале команду docker-compose up --build

## Использование

1. После успешного запуска приложения вы увидите следующие сообщения в терминале:
myapp  | time=2024-12-11T17:08:51.456Z level=INFO msg="Repository layer successfully create!"      
myapp  | time=2024-12-11T17:08:51.456Z level=INFO msg="Service layer successfully create!"
myapp  | time=2024-12-11T17:08:51.456Z level=INFO msg="GRPC-server started successfully"

Теперь вы можете войти внутрь контейнера при помощи: docker exec -it myapp /bin/sh
2. После этого рекомендую запустить следующие команды для тестирования:
/usr/local/bin/client --async https://www.youtube.com/watch?v=tPiagp9t5is
/usr/local/bin/client --async https://www.youtube.com/watch?v=tPiagp9t5is,https://www.youtube.com/watch?v=dmx_8jo0eqE
/usr/local/bin/client https://www.youtube.com/watch?v=tPiagp9t5is,https://www.youtube.com/watch?v=dmx_8jo0eqE
или любое кол-во видео в формате
/usr/local/bin/client --async https://www.youtube.com/watch?v=<id_video>,https://www.youtube.com/watch?v=<id_video>,...
/usr/local/bin/client https://www.youtube.com/watch?v=<id_video>,https://www.youtube.com/watch?v=<id_video>,...

3. При правильной работе программы вы увидите подобные сообщения в Stdout:
myapp  | time=2024-12-11T17:35:57.768Z level=ERROR msg="Cache miss for URL" url="https://www.youtube.com/watch?v=dmx_8jo0eqE"
myapp  | time=2024-12-11T17:35:57.797Z level=INFO msg="Saved thumbnail in Cache for url:" url="https://www.youtube.com/watch?v=dmx_8jo0eqE"
myapp  | time=2024-12-11T17:35:57.797Z level=INFO msg="Successfully take thumbnail from YouTube" url="https://www.youtube.com/watch?v=dmx_8jo0eqE"
myapp  | time=2024-12-11T17:36:01.697Z level=INFO msg="Received request" count=2
myapp  | time=2024-12-11T17:36:01.697Z level=INFO msg="Thumbnail Taken from Cache"
myapp  | time=2024-12-11T17:36:01.697Z level=INFO msg="Successfully take thumbnail from YouTube" url="https://www.youtube.com/watch?v=tPiagp9t5is"
