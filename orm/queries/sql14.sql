SELECT
    referral_id,
    referral_click_time,
    referral_status
FROM
    affiliate_manager_db.referral_table
WHERE
        affiliate_id = ?
  AND referral_click_time > ?
  AND referral_click_time <= ?
ORDER BY
    referral_click_time DESC
LIMIT
    10;