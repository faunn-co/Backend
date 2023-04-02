SELECT
    COUNT(r.referral_click_time) AS referral_clicks,
    COUNT(r.booking_id) AS referral_count,
    COALESCE(
            SUM(r.referral_commission),
            0
        ) AS referral_commission,
    b.total_revenue
FROM
    referral_table r,
    (
        SELECT
            COALESCE(
                        SUM(b.citizen_ticket_total) + SUM(b.tourist_ticket_total),
                        0
                ) AS total_revenue
        FROM
            booking_details_table b WHERE b.transaction_time <= ?
    ) AS b WHERE r.referral_click_time <= ? AND r.booking_time <= ?
GROUP BY
    total_revenue