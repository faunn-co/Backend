SELECT
    DATE_FORMAT(d.date_string, '%m-%d') AS date_string,
    COALESCE(
            SUM(r.referral_commission),
            0
        ) AS total_commission,
    COUNT(DISTINCT r.affiliate_id) AS total_active_affiliates,
    COUNT(r.referral_id) AS total_affiliate_bookings,
    COALESCE(
            SUM(b.citizen_ticket_total),
            0
        ) AS citizen_ticket_total,
    COALESCE(
            SUM(b.tourist_ticket_total),
            0
        ) AS tourist_ticket_total,
    COUNT(r.referral_id) AS total_affiliate_bookings
FROM
    date_ref_table d
        LEFT JOIN referral_table r ON d.date_string = DATE(
            FROM_UNIXTIME(r.booking_time)
        )
        LEFT JOIN booking_details_table b ON r.booking_id = b.booking_id
WHERE
    d.date_string BETWEEN ?
        AND ?
GROUP BY
    d.date_string
ORDER BY
    d.date_string;
