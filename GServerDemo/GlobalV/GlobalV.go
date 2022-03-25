package GlobalV

import (
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoGlobalV *mongo.Client
	RedisGlobalV *redis.Pool
)

