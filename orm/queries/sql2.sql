SELECT
    a.*,
    COUNT(r.referral_id) AS referral_count,
    u.user_name as affiliate_name
FROM
    affiliate_details_table a,
    referral_table r,
    user_table u
WHERE
        a.user_id = ?
  AND a.user_id = r.affiliate_id
  AND a.user_id = u.user_id
  AND r.referral_status != 4
GROUP BY
    r.affiliate_id;
