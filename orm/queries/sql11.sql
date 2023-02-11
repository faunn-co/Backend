SELECT
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
        ) AS tourist_ticket_total,
    (
        SELECT
            COUNT(r.referral_id) AS total_clicks
        FROM
            affiliate_manager_db.referral_table r
        WHERE
                r.affiliate_id = @id
          AND r.referral_click_time >= @startTime
          AND r.referral_click_time <= @endTime
    ) AS total_clicks
FROM
    affiliate_manager_db.referral_table r
        LEFT JOIN affiliate_manager_db.booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_status = 0
  AND r.affiliate_id = @id
  AND r.booking_time >= @startTime
  AND r.booking_time <= @endTime;
