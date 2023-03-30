SELECT
    a.*,
    u.user_name as affiliate_name,
    SUM(r.referral_commission) AS total_commission
FROM
    referral_table r,
    affiliate_details_table a,
    user_table u
WHERE
        r.referral_status = 0
  AND r.affiliate_id = a.user_id
  AND r.affiliate_id = u.user_id
  AND r.booking_time >= ?
  AND r.booking_time <= ?
  AND r.referral_status != 4
GROUP BY
    affiliate_id
ORDER BY
    SUM(r.referral_commission) DESC
LIMIT
    5;
