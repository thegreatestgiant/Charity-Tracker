ALTER TABLE ledgers
DROP CONSTRAINT fk_user_id;

ALTER TABLE Users
ALTER COLUMN user_id
DROP DEFAULT;

ALTER TABLE Users
ALTER COLUMN user_id TYPE uuid USING gen_random_uuid ();

ALTER TABLE ledgers
ALTER COLUMN user_id TYPE uuid USING gen_random_uuid ();

ALTER TABLE ledgers ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES Users (user_id);
