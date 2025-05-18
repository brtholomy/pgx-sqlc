create table users(
   id uuid PRIMARY KEY NOT NULL,
   email VARCHAR NOT NULL,
   name VARCHAR NOT NULL
);

create table products(
   id uuid PRIMARY KEY NOT NULL,
   user_id uuid REFERENCES users ON DELETE CASCADE,
   name VARCHAR NOT NULL,
   price NUMERIC(100, 2) NOT NULL
);

create table invoices(
   id uuid PRIMARY KEY NOT NULL,
   user_id uuid REFERENCES users ON DELETE CASCADE,
   invoice_number VARCHAR NOT NULL,
   total NUMERIC(100, 2) NOT NULL,
   created TIMESTAMPTZ DEFAULT now(),
   modified TIMESTAMPTZ DEFAULT now()
);

create table invoice_items(
   id uuid PRIMARY KEY NOT NULL,
   user_id uuid REFERENCES users ON DELETE CASCADE,
   product_id uuid REFERENCES products ON DELETE RESTRICT,
   invoice_id uuid REFERENCES invoices ON DELETE CASCADE,
   amount NUMERIC(100, 2) NOT NULL
);
