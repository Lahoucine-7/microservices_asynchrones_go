CREATE TABLE IF NOT EXISTS commandes (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  product TEXT NOT NULL,
  amount FLOAT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);