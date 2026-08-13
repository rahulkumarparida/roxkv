package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	redisparser "github.com/rahulkumarparida/roxkv/internal/resp-parser"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func RespServer(stre *store.MemoryAlloc, provider abstractor.Provider) error {

	listner, err := net.Listen("tcp", RespTCPAddr)

	if err != nil {
		return fmt.Errorf("start RESP server on %s: %w", RespTCPAddr, err)
	}
	fmt.Println("Listening Resp-Server at localhost" + RespTCPAddr)
	startupLog("server ports", "RESP TCP "+RespTCPAddr)

	go func() {
		defer listner.Close()

		for {
			conn, err := listner.Accept()

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				logger.ErrorLog("RESP TCP connection accept error: " + err.Error())
				continue
			}

			client := *utils.CreateClient(conn, "client")

			mutex.Lock()
			client.Initiator = utils.RedisSource
			if len(utils.TotalConnecntions) > MaxConnections {
				data := redisparser.EncodeBulkBytes([]byte("Max connections from the TCP server exceeded"))
				client.Conn.Write(data)
				client.Conn.Close()
				mutex.Unlock()
				continue
			}
			utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
			mutex.Unlock()
			logger.InfoLog("New RESP client connected: " + fmt.Sprintf("%v", client.ID))
			go HandleRespConnections(&client, stre, provider)
		}
	}()

	return nil
}

func HandleRespConnections(client *utils.NewClient, stre *store.MemoryAlloc, provider abstractor.Provider) {

	logger.InfoLog("Resp Client connected: " + client.ID.(string))

	defer client.Conn.Close()

	redisparser.ReadAndHandleConnection(client)

}
