CREATE TABLE saves (
	id SERIAL PRIMARY KEY,
	user_id UUID,
	data JSON NOT NULL,
	screenshot TEXT
);
CREATE INDEX idx_user_id ON saves (user_id);
