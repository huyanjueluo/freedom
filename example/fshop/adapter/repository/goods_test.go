package repository

import (
	"fmt"
	"testing"

	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/example/fshop/domain/entity"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getUnitTest() freedom.UnitTest {
	//创建单元测试工具
	unitTest := freedom.NewUnitTest()
	unitTest.InstallDB(func() interface{} {
		db, e := gorm.Open("mysql", "root:123123@tcp(127.0.0.1:3306)/fshop?charset=utf8&parseTime=True&loc=Local")
		if e != nil {
			freedom.Logger().Fatal(e.Error())
		}
		db = db.Debug()
		return db
	})

	opt := &redis.Options{
		Addr: "127.0.0.1:6379",
	}
	redisClient := redis.NewClient(opt)
	if e := redisClient.Ping().Err(); e != nil {
		freedom.Logger().Fatal(e.Error())
	}
	unitTest.InstallRedis(func() (client redis.Cmdable) {
		return redisClient
	})

	mockEvent := unitTest.NewDomainEventInfra(func(producer, topic string, data []byte, header map[string]string) {
		//mock 一个领域事件基础设施
		fmt.Println("mock-event-send", producer, topic, string(data), header)
	})
	unitTest.InstallDomainEventInfra(mockEvent)
	return unitTest
}

// TestGoodsEntity 商品实体单测
func TestGoodsEntity(t *testing.T) {
	//获取单测工具
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *Goods
	//获取资源库
	unitTest.GetRepository(&repo)
	goodsEnity, err := repo.Get(1)
	if err != nil {
		t.Error(err)
		return
	}
	err = goodsEnity.MarkedTag(entity.GoodsNewTag)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.Save(goodsEnity)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok", goodsEnity)
}
