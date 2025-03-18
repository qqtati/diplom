-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL,
    invite_code VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы связей учитель-ученик
CREATE TABLE IF NOT EXISTS teacher_student (
    id SERIAL PRIMARY KEY,
    teacher_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES "user"(id) ON DELETE CASCADE,
    UNIQUE(teacher_id, student_id)
);

-- Создание таблицы событий
CREATE TABLE IF NOT EXISTS "event" (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    duration INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    student_id INTEGER NOT NULL,
    description TEXT,
    approved_by_teacher BOOLEAN DEFAULT FALSE,
    skipped BOOLEAN DEFAULT FALSE,
    rating INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_user_username ON "user"(username);
CREATE INDEX IF NOT EXISTS idx_user_invite_code ON "user"(invite_code);
CREATE INDEX IF NOT EXISTS idx_teacher_student_teacher_id ON teacher_student(teacher_id);
CREATE INDEX IF NOT EXISTS idx_teacher_student_student_id ON teacher_student(student_id);
CREATE INDEX IF NOT EXISTS idx_event_teacher_id ON "event"(teacher_id);
CREATE INDEX IF NOT EXISTS idx_event_student_id ON "event"(student_id);
CREATE INDEX IF NOT EXISTS idx_event_start_time ON "event"(start_time);

-- Создание триггеров для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_event_updated_at
    BEFORE UPDATE ON "event"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 