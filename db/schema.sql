CREATE TABLE IF NOT EXISTS workers (
	 id uuid NOT NULL PRIMARY KEY,
	 name text NOT NULL
);
CREATE INDEX IF NOT EXISTS workers_name_idx ON workers(name COLLATE NOCASE);

CREATE TABLE IF NOT EXISTS shifts (
	 id uuid NOT NULL PRIMARY KEY,
	 worker_id text NOT NULL,
	 date date NOT NULL,
	 start_hour tinyint NOT NULL,
	 end_hour tinyint NOT NULL,
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
);
CREATE INDEX IF NOT EXISTS shifts_date_worker_id_idx ON shifts(date, worker_id);