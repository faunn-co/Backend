SELECT
    a.*,
    COUNT(r.referral_id) AS referral_count,
    u.user_name as affiliate_name
FROM
    affiliate_manager_db.affiliate_details_table a,
    affiliate_manager_db.referral_table r,
    affiliate_manager_db.user_table u
WHERE
        a.user_id = r.affiliate_id
  AND a.user_id = u.user_id
GROUP BY
    r.affiliate_id;