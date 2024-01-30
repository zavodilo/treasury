up:
	docker-compose -f docker-compose.yml up -d

build:
	docker-compose -f docker-compose.yml up -d --build

migrate:
	docker-compose -f docker-compose.yml exec go go run ./src/migrations/*.go init
	docker-compose -f docker-compose.yml exec go go run ./src/migrations/*.go up

down:
	docker-compose down

stop:
	docker-compose stop

restart:
	docker-compose restart

console:
	docker exec -it treasury bash
