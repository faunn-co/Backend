SELECT
    u.user_name AS affiliate_name,
    r.referral_id,
    r.referral_click_time,
    r.referral_status,
    r.referral_commission,
    r.booking_id AS booking_ref_id,
    r.booking_time,
    SUM(
                b.citizen_ticket_count + b.tourist_ticket_count
        ) AS total_ticket_count,
    SUM(
                b.citizen_ticket_total + b.tourist_ticket_total
        ) AS total_ticket_amount
FROM
    affiliate_manager_db.referral_table r
        LEFT JOIN affiliate_manager_db.user_table u ON u.user_id = r.affiliate_id
        LEFT JOIN affiliate_manager_db.booking_details_table b ON b.booking_id = r.booking_id
        WHERE r.referral_click_time > ? AND referral_click_time <= ?
GROUP BY
    r.referral_id,
    r.referral_click_time
ORDER BY
    r.referral_id DESC,
    r.referral_click_time DESC;
