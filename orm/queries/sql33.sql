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
            booking_details_table b
    ) AS b
GROUP BY
    total_revenue