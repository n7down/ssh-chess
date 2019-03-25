CREATE TABLE players 
(
 	id SERIAL, 
	username VARCHAR(100) PRIMARY KEY NOT NULL,
	secret VARCHAR(100) NOT NULL
    -- CONSTRAINT PK_Person PRIMARY KEY (id, username)
	-- PRIMARY KEY (id, username)
);

CREATE TABLE games
(
	id uuid PRIMARY KEY NOT NULL, 
	name VARCHAR(100) NOT NULL,
	black_player_id int NOT NULL,
	white_player_id int NOT NULL,
	-- black_player_name VARCHAR(100) NOT NULL,
	-- white_player_name VARCHAR(100) NOT NULL,
	start_time VARCHAR(100) NOT NULL,
	end_time VARCHAR(100) DEFAULT NULL,
	outcome VARCHAR(100) DEFAULT NULL,
	pgn VARCHAR(5000) NOT NULL
);

