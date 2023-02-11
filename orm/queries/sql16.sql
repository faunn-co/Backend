SELECT
    a.affiliate_id,
    a.user_name AS affiliate_name,
    a.affiliate_type,
    a.unique_referral_code,
    COUNT(a.referral_click_time) AS referral_clicks,
    COALESCE(
                SUM(citizen_ticket_total) + SUM(tourist_ticket_total),
                0
        ) AS total_revenue,
    COUNT(a.booking_id) AS referral_count,
    COALESCE(
            SUM(a.referral_commission),
            0
        ) AS referral_commission
FROM
    (
        SELECT
            a.affiliate_type,
            a.unique_referral_code,
            u.user_name,
            u.user_email,
            u.user_contact,
            r.affiliate_id,
            IF(
                            r.referral_click_time <= @startTime
                        OR r.referral_click_time > @endTime,
                            NULL,
                            r.referral_click_time
                ) AS referral_click_time,
            IF(
                            r.booking_time <= @startTime
                        OR r.booking_time > @endTime,
                            0,
                            r.referral_commission
                ) AS referral_commission,
            IF(
                            r.booking_time <= @startTime
                        OR r.booking_time > @endTime,
                            NULL,
                            b.booking_id
                ) AS booking_id,
            IF(
                            r.booking_time <= @startTime
                        OR r.booking_time > @endTime,
                            0,
                            b.citizen_ticket_total
                ) AS citizen_ticket_total,
            IF(
                            r.booking_time <= @startTime
                        OR r.booking_time > @endTime,
                            0,
                            b.tourist_ticket_total
                ) AS tourist_ticket_total
        FROM
            affiliate_manager_db.affiliate_details_table a
                LEFT JOIN affiliate_manager_db.user_table u ON a.user_id = u.user_id
                LEFT JOIN affiliate_manager_db.referral_table r ON a.user_id = r.affiliate_id
                LEFT JOIN affiliate_manager_db.booking_details_table b ON b.booking_id = r.booking_id
        GROUP BY
            a.user_id,
            u.user_id,
            r.affiliate_id,
            b.booking_id,
            r.referral_id
    ) AS a
GROUP BY
    a.affiliate_id,
    a.user_name,
    a.user_email,
    a.user_contact,
    a.affiliate_type,
    a.unique_referral_code;
