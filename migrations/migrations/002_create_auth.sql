-- Удаляем триггер
DROP TRIGGER IF EXISTS trigger_set_auth_updated_at ON auth;

-- Удаляем функцию
DROP FUNCTION IF EXISTS set_auth_updated_at();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_auth_user_id;
DROP INDEX IF EXISTS idx_auth_email;
DROP INDEX IF EXISTS idx_auth_created_at;
DROP INDEX IF EXISTS idx_auth_updated_at;

DROP INDEX IF EXISTS idx_auth_tokens_user_id;
DROP INDEX IF EXISTS idx_auth_tokens_expires_at;

-- Удаляем таблицы
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS auth;