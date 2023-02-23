SELECT Count(r.referral_id) AS previous_cycle_referrals
FROM   referral_table r
WHERE  r.affiliate_id = ?
  AND r.booking_time >= ?
  AND r.booking_time <= ?
GROUP  BY affiliate_id;