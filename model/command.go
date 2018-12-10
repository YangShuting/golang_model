package model

import (
	"fmt"
	"strings"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/storage"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2/bson"
)

const CNAME_COMMAND = "commands"
const COMAMND_PREFIX = "command_"

type Command struct {
	ID      string `bson:"_id" json:"_id"`
	UserID  string `bson:"user_id" json:"user_id"`
	Command string `bson:"command" json:"command"`
}

func NewCommand(cc *lib.Cusctx, userID string, command string) (*Command, error) {
	commandIn := &Command{
		ID:      bson.NewObjectId().Hex(),
		UserID:  userID,
		Command: command,
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_COMMAND).Insert(commandIn)
	if err != nil {
		return nil, err
	}
	return commandIn, nil
}

func FindCommands(cc *lib.Cusctx, userID string, limit int) []*Command {
	var commands []*Command
	findC := bson.M{
		"user_id": userID,
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_COMMAND).Find(findC).Limit(limit).All(&commands)
	return commands
}

// redis 的应用：取出最新的 从  start 到 endNum 的数据。
//先从redis 读取， 如果读出来的小于 endNum 的话，就从数据库中读取数据。
//数据推入 list
func RedisListPush(listName string, values ...interface{}) *redis.IntCmd {
	lists := storage.RedisConn().LPush(listName, values...)
	fmt.Printf("redis_push_list_%+v", lists)
	return lists
}

//读取 list
func GetRedisList(cc *lib.Cusctx, listName string, start, num_items int64) (*redis.StringSliceCmd, []*Command) {
	lists := storage.RedisConn().LRange(listName, start, start+num_items-1)
	if len(lists.Val()) < int(num_items) {
		ss := strings.Split(listName, "_")
		userID := ss[len(ss)-1]
		var listMongo []*Command
		listMongo = FindCommands(cc, userID, 10)
		cc.Logf("commands_from_mongo_")
		return nil, listMongo
	}
	cc.Logf("comamnd_from_redis_")
	return lists, nil
}

//查看一个评论
func FindACommandById(cc *lib.Cusctx, id string) (*Command, error) {
	var command *Command
	findC := bson.M{
		"_id": id,
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_COMMAND).Find(findC).One(&command)
	if err != nil {
		return nil, err
	}
	return command, nil
}

// Inc 的应用:
// 计算次数
// 限制 api 的访问次数
//查看一个评论的点击次数
// key : command_id
func FindACommandClickCount(cc *lib.Cusctx, id string) (string, error) {
	// first from redis
	stringcmd := storage.RedisConn().Get(COMAMND_PREFIX + id)
	if stringcmd.Val() == "" {
		// from mongo
		count, findErr := FindAClick(cc, id)
		if findErr != nil {
			return "", findErr
		}
		return count.Num, nil
	}
	return stringcmd.Val(), nil
}

// 添加一个点击的次数
func AddCommandClickCount(cc *lib.Cusctx, id string) (int64, error) {
	intcmd := storage.RedisConn().Incr(COMAMND_PREFIX + id)
	if intcmd.Err() != nil {
		// read from mongo
		clicknum, incErr := IncClick(cc, COMAMND_PREFIX+id)
		if incErr != nil {
			return 0, incErr
		}
		return clicknum, nil
	}
	return intcmd.Val(), intcmd.Err()
}

// CAS 机制实现 乐观锁
// redis 的 WATCH 命令是用 CAS 来实现 事物的
// WATCH 命令：watched key 被监听。如果一个 watched key 在EXEC 之前被修改，
// 那么整个 transaction 就会被 aborted 掉。返回一个 null reply 去通知失败
// 比如： val = GET mykey
// val = val + 1
// SET mykey $val
// 当我们只有一个客户端的时候，这样修改是没有问题的。但是多个客户端同时修改key的时候，会出现竞争条件。比如 客户端 A 和客户端 B 将会读取旧的 值，比如10
// 两边的客户端值将会变成11而不是12

//如果使用 WATCH 的话，就会变成这样

// WATCH mykey
// val = GET mykey
// val = val + 1
// MULTI
// SET mykey $val
// EXEC

// 上面的代码，如果出现了竞争条件，另一个客户端在 WATCH 命令和 EXEC 命令中修改了 val 的值，那么 transaction 将会失败。
// 失败了之后会重试，直到没有出现竞争为止。

// WACTH命令。当没有 WATCHED key 被修改的时候，redis 开始执行transaction。不然transaction不会进来。
// (WATCH 一个无常的key的时候，redis expired了它，EXEC还将会继续执行。)

// 当 EXEC 执行的时候，所有的 key 将会被 UNWATCHED。不管 transaction有没有被 aborted掉。需要 flushed 所有的key的时候，可以使用 UNWATCH命令。
// 使用 WATCH 来实现 ZPOP
// WATCH zset
// elemnt = ZRANGE zset 0 0
// MULTI
// ZREM zset element
// EXEC

// 使用 watch 和 expired 来实现一个限速器
// 限速器需要一个过期的时间，过期时间的写入和计数同样重要
// 使用 transaction同步写入数据
// go-redis 的 TxPipeline 方法
func SpeedLimiter(key string, expiredTime int64) (bool, error) {
	// pipe := storage.RedisConn().TxPipeline()
	// fmt.Printf("pipe_%+v \n", pipe)
	// boolcmd := pipe.Expire(key, time.Duration(expiredTime)*time.Second)
	// fmt.Printf("过期时间_%+v", time.Duration(expiredTime)*time.Second)
	// fmt.Printf("boolcmd_%+v \n", boolcmd)
	// return boolcmd.Val(), boolcmd.Err()
	// limiter := redis.NewScript(`FUNCTION LIMIT_API_CALL(KEYS[0])
	// ts = CURRENT_UNIX_TIME()
	// keyname = KEYS[0] + ":" + ts
	// current = GET(keyname)
	// IF current != NULL AND current > 10 THEN
	// 	ERROR "too many requests per second"
	// ELSE MULTI
	// 		INCR(keyname, 1)
	// 		EXPIRE(keyname, 10)
	// 	 EXEC
	// 	 PERFORM_API_CALL()
	// END
	// `)
	// n, err := limiter.Run(storage.RedisConn(), []string{"ip"}, "10.0.0.10").Result()
	// if err != nil {
	// 	fmt.Printf("err_%+v", err)
	// }
	// fmt.Printf("n_%+v", n)
	// return false, nil
	Publish("redisClient2", "from redisClient2 msg.")
	Sub("redisClient2")
	return false, nil
}

func Publish(client, msg string) {
	incmd := storage.RedisConn().Publish(client, msg)
	fmt.Printf("*********************** incmd_%+v", incmd)
}

func Sub(client string) {
	pubS := storage.RedisConn().Subscribe(client)
	msg, err := pubS.ReceiveMessage()
	fmt.Printf("pubS_message_%+v", msg)
	fmt.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^ err_%+v", err)
}
