-- First, add a new column to the counters_users table
ALTER TABLE counters_users 
ADD COLUMN entry_count INTEGER DEFAULT 0;

-- Create a function to update the count
CREATE OR REPLACE FUNCTION update_counter_entry_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment the entry_count in the counters_users table
    UPDATE counters_users
    SET entry_count = entry_count + 1
    WHERE user_id = NEW.user_id AND counter_id = NEW.counter_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function after insert
CREATE TRIGGER increment_counter_entry_count
AFTER INSERT ON counters_users_events
FOR EACH ROW EXECUTE FUNCTION update_counter_entry_count();
