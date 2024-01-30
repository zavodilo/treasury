# Реализовал алгоритм более быстрого обновления данных при повторном вызове
* Тк записей в таблице у нас немного, то мы получаем все записи из таблицы.
* Обходим все Individual
* Если uid нет в таблице, то сохраняем его
* Сравниваем по uid с тем, что было получено из таблицы в полях first_name, last_name, если не совпадает, то пересохраняем

# Запуск приложения
make build
make migrate

# curl commands  

curl -X GET localhost:8080/state

# Migrations 
Migrations work using https://github.com/go-pg/migrations. 
For apply a migration   in folder cmd/migrations run
`` go run *.go <arguments>``

Currently, the following arguments are supported:

* ``up`` - runs all available migrations;
* ``up [target] `` - runs available migrations up to the target one;
* ``down `` - reverts last migration;
* ``reset`` - reverts all migrations;
* ``version`` - prints current db version;
* ``set_version [version] `` - sets db version without running migrations.

