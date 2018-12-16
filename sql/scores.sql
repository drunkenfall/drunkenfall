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
   ORDER BY m.id, p.index);
  END $$
LANGUAGE plpgsql;

-- Same as above, but always from the most recent tournament
CREATE OR REPLACE FUNCTION scores() RETURNS
  TABLE(tid INTEGER, match_id integer, kind match_kind, nick text, color color, kills smallint, shots smallint, self smallint, sweeps smallint, total_score INTEGER) AS $$
  DECLARE tid INTEGER;
  BEGIN
   SELECT t.ID FROM tournaments t ORDER BY ID DESC LIMIT 1 INTO tid;\
   RETURN QUERY SELECT * FROM scores(tid);
  END $$
LANGUAGE plpgsql;

CREATE TABLE points (
       ID SERIAL PRIMARY KEY,
       kind text,
       score INTEGER
);

INSERT INTO points (kind, score) VALUES
  ('kill', 147),
  ('self', -245),
  ('sweep', 679),
  ('winner', 2450),
  ('second', 1050),
  ('third', 490),
  ('fourth', 210);

CREATE OR REPLACE FUNCTION update_summary(tid INTEGER, pid TEXT) RETURNS void AS $$
  BEGIN
   UPDATE player_summaries ps
   SET (kills, sweeps, self, shots, matches, total_score, skill_score)
   =
   (SELECT SUM(p.kills),
           SUM(p.sweeps),
           SUM(p.self),
           SUM(p.shots),
           COUNT(*),
           calculate_score(tid, pid, SUM(p.kills)::INTEGER, SUM(p.sweeps)::INTEGER, SUM(p.self)::INTEGER),
          (calculate_score(tid, pid, SUM(p.kills)::INTEGER, SUM(p.sweeps)::INTEGER, SUM(p.self)::INTEGER) / COUNT(*))
      FROM players p
      INNER JOIN matches m ON p.match_id = m.id
      WHERE m.tournament_id = tid
        AND m.started IS NOT NULL
        AND person_id = pid)
    WHERE person_id = pid
      AND tournament_id = tid;

   RETURN;
  END $$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION calculate_score(tid INTEGER, pid TEXT, kills INTEGER, sweeps INTEGER, self INTEGER) RETURNS INTEGER AS $$
  BEGIN
   RETURN (kills *  (SELECT score FROM points WHERE kind = 'kill')) +
          (sweeps * (SELECT score FROM points WHERE kind = 'sweep')) +
          (SELF *   (SELECT score FROM points WHERE kind = 'self')) +
          match_scores(tid, pid);
  END $$
LANGUAGE plpgsql;

-- Get the scores based on the final position in matches
CREATE OR REPLACE FUNCTION match_scores(tid INTEGER, pid TEXT) RETURNS INTEGER AS $$
  DECLARE first INTEGER;
  DECLARE second INTEGER;
  DECLARE third INTEGER;
  DECLARE fourth INTEGER;

  BEGIN
   SELECT count_positions(tid, pid, 1) INTO first;
   SELECT count_positions(tid, pid, 2) INTO second;
   SELECT count_positions(tid, pid, 3) INTO third;
   SELECT count_positions(tid, pid, 0) INTO fourth;

   RETURN (first *  (SELECT score FROM points WHERE kind = 'winner')) +
          (second * (SELECT score FROM points WHERE kind = 'second')) +
          (third *  (SELECT score FROM points WHERE kind = 'third')) +
          (fourth * (SELECT score FROM points WHERE kind = 'fourth'));
  END $$
LANGUAGE plpgsql;

-- Because of the nature of the 1-indexed ROW_NUMBER, the positions
-- make sense until the last; 4. Whoever came last will have the pos
-- 0, so the calling functions will need to take that into
-- consideration as well.
--
-- This makes sure to only consider matches that have ended, so this
-- will not affect maches that are ongoing.
CREATE OR REPLACE FUNCTION count_positions(tid INTEGER, pid TEXT, pos INTEGER) RETURNS INTEGER AS $$
  BEGIN
    RETURN (SELECT COUNT(person_id)
      FROM (SELECT ROW_NUMBER() OVER(ORDER BY p.match_id, p.kills DESC, p.SELF) AS n, p.person_id
              FROM players p WHERE p.match_id IN (
                  SELECT ID FROM matches
                   WHERE tournament_id = tid
                     AND ended IS NOT NULL)
      ) AS t WHERE t.n % 4 = pos AND person_id = pid);
  END $$
LANGUAGE plpgsql;
