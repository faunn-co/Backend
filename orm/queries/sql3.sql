SELECT
    r.*
FROM
    affiliate_manager_db.referral_table r
WHERE
        r.affiliate_id = ?;
