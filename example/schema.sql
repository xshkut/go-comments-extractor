-- Generated SQL Schema. Do not edit.

-- source: ../go/file.go

CREATE TYPE file_source AS ENUM ('aws_s3', 'google_slides', 'text');
CREATE TABLE file (
    id SERIAL PRIMARY KEY,
    source file_source NOT NULL,
    key TEXT NOT NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- source: ../go/page.go

CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    file_id INTEGER REFERENCES files(id) NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- source: ../go/task.go

CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    age_min INTEGER NOT NULL DEFAULT 0,
    age_max INTEGER NOT NULL DEFAULT 100,
    duration_sec INTEGER NOT NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
