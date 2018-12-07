CREATE TABLE player_states (
  ID SERIAL PRIMARY KEY,
  player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
  arrows    INTEGER[],
	shield    BOOLEAN NOT NULL DEFAULT FALSE,
	wings     BOOLEAN NOT NULL DEFAULT FALSE,
	hat       BOOLEAN NOT NULL DEFAULT TRUE,
	invisible BOOLEAN NOT NULL DEFAULT FALSE,
	speed     BOOLEAN NOT NULL DEFAULT FALSE,
	alive     BOOLEAN NOT NULL DEFAULT TRUE,
	lava      BOOLEAN NOT NULL DEFAULT FALSE,
	killer    INTEGER NOT NULL DEFAULT -2
);

CREATE INDEX player_states_id_idx ON player_states(id);

-- Trigger function to be used on insertion of players. Makes player_state objects for the player
-- that has just been added.
CREATE OR REPLACE FUNCTION auto_add_player_state() RETURNS trigger AS $$
  DECLARE tid INTEGER;
  BEGIN
    -- First grab the existing data from the person
    INSERT INTO player_states (player_id, index) VALUES (NEW.id, NEW.index);
    RETURN NEW;
  END;$$
LANGUAGE plpgsql;

CREATE TRIGGER on_player_make_player_state AFTER INSERT ON players FOR EACH ROW EXECUTE PROCEDURE auto_add_player_state();

-- Show the current state of a match
CREATE OR REPLACE FUNCTION states(mid INTEGER) RETURNS
  TABLE(ID INTEGER, m INTEGER, pid INTEGER, nick TEXT, color color, arrows INTEGER[], shield BOOLEAN, wings BOOLEAN, HAT BOOLEAN, invisible BOOLEAN, speed BOOLEAN, alive BOOLEAN, lava BOOLEAN, killer int) AS $$
  BEGIN
   RETURN QUERY (
      SELECT ps.id, mid, p.id, p.nick, p.color, ps.arrows, ps.shield, ps.wings, ps.hat, ps.invisible, ps.speed, ps.alive, ps.lava, ps.killer
        FROM player_states ps
       INNER JOIN players P ON p.ID = ps.player_id
       INNER JOIN matches M ON p.match_id = m.id
       WHERE m.id = mid
       ORDER BY ps.index
   );
  END $$
LANGUAGE plpgsql;

-- Same as above, but get the latest current match
CREATE OR REPLACE FUNCTION states() RETURNS
  TABLE(ID INTEGER, mid INTEGER, pid INTEGER, nick TEXT, color color, arrows INTEGER[], shield BOOLEAN, wings BOOLEAN, HAT BOOLEAN, invisible BOOLEAN, speed BOOLEAN, alive BOOLEAN, lava BOOLEAN, killer int) AS $$
  DECLARE tid INTEGER;
  DECLARE mid INTEGER;
  BEGIN
   SELECT t.ID FROM tournaments t ORDER BY ID DESC LIMIT 1 INTO tid;
   SELECT m.ID FROM matches m WHERE m.tournament_id = tid AND ended IS NULL ORDER BY ID ASC LIMIT 1 INTO mid;
   RETURN QUERY SELECT * FROM states(mid);
  END $$
LANGUAGE plpgsql;
