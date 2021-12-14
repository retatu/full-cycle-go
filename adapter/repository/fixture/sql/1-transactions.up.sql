CREATE TABLE transactions(
  id            text not null,
  account_id    text not null,
  amount        real not null,
  status        text not null,
  error_message text not null,
  created_at    text not null,
  updated_at    text not null
);