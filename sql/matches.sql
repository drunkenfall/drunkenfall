CREATE OR REPLACE FUNCTION reset_match(mid INTEGER) RETURNS void AS $$
  BEGIN
    -- Reset the start date of the match
    UPDATE matches SET started = NULL WHERE id = mid;

    -- Reset all the scores
    UPDATE players SET shots = 0,
                      sweeps = 0,
                       kills = 0,
                        self = 0,
                 total_score = 0,
                 match_score = 0
     WHERE match_id = mid;

    -- Reset all player states to the default
    UPDATE player_states SET arrows = DEFAULT, shield = DEFAULT, wings = DEFAULT,
                             hat = DEFAULT, invisible = DEFAULT, speed = DEFAULT,
                             alive = DEFAULT, lava = DEFAULT, killer = DEFAULT
     WHERE player_id IN (SELECT id FROM players WHERE match_id = mid);

  END;$$
LANGUAGE plpgsql;
