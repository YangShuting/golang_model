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
