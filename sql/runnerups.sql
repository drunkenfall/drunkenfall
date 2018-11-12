-- Prints the list of players eligible for runnerup in the given
-- tournament. Excludes those already booked into matches that haven't
-- ended yet (i.e. scheduled matches)
CREATE OR REPLACE FUNCTION runnerups(tid INTEGER) RETURNS
  TABLE(ID INTEGER, tournament_id integer, nick text, shots smallint, sweeps SMALLINT, kills smallint, self smallint, matches SMALLINT, total_score INTEGER, skill_score INTEGER) AS $$
  BEGIN
   RETURN QUERY (SELECT ps.ID, ps.tournament_id, p.nick, ps.shots, ps.sweeps, ps.kills, ps.self, ps.matches, ps.total_score, ps.skill_score
    FROM player_summaries ps
    INNER JOIN people p ON p.person_id = ps.person_id
   WHERE ps.tournament_id = tid
     AND ps.person_id NOT IN (
            SELECT p.person_id
              FROM players P INNER JOIN matches M ON m.ID = p.match_id
              WHERE m.ended IS NULL AND m.tournament_id = tid
         )
   ORDER BY ps.matches ASC, ps.skill_score DESC);
  END $$
LANGUAGE plpgsql;

-- Same as above, but for latest tournament only
CREATE OR REPLACE FUNCTION runnerups() RETURNS
  TABLE(ID INTEGER, tournament_id integer, nick text, shots smallint, sweeps SMALLINT, kills smallint, self smallint, matches SMALLINT, total_score INTEGER, skill_score INTEGER) AS $$
  DECLARE tid INTEGER;
  BEGIN
   SELECT t.ID FROM tournaments t ORDER BY ID DESC LIMIT 1 INTO tid;
   RETURN QUERY SELECT * FROM runnerups(tid);
  END $$
LANGUAGE plpgsql;
