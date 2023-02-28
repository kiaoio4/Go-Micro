package micro

// ElementKey is type of key of element (not component)
type ElementKey string

// LoggingElementKey is ElementKey for logging
var LoggingElementKey = ElementKey("Go-Micro/LoggingComponent")

// LoggerGroupElementKey is ElementKey for LoggerGroup
var LoggerGroupElementKey = ElementKey("Go-Micro/LoggerGroupComponent")

// TracingElementKey is ElementKey for tracing
var TracingElementKey = ElementKey("Go-Micro/TracingComponent")

// NacosClientElementKey is ElementKey for nacos client
var NacosClientElementKey = ElementKey("Go-Micro/NacosClient")

// MysqlElementKey is ElementKey for mysql
var MysqlElementKey = ElementKey("Go-Micro/MysqlComponent")

// RedisElementKey is ElementKey for redis
var RedisElementKey = ElementKey("Go-Micro/RedisComponent")

// GossipKVCacheElementKey is ElementKey for GossipKVCache
var GossipKVCacheElementKey = ElementKey("Go-Micro/GossipKVCacheComponent")
