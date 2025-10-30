-- Очистка в порядке зависимостей (сначала дочерние таблицы, потом родительские)
TRUNCATE TABLE installation_times;
TRUNCATE TABLE software_requests;
TRUNCATE TABLE softwares;
TRUNCATE TABLE users;