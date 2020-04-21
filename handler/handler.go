package handler

import (
	"chat/dao"
	"chat/model"
	"chat/util"
	"crypto/md5"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type ConnectedClient struct {
	Lock       sync.Mutex
	Connection map[string]*websocket.Conn
}

var connection =ConnectedClient{
	Lock:       sync.Mutex{},
	Connection: make(map[string]*websocket.Conn, 2),
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func chat(ctx *gin.Context) {
	if len(connection.Connection) >= 2 {
		util.FailedResponse(ctx, util.PersonExceed, util.PersonExceedMsg)
		return
	}
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("upgrade:%s", err)
		return
	}
	defer c.Close()
	name, _ := ctx.Get("name")
	connection.Lock.Lock()
	connection.Connection[name.(string)] = c
	connection.Lock.Unlock()
	log.Infof("%s success to join in room",name)
	err, list := dao.GetNotReadMessage(name.(string))
	if err != nil {
		log.Errorf("error:%s", err.Error())
		util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
		return
	}
	if len(list) > 0 {
		data, err := json.Marshal(list)
		if err != nil {
			log.Errorf("error:%s", err.Error())
			util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
			return
		}
		err = c.WriteMessage(1, data)
		if err != nil {
			log.Errorf("error:%s", err.Error())
			util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
			return
		}
	}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		recs := model.Message{}
		err = json.Unmarshal(message, &recs)
		if err != nil {
			log.Errorf("error:%s", err.Error())
			util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
			return
		}
		if len(connection.Connection) == 1 {
			recs.IsRead = 0
		} else {
			for who, conn := range connection.Connection {
				if who != name.(string) {
					data, err := json.Marshal([]model.Message{recs})
					if err != nil {
						recs.IsRead = 0
						log.Errorf("error:%s", err.Error())
						util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
						return
					}
					err = conn.WriteMessage(mt, data)
					if err != nil {
						recs.IsRead = 0
						log.Errorf("error:%s", err.Error())
						util.FailedResponse(ctx, util.InternalError, util.InternalErrorMsg)
						return
					}
					recs.IsRead = 1
				}
			}
		}
		if err := dao.InsertMessage(recs.From, recs.To, recs.Content, recs.IsRead); err != nil {
			log.Errorf("error:%s", err.Error())
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	connection.Lock.Lock()
	delete(connection.Connection, name.(string))
	connection.Lock.Unlock()
}

type LoginReq struct {
	Username string `json:"username"`
	Ticket   string `json:"ticket"`
}

func login(ctx *gin.Context) {
	req := LoginReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("login para error:%s", err.Error())
		util.FailedResponse(ctx, util.LoginFailed, util.LoginFailedMsg)
		return
	}
	mymd5 := md5.New()
	ticket:=req.Username + time.Now().Format("20060102")
	mymd5.Write([]byte(req.Username + time.Now().Format("20060102")))
	//ticket := hex.EncodeToString(mymd5.Sum(nil))
	if ticket == req.Ticket {
		log.Infof("%s login successfully", req.Username)
		ctx.SetCookie("username", req.Username, 7*24*3600, "/", "", false, false)
		util.OKResponse(ctx, nil)
	} else {
		log.Infof("%s failed to login", req.Username)
		util.FailedResponse(ctx,util.LoginFailed,util.LoginFailedMsg)
	}
	return
}
