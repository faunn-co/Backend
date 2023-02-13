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
    affiliate_manager_db.referral_table r
        LEFT JOIN affiliate_manager_db.user_table u ON r.affiliate_id = u.user_id
        LEFT JOIN affiliate_manager_db.booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_id = ?;