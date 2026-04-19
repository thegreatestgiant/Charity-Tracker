CREATE type entry AS enum ('paycheck', 'donation');

CREATE TABLE Ledgers (
  user_id INT NOT NULL,
  tranaction_id SERIAL PRIMARY KEY,
  ledger_entry entry,
  amount DECIMAL(18, 2) NOT NULL,
  charity_owed DECIMAL(18, 2),
  charity_fulfilled DECIMAL(18, 2),
  transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES Users (user_id)
);
