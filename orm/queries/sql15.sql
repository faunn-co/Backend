SELECT
    referral_id,
    booking_time,
    referral_commission,
    booking_id AS booking_ref_id
FROM
    affiliate_manager_db.referral_table
WHERE
        affiliate_id = ?
  AND booking_time > ?
  AND booking_time <= ?
ORDER BY
    booking_time DESC
LIMIT
    10;
