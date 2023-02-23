SELECT
    r.referral_id,
    r.affiliate_id,
    u.user_name,
    r.referral_click_time,
    r.referral_status,
    r.booking_id,
    r.referral_commission,
    b.*
FROM
    referral_table r
        LEFT JOIN user_table u ON r.affiliate_id = u.user_id
        LEFT JOIN booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_id = ?;