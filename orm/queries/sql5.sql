SELECT
    COALESCE(SUM(r.referral_commission),0) AS total_commission ,
    COUNT(DISTINCT r.affiliate_id) AS total_active_affiliates ,
    COUNT(r.referral_id) AS total_affiliate_bookings,
    COALESCE(SUM(b.citizen_ticket_total),0) AS citizen_ticket_total,
    COALESCE(SUM(b.tourist_ticket_total),0) AS tourist_ticket_total
FROM
    referral_table r
        LEFT JOIN booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_status = 0 AND r.affiliate_id IS NOT NULL AND r.referral_status != 4
  AND r.booking_time >= ?
  AND r.booking_time <= ?;
