-- Apply inside an isolated test schema. READ COMMITTED with serialized counter writers.
CREATE TABLE business(id integer PRIMARY KEY CHECK(id=1), value bigint NOT NULL);
INSERT INTO business VALUES(1,0);
CREATE TABLE requests(scope text NOT NULL, key text NOT NULL, amount bigint NOT NULL CHECK(amount>0), response bigint NOT NULL, PRIMARY KEY(scope,key));
CREATE TABLE outgoing(scope text NOT NULL, key text NOT NULL, amount bigint NOT NULL, sent boolean NOT NULL DEFAULT false, PRIMARY KEY(scope,key));
CREATE TABLE inbox(scope text NOT NULL, key text NOT NULL, amount bigint NOT NULL, PRIMARY KEY(scope,key));
CREATE TABLE effects(scope text NOT NULL, key text NOT NULL, amount bigint NOT NULL, PRIMARY KEY(scope,key));
