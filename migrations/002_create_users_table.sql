-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем колонку user_id в таблицу tasks
ALTER TABLE tasks ADD COLUMN user_id INTEGER REFERENCES users(id);

-- Создаем индекс для быстрого поиска задач пользователя
CREATE INDEX idx_tasks_user_id ON tasks(user_id); 