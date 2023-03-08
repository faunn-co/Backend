SELECT
    DISTINCT date_string,
    IF(
                a.affiliate_id <> @id, 0, total_commission
        ) AS total_commission,
    IF(
                a.affiliate_id <> @id, 0, total_bookings
        ) AS total_bookings,
    IF(
                a.affiliate_id <> @id, 0, citizen_ticket_total
        ) AS citizen_ticket_total,
    IF(
                a.affiliate_id <> @id, 0, tourist_ticket_total
        ) AS tourist_ticket_total
FROM
    (
        SELECT
            r.affiliate_id,
            date_string,
            COALESCE(
                    SUM(r.referral_commission),
                    0
                ) AS total_commission,
            COUNT(r.referral_id) AS total_bookings,
            COALESCE(
                    SUM(b.citizen_ticket_total),
                    0
                ) AS citizen_ticket_total,
            COALESCE(
                    SUM(b.tourist_ticket_total),
                    0
                ) AS tourist_ticket_total
        FROM
            date_ref_table d
                LEFT JOIN (SELECT * FROM referral_table WHERE affiliate_id= @id) AS r ON d.date_string = DATE(
                    FROM_UNIXTIME(r.booking_time)
                )
                LEFT JOIN booking_details_table b ON r.booking_id = b.booking_id
        WHERE
            d.date_string BETWEEN @startTime
                AND @endTime
        GROUP BY
            d.date_string,
            r.affiliate_id
        ORDER BY
            d.date_string
    ) a;