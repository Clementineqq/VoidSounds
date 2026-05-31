# VoidSounds
 go run ./cmd/server
 # или 
 air 
 # если air не запущен при добавлении шаблонов
 templ generate


# migr
go run migrations/migrate.go

# собрать образы  (сервисы в фоне)
docker-compose up -d

# пересобрать образы без кеша (при изменениях в Dockerfile)
docker-compose build --no-cache

# Остановить все контейнеры без удаления
docker-compose stop

# остановить и удалить контейнеры но сохранить тома с бд	
docker-compose down

# удалить контейнеры и тома с данными (сброс бд)	
docker-compose down -v

# env 
DB_HOST=postgres для докера
DB_HOST=localhost для натива


# TODO
фикс оргов в карточке детальной, фикс фильтров, подумать над картами и админ панелью, фикс редактирование мероприятий