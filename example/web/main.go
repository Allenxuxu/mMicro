package main

import (
	micro "github.com/Allenxuxu/mMicro"
	"github.com/Allenxuxu/mMicro/transport/grpc"
	"github.com/Allenxuxu/mMicro/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.api.agent"),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	service.Options().Service.Init(micro.Transport(grpc.NewTransport()))

	router := gin.Default()
	router.GET("/agent/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello mMicro")
	})
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}
