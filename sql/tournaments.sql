-- Clean up the join tables for levels and nicknames once a tournament
-- is done
CREATE OR REPLACE FUNCTION cleanup_after_tournament() RETURNS TRIGGER AS $$
 BEGIN
   DELETE FROM tournament_levels WHERE tournament_id = NEW.id;
   DELETE FROM tournament_names WHERE tournament_id = NEW.id;
   RETURN NEW;
 END;$$
LANGUAGE plpgsql;

-- Set the trigger to only happen when the ended date is set
CREATE TRIGGER on_tournament_end_cleanup
         AFTER UPDATE ON tournaments
           FOR EACH ROW
          WHEN (NEW.ended IS NOT NULL)
       EXECUTE PROCEDURE cleanup_after_tournament();
