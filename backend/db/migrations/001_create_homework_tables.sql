-- Создание таблицы домашних заданий
CREATE TABLE IF NOT EXISTS homeworks (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    student_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    rating INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- Создание таблицы файлов домашних заданий
CREATE TABLE IF NOT EXISTS homework_files (
    id SERIAL PRIMARY KEY,
    homework_id INTEGER NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (homework_id) REFERENCES homeworks(id) ON DELETE CASCADE
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_homeworks_student_id ON homeworks(student_id);
CREATE INDEX IF NOT EXISTS idx_homeworks_teacher_id ON homeworks(teacher_id);
CREATE INDEX IF NOT EXISTS idx_homework_files_homework_id ON homework_files(homework_id);

-- Создание триггера для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_homeworks_updated_at
    BEFORE UPDATE ON homeworks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 