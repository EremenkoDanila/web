
-- Вставка данных в таблицу users
INSERT INTO users (login, password, moderator_flg) VALUES
('admin', 'admin123', true),    -- Администратор (модератор)
('client', 'client123', false); -- Обычный клиент


-- Вставка данных в таблицу software (из orders)
INSERT INTO softwares (deleted_flg, title, description, full_description, img_url, os, size, version, functions) VALUES
(false, 'GIT', 'GIT - система контроля версий', 'Git — система контроля версий для отслеживания изменений в файлах и совместной работы.', 'http://127.0.0.1:9000/lb1/Git_icon.png', 'Linux', 1.5, '2.45', '1. Создание репозитория\n2. Добавление изменений в индекс\n3. Фиксация изменений\n4. Просмотр статуса и истории изменений\n5. Ветвление\n6. Слияние веток'),
(false, 'Open Shift', 'Open Shift - управление контейнерами', 'Платформа контейнеризации и оркестрации приложений на базе Kubernetes.', 'http://127.0.0.1:9000/lb1/open_shift.png', 'Linux', 18.2, '4.15', '1. Оркестрация контейнеров\n2. Автоматическое масштабирование\n3. Непрерывное развертывание\n4. Мониторинг и логирование\n5. Управление сетями\n6. Балансировка нагрузки'),
(false, 'PostgreSQL', 'PostgreSQL - система управления БД', 'Надёжная объектно-реляционная СУБД с открытым исходным кодом.', 'http://127.0.0.1:9000/lb1/Postgre_SQL.png', 'Linux', 22.3, '15.1', '1. Подключение к различным СУБД\n2. Редактор SQL запросов\n3. Визуальное построение запросов\n4. Управление схемами БД\n5. Экспорт и импорт данных\n6. Генерация ER-диаграмм'),
(false, 'DBeaver', 'DBeaver - инструмент для работы с БД', 'Многоплатформенный инструмент для управления и проектирования СУБД.', 'http://127.0.0.1:9000/lb1/dbeaver.png', 'Linux', 22.3, '15.1', '1. Подключение к различным СУБД\n2. Редактор SQL запросов\n3. Визуальное построение запросов\n4. Управление схемами БД\n5. Экспорт и импорт данных\n6. Генерация ER-диаграмм');

-- Вставка данных в таблицу software_request (из cart_data)
-- Предполагается, что есть пользователи с ID 1 (CreatorID) и 2 (ModeratorID)
INSERT INTO software_requests (status, create_dt, update_dt, finish_dt, creator_id, moderator_id, phone) VALUES
('completed', NOW(), NOW(), NULL, 2, 1, '8-800-555-35-35');

-- Вставка данных в таблицу installation_time (из cart_to_service)
INSERT INTO installation_times (software_id, request_id, install_time, final_time) VALUES
(1, 1, now(), now()),
(2, 1, now(), now());



черновик - draft
удалён - deleted
сформирован - formed
завершён - completed
отклонён - rejected