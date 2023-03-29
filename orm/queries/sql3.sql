SELECT
    r.*
FROM
    referral_table r
WHERE
        r.affiliate_id = ? AND r.referral_status != 4;
