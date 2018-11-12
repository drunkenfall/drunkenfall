-- Show the total scores of a tourmament
CREATE OR REPLACE FUNCTION overview(tid INTEGER) RETURNS
  TABLE(ID INTEGER, tournament_id integer, nick text, shots smallint, sweeps SMALLINT, kills smallint, self smallint, matches SMALLINT, total_score INTEGER, skill_score INTEGER) AS $$
  BEGIN
   RETURN QUERY (SELECT ps.ID, ps.tournament_id, p.nick, ps.shots, ps.sweeps, ps.kills, ps.self, ps.matches, ps.total_score, ps.skill_score
    FROM player_summaries ps
    INNER JOIN people p ON p.person_id = ps.person_id
   WHERE ps.tournament_id = tid
   ORDER BY ps.total_score DESC);
  END $$
LANGUAGE plpgsql;

-- Same as above, but for latest tournament only
CREATE OR REPLACE FUNCTION overview() RETURNS
  TABLE(ID INTEGER, tournament_id integer, nick text, shots smallint, sweeps SMALLINT, kills smallint, self smallint, matches SMALLINT, total_score INTEGER, skill_score INTEGER) AS $$
  DECLARE tid INTEGER;
  BEGIN
   SELECT t.ID FROM tournaments t ORDER BY ID DESC LIMIT 1 INTO tid;
   RETURN QUERY SELECT * FROM overview(tid);
  END $$
LANGUAGE plpgsql;
