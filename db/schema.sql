create table bank(
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   first VARCHAR NOT NULL,
   last VARCHAR NOT NULL,
   email VARCHAR NOT NULL,
   amount NUMERIC(100, 2) DEFAULT 0,
   creation TIMESTAMPTZ DEFAULT now()
);
