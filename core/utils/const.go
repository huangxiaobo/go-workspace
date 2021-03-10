package utils

type DomainType int

const (
	KUAIDAILI DomainType = 1 << iota
	XICI

	REDIS_RAW_PROXY_LIST = "proxy_raw_list"
)
