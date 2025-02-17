package taskgin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	goservice "github.com/phathdt/libs/go-sdk"
	"github.com/phathdt/libs/go-sdk/sdkcm"
	"gorm.io/gorm"
	"togo/common"
	"togo/modules/task/taskhandler"
	"togo/modules/task/taskrepo"
	"togo/modules/task/taskstorage"
)

func DeleteTask(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("current_user").(*sdkcm.SimpleUser)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		db := sc.MustGet(common.DBMain).(*gorm.DB)
		redisClient := sc.MustGet(common.PluginRedis).(*redis.Client)
		store := taskstorage.NewSQLStore(db)
		redisStore := taskstorage.NewRedisStore(redisClient)
		repo := taskrepo.NewRepo(store, redisStore)

		hdl := taskhandler.NewDeleteTaskHdl(repo, user)
		if err = hdl.Response(c.Request.Context(), id); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, sdkcm.SimpleSuccessResponse("ok"))
	}
}
