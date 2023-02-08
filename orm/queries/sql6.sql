SELECT
    a.*,
    u.user_name as affiliate_name,
    COUNT(referral_id) AS total_referrals
FROM
    affiliate_manager_db.referral_table r,
    affiliate_manager_db.affiliate_details_table a,
    affiliate_manager_db.user_table u
WHERE
        r.referral_status = 0
  AND r.affiliate_id = a.user_id
  AND r.affiliate_id = u.user_id
  AND r.booking_time >= ?
  AND r.booking_time <= ?
GROUP BY
    affiliate_id
ORDER BY
    COUNT(referral_id) DESC
LIMIT
    5;
