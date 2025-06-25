CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS goods (
    id SERIAL,
    project_id INT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    priority INT NOT NULL,
    removed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, project_id)
);

CREATE INDEX idx_goods_id ON goods(id);
CREATE INDEX idx_goods_project_id ON goods(project_id);
CREATE INDEX idx_goods_name ON goods(name);

INSERT INTO projects (name) VALUES ('Первая запись');
