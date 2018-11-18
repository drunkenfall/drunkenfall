-- Returns an array of two strings, to be displayed on the character
-- select screen in the game.
--
-- Players with two words in their name (be it via spaces or
-- CamelCase) will get just those. Others will get a random prefix or
-- suffix (50/50 chance of each).
--
-- Assumes that player nicks are only allowed a maximum of one space
-- and that input validation handles those cases.
CREATE OR REPLACE FUNCTION display_names(tid INTEGER, nick TEXT) RETURNS TABLE (words TEXT []) AS $$
  DECLARE clean_nick TEXT;
  DECLARE other TEXT;
  BEGIN
    -- Clean the nick by un-camelcasing and replacing all double spaces
    SELECT TRIM(regexp_replace(regexp_replace(nick, '([^ A-Z][^A-Z]*)','\1 '), '[^0-9A-Za-z]+', ' ', 'g')) INTO clean_nick;

    -- When we already have two other (contains a space), use those straight up
    IF clean_nick LIKE '% %' THEN
      RETURN QUERY SELECT string_to_array(clean_nick, ' ');
    ELSE
      -- If we have one word only, randomly add a prefix or a suffix
      IF random() * 100 >= 50 THEN
        -- Prefix
        SELECT random_name FROM random_name(tid, TRUE) INTO other;
        RETURN QUERY SELECT ARRAY[other, clean_nick];
      ELSE
        -- Suffix
        SELECT random_name FROM random_name(tid, FALSE) INTO other;
        RETURN QUERY SELECT ARRAY[clean_nick, other];
      END IF;
    END IF;
  END;$$
LANGUAGE plpgsql;

-- Trigger function to be used on insertion of players. Generates
-- display names as by the function defined above, or grabs from the
-- person object if already set.
CREATE OR REPLACE FUNCTION player_insert() RETURNS trigger AS $$
  DECLARE nick TEXT;
  DECLARE dn TEXT [];
  DECLARE tid INTEGER;
  BEGIN
    -- First grab the existing data from the person
    SELECT m.tournament_id FROM matches M WHERE id = NEW.match_id INTO tid;
    SELECT p.nick FROM people p WHERE p.person_id = NEW.person_id INTO nick;
    SELECT display_names FROM people WHERE person_id = NEW.person_id INTO dn;

    IF dn IS NOT NULL THEN
      -- If the display names are set on the player, just use those
      NEW.display_names = dn;
    ELSE
      -- If not, grab random ones
      NEW.display_names = display_names(tid, nick);
    END IF;

    RETURN NEW;
  END;$$
LANGUAGE plpgsql;

CREATE TRIGGER on_player_set_display_names BEFORE INSERT ON players FOR EACH ROW EXECUTE PROCEDURE player_insert();

-- Fisher-Yates for name prefixes and suffixes. See `random_level` in
-- levels.sql for a more in-depth explanation of what's going on.
CREATE OR REPLACE FUNCTION random_name(tid INTEGER, aprefix bool) RETURNS text AS $$
 DECLARE nid INTEGER;
 DECLARE remaining SMALLINT;

 BEGIN
   SELECT COUNT(*)
     FROM tournament_names tn
    INNER JOIN names n ON n.id = tn.name_id
    WHERE tournament_id = tid
      AND n.prefix = aprefix
     INTO remaining;

   IF remaining = 0 THEN
     INSERT INTO tournament_names (tournament_id, name_id)
          SELECT tid, n.id FROM names n
           WHERE n.prefix = aprefix
           ORDER BY random();
   END IF;

   DELETE FROM tournament_names tn
         WHERE id IN (
               SELECT tn.id FROM tournament_names tn
                INNER JOIN names n ON n.id = tn.name_id
                WHERE tn.tournament_id = tid
                  AND n.prefix = aprefix
             ORDER BY tn.id
                LIMIT 1
         )
         RETURNING tn.name_id
   INTO nid;

   RETURN (SELECT name FROM names WHERE id = nid);
 END;$$
LANGUAGE plpgsql;

ALTER TABLE NAMES ADD COLUMN ID SERIAL PRIMARY KEY;
CREATE TABLE tournament_names (
    ID SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(ID) ON DELETE CASCADE,
    name_id INTEGER REFERENCES names(ID)
);
CREATE INDEX tournament_id_idx ON tournament_names(tournament_id);
