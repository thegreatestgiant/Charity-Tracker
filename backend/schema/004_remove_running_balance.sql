ALTER TYPE entry
RENAME TO entry_old;

CREATE type entry AS enum ('paycheck', 'donation');

ALTER TABLE Ledgers
ALTER COLUMN ledger_entry Type entry USING ledger_entry::text::entry;

DROP TYPE entry_old;
