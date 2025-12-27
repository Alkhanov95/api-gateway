CREATE TABLE IF NOT EXISTS public.users (
  id   UUID PRIMARY KEY,
  name TEXT NOT NULL,
  age  INT  NOT NULL CHECK (age >= 0)
);
