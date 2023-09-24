CREATE OR REPLACE FUNCTION insert_data_hash() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.screenshot = md5(NEW.data::text);
    RETURN NEW; 
END;
$$ language 'plpgsql';

CREATE TABLE saves (
	id SERIAL PRIMARY KEY,
	user_id UUID,
	name TEXT NOT NULL,
	data JSON NOT NULL,
	screenshot TEXT
);
CREATE INDEX idx_user_id ON saves (user_id);
CREATE TRIGGER set_screenshot_value BEFORE INSERT ON saves FOR EACH ROW EXECUTE PROCEDURE insert_data_hash();
