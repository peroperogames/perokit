package ginm

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/peroperogames/perokit/core/hardware_percent"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// InitWebSocketConnMiddle Websocket连接初始化中间件
// memPercentLimit传入内存占用上限值，如果当前服务器的内存占用大于该值，则返回错误并断开ws连接
// Websocket.Conn连接对象存储于gin.Context中，key为'ContextWebsocketConn'，使用Get方法取出
func InitWebSocketConnMiddle(memPercentLimit float64) gin.HandlerFunc {
	type ApiJson struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	type ErrJson struct {
		Err string `json:"err"`
	}

	type WebsocketResp struct {
		MessageType int32  `json:"message_type"`
		Code        int32  `json:"code"`
		Data        string `json:"data"`
	}

	return func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, ApiJson{
				Code: -1,
				Msg:  "err",
				Data: ErrJson{Err: err.Error()},
			})
			c.Abort()
			return
		}

		//连接创建时检查内存占用率，如果大于设定值则断开连接
		memPercent := hardware_percent.GetMemPercent()
		if memPercent >= memPercentLimit {
			resp := WebsocketResp{
				MessageType: 8,
				Code:        -1,
				Data:        "memory limit",
			}
			data, _ := json.Marshal(&resp)
			_ = ws.WriteMessage(websocket.TextMessage, data)
			c.Abort()
			_ = ws.Close()
			return
		}

		c.Set("ContextWebsocketConn", ws)
		c.Next()
	}
}
