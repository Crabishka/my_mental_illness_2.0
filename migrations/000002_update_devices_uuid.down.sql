-- Возвращаем id колонку
ALTER TABLE devices 
    ADD COLUMN id SERIAL PRIMARY KEY;

-- Удаляем ограничения
ALTER TABLE devices 
    DROP CONSTRAINT devices_token_unique,
    DROP CONSTRAINT devices_uuid_unique,
    DROP CONSTRAINT devices_pkey;

-- Удаляем uuid колонку
ALTER TABLE devices 
    DROP COLUMN uuid; 