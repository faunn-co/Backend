SELECT
    SUM(r.referral_commission) AS total_commission,
    COUNT(r.affiliate_id) AS total_active_affiliates,
    COUNT(r.referral_id) AS total_affiliate_bookings,
    SUM(b.citizen_ticket_total) AS citizen_ticket_total,
    SUM(b.tourist_ticket_total) AS tourist_ticket_total
FROM
    affiliate_manager_db.referral_table r
        LEFT JOIN affiliate_manager_db.booking_details_table b ON r.booking_id = b.booking_id
WHERE
        r.referral_status = 0
  AND r.booking_time >= ?
  AND r.booking_time <= ?;
