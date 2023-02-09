SELECT
    COALESCE(SUM(r.referral_commission),0) AS total_commission ,
    COUNT(r.affiliate_id) AS total_active_affiliates ,
    COUNT(r.referral_id) AS total_affiliate_bookings,
    COALESCE(SUM(b.citizen_ticket_total),0) AS citizen_ticket_total,
    COALESCE(SUM(b.tourist_ticket_total),0) AS tourist_ticket_total
FROM
    affiliate_manager_db.referral_table r
        LEFT JOIN affiliate_manager_db.booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_status = 0
  AND r.booking_time >= ?
  AND r.booking_time <= ?;
