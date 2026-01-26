-- Drop the trigger first
DROP TRIGGER IF EXISTS increment_counter_entry_count ON counters_users_events;

-- Drop the function
DROP FUNCTION IF EXISTS update_counter_entry_count();

-- Optional: Drop the decrement trigger and function if you created it earlier
DROP TRIGGER IF EXISTS decrement_counter_entry_count ON counters_users_events;
DROP FUNCTION IF EXISTS decrement_counter_entry_count();

-- Remove the entry_count column from the counters_users table
ALTER TABLE counters_users 
DROP COLUMN IF EXISTS entry_count;
