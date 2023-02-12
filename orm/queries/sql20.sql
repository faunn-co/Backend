SELECT
    booking_id,
    booking_status,
    DATE_FORMAT(booking_day, '%d/%m/%Y') AS booking_day,
    booking_slot,
    transaction_time,
    payment_status,
    citizen_ticket_count,
    tourist_ticket_count,
    SUM(
                citizen_ticket_total + tourist_ticket_total
        ) AS ticket_total
FROM
    affiliate_manager_db.booking_details_table
WHERE
        transaction_time > ?
  AND transaction_time <= ?
GROUP BY
    booking_id
ORDER BY
    booking_id DESC;
