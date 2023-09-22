CREATE TABLE saves (
	id SERIAL PRIMARY KEY,
	user_id UUID NOT NULL,
	data JSON NOT NULL
);
CREATE INDEX idx_user_id ON saves (user_id);
