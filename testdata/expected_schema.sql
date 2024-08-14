-- Code generated by github.com/xshkut/go-comments-extractor, DO NOT EDIT.

-- source: file://./pkg/file.go

CREATE TYPE file_source AS ENUM ('aws_s3', 'google_slides', 'text');

CREATE TABLE file (
    id SERIAL PRIMARY KEY,
    source file_source NOT NULL,
    key TEXT NOT NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- source: file://./pkg/page.go

CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    file_id INTEGER REFERENCES files(id) NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE task_page (
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    page_id INTEGER NOT NULL REFERENCES pages(id),
    position INTEGER NOT NULL,
    PRIMARY KEY (task_id, page_id, position)
);

CREATE UNIQUE INDEX idx_page_title ON page (title);

-- source: file://./pkg/task.go

CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL, -- unique
    description TEXT NOT NULL,
    age_min INTEGER NOT NULL DEFAULT 0,
    age_max INTEGER NOT NULL DEFAULT 100,
    duration_sec INTEGER NOT NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_task_title ON task (title);

