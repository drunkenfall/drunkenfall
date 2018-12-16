-- Loop through all the messages for a match and assign them to their rounds accordingly
CREATE OR REPLACE FUNCTION bind_messages_to_rounds(mid INTEGER) RETURNS void AS $$
  DECLARE rid INTEGER;
  DECLARE M messages%rowtype;

  BEGIN
    -- Store the first ID of the rounds
    SELECT ID FROM rounds WHERE match_id = mid ORDER BY ID LIMIT 1 INTO rid;

    -- Loop all the messages of the match
    FOR M IN SELECT * FROM messages WHERE match_id = mid ORDER BY TIMESTAMP LOOP
      -- Update the current message to the round that we are on
      UPDATE messages SET round_id = rid WHERE id = m.ID;

      -- If we've just updated a `round_end`, then we're about to go to the next round
      -- so we increase the current round_id
      IF M.TYPE = 'round_end' THEN
         SELECT rid+1 INTO rid;
      END IF;
    END LOOP;

    -- Lastly, set the match_end to be on the same round as the round_end
    UPDATE messages SET round_id = rid-1 WHERE id = m.ID AND TYPE = 'match_end';
  END;$$
LANGUAGE plpgsql;


-- Updates the data in a `round` based on the messages that are in it
CREATE OR REPLACE FUNCTION update_round(rid INTEGER) RETURNS void AS $$
  BEGIN
    UPDATE rounds
       SET p1up = sum_kills(rid, 0), p1down = sum_self(rid, 0),
           p2up = sum_kills(rid, 1), p2down = sum_self(rid, 1),
           p3up = sum_kills(rid, 2), p3down = sum_self(rid, 2),
           p4up = sum_kills(rid, 3), p4down = sum_self(rid, 3)
     WHERE ID = rid;
  END;$$
LANGUAGE plpgsql;

-- Get the amount of kills a player has done, and disregard selfs
CREATE OR REPLACE FUNCTION sum_kills(rid INTEGER, P INTEGER) RETURNS INTEGER AS $$
  BEGIN
    RETURN (SELECT COUNT(*) FROM messages
            WHERE round_id = rid
              AND TYPE = 'kill'
              AND CAST(json->>'player' AS INTEGER) != P
              AND CAST(json->>'killer' AS INTEGER) = P);
  END;$$
LANGUAGE plpgsql;

-- Get the amount of selfs a player has done in a round.
-- Can only be one or zero, but for the sake of the data model it is kept as an
-- integer. Parts of the interface expects it to be numbers.
CREATE OR REPLACE FUNCTION sum_self(rid INTEGER, P INTEGER) RETURNS INTEGER AS $$
  DECLARE C INTEGER;
  BEGIN
    SELECT COUNT(*) FROM messages
            WHERE round_id = rid
              AND TYPE = 'kill'
              AND CAST(json->>'player' AS INTEGER) = p
              AND CAST(json->>'killer' AS INTEGER) = p
      INTO C;
    RETURN C;
  END;$$
LANGUAGE plpgsql;
