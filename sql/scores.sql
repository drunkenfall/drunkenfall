-- Prints a human friendly list (i.e. nicks instead of person_ids) of
-- all matches in a tournament
CREATE OR REPLACE FUNCTION scores(tid INTEGER) RETURNS
  TABLE(id INTEGER, match_id integer, kind match_kind, nick text, color color, kills smallint, shots smallint, self smallint, sweeps smallint, total_score INTEGER) AS $$
  BEGIN
   RETURN QUERY (SELECT tid, m.id, m.kind, ps.nick, p.color, p.kills, p.shots, p.self, p.sweeps, p.total_score
    FROM players P
    INNER JOIN matches m ON p.match_id = m.ID
    INNER JOIN people ps ON p.person_id = ps.person_id
   WHERE m.tournament_id = tid
   ORDER BY m.id);
  END $$
LANGUAGE plpgsql;

-- Same as above, but always from the most recent tournament
CREATE OR REPLACE FUNCTION scores() RETURNS
  TABLE(tid INTEGER, match_id integer, kind match_kind, nick text, color color, kills smallint, shots smallint, self smallint, sweeps smallint, total_score INTEGER) AS $$
  DECLARE tid INTEGER;
  BEGIN
   SELECT t.ID FROM tournaments t ORDER BY ID DESC LIMIT 1 INTO tid;
   RETURN QUERY SELECT * FROM scores(tid);
  END $$
LANGUAGE plpgsql;
