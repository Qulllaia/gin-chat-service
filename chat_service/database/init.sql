CREATE TABLE "Chat" (
    chat_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    chat_type VARCHAR(50),
    users BIGINT[],
    chat_background VARCHAR(500)
);
CREATE TABLE "ChatHistory" (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    "timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ChatType" (
    id VARCHAR(20) PRIMARY KEY,
    type_name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO "ChatType" (id, type_name) VALUES 
    ('GROUPCHAT', 'Group Chat'),
    ('PRIVATECHAT', 'Private Chat');

ALTER TABLE "Chat" DROP COLUMN chat_type;

-- Добавляем новую колонку с foreign key (varchar)
ALTER TABLE "Chat" ADD COLUMN chat_type_id VARCHAR(20) NOT NULL;

-- Создаем foreign key constraint
ALTER TABLE "Chat" 
ADD CONSTRAINT fk_chat_chat_type 
FOREIGN KEY (chat_type_id) REFERENCES "ChatType"(id);
