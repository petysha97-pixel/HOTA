-- ===== Этап 1: изменения схемы (запускать вручную в SQLite) =====

-- пользователь: ФИО и грейд
ALTER TABLE users ADD COLUMN Name  TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN Grade TEXT NOT NULL DEFAULT '';

-- слот: название, описание (техзадание), исполнитель
ALTER TABLE slots ADD COLUMN name        TEXT NOT NULL DEFAULT '';
ALTER TABLE slots ADD COLUMN description TEXT NOT NULL DEFAULT '';
-- при удалении пользователя исполнитель слота сбрасывается в NULL
ALTER TABLE slots ADD COLUMN user_id     INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- отклик: сообщение создателю проекта
ALTER TABLE applicationsSlot ADD COLUMN message TEXT NOT NULL DEFAULT '';
-- когда последний раз менялся статус: от него считается пауза перед повторным откликом
ALTER TABLE applicationsSlot ADD COLUMN updated_at DATETIME;

-- приглашения: создатель проекта зовёт разработчика в слот
CREATE TABLE invites (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    slot_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    status     TEXT    NOT NULL DEFAULT 'pending',
    message    TEXT    NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- каталог ролей (как stacks)
CREATE TABLE roles (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO roles (name) VALUES
 ('Backend'),('Frontend'),('Fullstack'),('DevOps'),('Mobile'),('iOS'),('Android'),
 ('QA'),('QA Automation'),('Design'),('UI/UX'),('Data Science'),('Data Engineer'),
 ('ML Engineer'),('GameDev'),('Security'),('Embedded'),('SRE'),('1C'),
 ('Системный аналитик'),('Tech Lead'),('Архитектор'),('Product Manager');
