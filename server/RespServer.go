package server

import (
	"fmt"
	"net"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	redisparser "github.com/rahulkumarparida/roxkv/internal/resp-parser"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func RespServer(stre *store.MemoryAlloc, agent *api.Client) {

	listner, err := net.Listen("tcp", ":6973")

	if err != nil {
		fmt.Println("Server is busy and not listening at port 6973:\n", err)
		logger.ErrorLog("error while connecting to the TCP server, port 6973 is busy")
		return
	}
	defer listner.Close()

	fmt.Println("Listening Resp-Server at localhost:6973")

	for {

		conn, err := listner.Accept()

		if err != nil {
			fmt.Println("Connection could not be established:", err)
			logger.ErrorLog("Connection failed could not be established")
			continue
		}

		client := *utils.CreateClient(conn, "client")

		mutex.Lock()
		if len(utils.TotalConnecntions) > MaxConnections {
			data := redisparser.EncodeBulkBytes([]byte("Max connections from the TCP server exceeded"))
			client.Conn.Write(data)
			client.Conn.Close()
			continue
		}
		utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
		mutex.Unlock()
		fmt.Println("Connected: ", client.ID)
		go HandleRespConnections(&client, stre, agent)

	}

}

func HandleRespConnections(client *utils.NewClient, stre *store.MemoryAlloc, agent *api.Client) {

	logger.InfoLog("Resp Client connected: " + client.ID.(string))


	defer client.Conn.Close()


	redisparser.ReadAndHandleConnection(client)
	
}
