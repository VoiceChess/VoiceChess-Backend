CREATE OR REPLACE FUNCTION set_move_order()
RETURNS TRIGGER AS $$
BEGIN
SELECT COALESCE(MAX(move_order), 0) + 1
INTO NEW.move_order
FROM moves
WHERE game_id = NEW.game_id;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER moves_set_order
    BEFORE INSERT ON moves
    FOR EACH ROW
    EXECUTE FUNCTION set_move_order();
