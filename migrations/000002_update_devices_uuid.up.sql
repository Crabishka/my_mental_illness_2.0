-- Создаем временную колонку для UUID
ALTER TABLE devices ADD COLUMN uuid VARCHAR(36);

-- Заполняем UUID для существующих записей
UPDATE devices SET uuid = gen_random_uuid()::text WHERE uuid IS NULL;

-- Делаем UUID NOT NULL и уникальным
ALTER TABLE devices 
    ALTER COLUMN uuid SET NOT NULL,
    ADD CONSTRAINT devices_uuid_unique UNIQUE (uuid);

-- Делаем token уникальным
ALTER TABLE devices 
    ADD CONSTRAINT devices_token_unique UNIQUE (token);

-- Удаляем старый primary key и id колонку
ALTER TABLE devices 
    DROP CONSTRAINT devices_pkey,
    ADD PRIMARY KEY (uuid),
    DROP COLUMN id; 