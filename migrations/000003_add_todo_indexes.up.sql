CREATE INDEX IF NOT EXISTS todo_task_name_idx ON notes USING GIN(to_tsvector('simple', task_name));
CREATE INDEX IF NOT EXISTS todo_priority_idx ON notes USING GIN(to_tsvector('simple', priority));
CREATE INDEX IF NOT EXISTS todo_status_idx ON notes USING GIN(status);