package pkg

/* SQL:
CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    file_id INTEGER REFERENCES files(id) NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
*/

/* SQL:
CREATE TABLE task_page (
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    page_id INTEGER NOT NULL REFERENCES pages(id),
    position INTEGER NOT NULL,
    PRIMARY KEY (task_id, page_id, position)
);
*/

// All code below was generated by AI using the content above.

// Page represents the page entity.
type Page struct {
	ID int `json:"id"`
	/* SQL:
	CREATE UNIQUE INDEX idx_page_title ON page (title);
	*/
	Title     string `json:"title"`
	FileID    int    `json:"file_id,omitempty"`
	CreatedBy int    `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

// TaskPage represents the task_page entity.
type TaskPage struct {
	TaskID   int `json:"task_id"`
	PageID   int `json:"page_id"`
	Position int `json:"position"`
}
