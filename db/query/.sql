
SELECT INTO auth_audit (
		user_id,
		event,
		event_time
	) VALUES (
		$1, $2, $3
);