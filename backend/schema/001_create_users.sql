CREATE TABLE Users (
  user_id UUID PRIMARY KEY,
  email text,
  username text NOT NULL,
  password_hash VARCHAR(72) NOT NULL,
  date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  donation_percentage INT DEFAULT 10
)
