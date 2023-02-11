SELECT
    COUNT(r.referral_id) AS total_clicks
FROM
    affiliate_manager_db.referral_table r
WHERE
        r.referral_click_time > UNIX_TIMESTAMP(
            CONCAT(@startTime, ' 00:00:00')
        )
  AND r.referral_click_time <= UNIX_TIMESTAMP(
        CONCAT(@endTime, ' 23:59:59')
    )
  AND r.affiliate_id = @id;
