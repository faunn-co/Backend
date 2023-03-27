UPDATE booking_slots_table
SET citizen_slot = citizen_slot + ?, tourist_slot = tourist_slot + ?
WHERE date = ? AND slot = ? AND citizen_slot <= 50 AND tourist_slot <= 50