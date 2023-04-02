SELECT
    a.user_id as affiliate_id,
    a.entity_name,
    a.entity_identifier,
    a.affiliate_type,
    a.unique_referral_code,
    COUNT(r.referral_click_time) AS referral_clicks,
    COALESCE(
                SUM(b.citizen_ticket_total) + SUM(b.tourist_ticket_total),
                0
        ) AS total_revenue,
    COUNT(r.booking_id) AS referral_count,
    COALESCE(
            SUM(r.referral_commission),
            0
        ) AS referral_commission
FROM
    affiliate_details_table a
        LEFT JOIN user_table u ON a.user_id = u.user_id
        LEFT JOIN referral_table r ON a.user_id = r.affiliate_id
        LEFT JOIN booking_details_table b ON b.booking_id = r.booking_id
WHERE
        a.user_id = ? AND r.referral_click_time <= ? AND r.booking_time <= ? AND b.transaction_time <= ?
GROUP BY a.user_id