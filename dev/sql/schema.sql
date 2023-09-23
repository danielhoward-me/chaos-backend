CREATE TABLE saves (
	id SERIAL PRIMARY KEY,
	user_id UUID,
	name TEXT NOT NULL,
	data JSON NOT NULL,
	screenshot TEXT
);
CREATE INDEX idx_user_id ON saves (user_id);
