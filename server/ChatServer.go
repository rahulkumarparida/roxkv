package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/ollama/ollama/api"
	master "github.com/rahulkumarparida/roxkv/agents/Master"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var cmutex = sync.Mutex{}

func handleChatConnection(user *utils.NewClient, stre *store.MemoryAlloc, namespace *store.NameSpace, agent *api.Client) {

	reader := bufio.NewReader(user.Conn)

	msg := fmt.Sprintln("Total Users Connected: ", len(utils.TotalConnecntions))
	user.Conn.Write([]byte(msg))

	defer user.Conn.Close()

	for {

		input, err := reader.ReadString('\n')
		print(input)

		cmutex.Lock()
		user.Interactions += 1
		utils.TotalInputs += 1
		cmutex.Unlock()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client Disconnected:", err)
				fmt.Println("Client name: ", user.ID)
			} else {
				fmt.Println("err:", err)
				logger.ErrorLog("Error while reading data")
			}
			break
		}

		// Route all requests through the master orchestrator so specialist-agent results are collected and summarized centrally.
		data := master.MasterAgent(input, stre, user, namespace, agent)

		dataString := fmt.Sprintf("%v",data)
		user.Conn.Write([]byte("\nroxai> "+dataString+"\n"))
	}

}

func ChatServer(stre *store.MemoryAlloc, namespace *store.NameSpace, agent *api.Client) {
	
	listner, err := net.Listen("tcp", ":6970")

	if err != nil {
		fmt.Println("Server is busy and not listening at port 6970:", err)
		logger.ErrorLog("error while connecting to the TCP server, port 6969 is busy")
		return
	}
	defer listner.Close()
	fmt.Println("Listening RoxAI Connections at localhost:6970")

	for {
		conn, err := listner.Accept()

		if err != nil {
			fmt.Println("Connection could not be established:", err)
			logger.ErrorLog("Connection failed could not be established")
			continue
		}

		client := *utils.CreateClient(conn, "system")

		cmutex.Lock()
		if len(utils.TotalConnecntions) > MaxConnections {
			client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
			client.Conn.Close()
			continue
		}
		utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
		cmutex.Unlock()

		fmt.Println("Connected: ", client.ID)
		go handleChatConnection(&client, stre, namespace, agent)

	}

}
