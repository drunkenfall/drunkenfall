-- Use a Fisher-Yates shuffle bucket of levels to determine what
-- levels to play on. Uses a volatile join table `tournament_levels`
-- as the bucket container, and inserts a new batch of shuffled levels
-- whenever the bucket is exhausted. Hence, the join table can be
-- truncated at any time to reset the randomness, and things will just
-- carry on.
--
-- Realistically, the `final` bucket will never be seen, because as of
-- forever it only has "cataclysm" in it. So, an empty bucket will be
-- seen, a single item will be inserted, and it will immediately be
-- consumed.
--
-- The join table also makes the data completely normalized, which
-- gives warm fuzzy feelings to those that care about such things.
CREATE OR REPLACE FUNCTION random_level(tid INTEGER, akind match_kind) RETURNS level AS $$
 DECLARE lid INTEGER;
 DECLARE remaining SMALLINT;

 BEGIN
   -- First, count how many levels we have left for the tournament in
   -- the current bucket
   SELECT COUNT(*)
     FROM tournament_levels tl
    INNER JOIN levels l ON l.id = tl.level_id
    WHERE tournament_id = tid
      AND l.kind = akind
     INTO remaining;

   -- If the bucket is empty, then we should enter a full new set of
   -- shuffled levels into the joining table.
   --
   -- Essentially, putting everything in the hat.
   IF remaining = 0 THEN
     INSERT INTO tournament_levels (tournament_id, level_id)
          SELECT tid, l.ID FROM levels l
           WHERE l.kind = akind
           ORDER BY random();
   END IF;

   -- Delete the top row from the current bucket, and return the level
   -- id into the `lid` variable.
   --
   -- Essentially, picking a piece from the hat.
   DELETE FROM tournament_levels tl
         WHERE id IN (
               SELECT tl.id FROM tournament_levels tl
                INNER JOIN levels l ON l.id = tl.level_id
                WHERE tl.tournament_id = tid
                  AND l.kind = akind
             ORDER BY tl.id
                LIMIT 1
         )
         RETURNING tl.level_id
   INTO lid;

   -- Open the piece and read what you got!
   RETURN (SELECT level FROM levels WHERE id = lid);
 END;$$
LANGUAGE plpgsql;

-- Set the ruleset based on the match kind.
--
-- Almost too simple, but it's hella nice to have it automatic.
CREATE OR REPLACE FUNCTION level_rules(akind match_kind) RETURNS ruleset AS $$
 BEGIN
  IF akind = 'final' THEN
     RETURN 'B';
  ELSE
     RETURN 'A';
  END IF;
 END;$$
LANGUAGE plpgsql;

-- Sets a level and a ruleset on a match when we are inserting it
CREATE OR REPLACE FUNCTION match_insert() RETURNS trigger AS $$
  BEGIN
   NEW.level = random_level(NEW.tournament_id, NEW.kind);
   NEW.ruleset = level_rules(NEW.kind);
   RETURN NEW;
  END;$$
LANGUAGE plpgsql;

CREATE TRIGGER on_match_insert_ruleset BEFORE INSERT ON matches FOR EACH ROW EXECUTE PROCEDURE match_insert();
