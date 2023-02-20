package repository

const usersTableSchema = `
CREATE TABLE if not exists users_table
(
    id serial not null unique,
    login varchar(255) not null unique,
    password_hash varchar(255) not null
);`

var Index = `
CREATE UNIQUE INDEX if not exists login_index_unique
  ON users_table
  USING btree(login);
`

const txLogSchema = `
CREATE TABLE if not exists transaction_log
(
    User_id int references users_table(id) on delete cascade not null,
	Sender_wallet_id integer,
    Receiver_wallet_id integer,
	TxID varchar (40),
	Status varchar (10),
	Amount integer,
    Note varchar(50),
    Initiated_at timestamp
   
);`
