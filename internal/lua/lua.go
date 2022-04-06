package lua

const (
	IPLimitPeriod        = 60
	IPLimitMaximum int64 = 60
	Script               = `
		local count = redis.call('incr', KEYS[1])
		if count == 1 then
		redis.call('expire', KEYS[1], tonumber(KEYS[2]))
		end
		local reset = redis.call('ttl', KEYS[1])
		return {
			count,
			reset
		}
	`
)
