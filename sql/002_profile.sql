-- ===== Профиль: описание опыта у каждой технологии пользователя =====
-- запускать вручную в SQLite после 001_stage1.sql

-- описание опыта с технологией (пишется в профиле после регистрации)
ALTER TABLE user_stacks ADD COLUMN description TEXT NOT NULL DEFAULT '';
