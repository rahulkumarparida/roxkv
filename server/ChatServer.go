package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/ollama/ollama/api"
	storageagent "github.com/rahulkumarparida/roxkv/agents/StorageAgent"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)


func handleChatConnection(user *utils.NewClient,stre *store.MemoryAlloc,namespace *store.NameSpace,agent *api.Client){

	reader := bufio.NewReader(user.Conn)

	msg := fmt.Sprintln("Total Users Connected: ", len(utils.TotalConnecntions))
	user.Conn.Write([]byte(msg))
	mutex := sync.Mutex{}

	for {
		
		input, err := reader.ReadString('\n')
		print(input)

		mutex.Lock()
		user.Interactions += 1
		utils.TotalInputs += 1
		mutex.Unlock()
		


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

		// Gets the data from type interface{}/any to string and then writes to byte
		// kvagent.KvAgent(input,stre,user,namespace,agent)
		// monitoragent.MonitorAgent(input,stre,user,agent)
		storageagent.StorageAgent(input,stre,namespace,user,agent)
		
		

	}

}


func ChatServer(stre *store.MemoryAlloc,namespace *store.NameSpace){
	agent,err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	listner, err := net.Listen("tcp", ":6970")

	if err != nil {
		fmt.Println("Server is busy and not listening at port 6970:", err)
		logger.ErrorLog("error while connecting to the TCP server, port 6969 is busy")
		return
	}
	defer listner.Close()
	fmt.Println("Listening at localhost:6970")

	for {
		conn, err := listner.Accept()

		if err != nil {
			fmt.Println("Connection could not be established:", err)
			logger.ErrorLog("Connection failed could not be established")
			continue
		}

		client := *utils.CreateClient(conn)

		if len(utils.TotalConnecntions) > MaxConnections {
			client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
			client.Conn.Close()
			continue
		}

		utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
		fmt.Println("Connected: ", client.ID)
		go handleChatConnection(&client, stre,namespace, agent)

	}

}