package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	master "github.com/rahulkumarparida/roxkv/agents/Master"
	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var cmutex = sync.Mutex{}

func handleChatConnection(user *utils.NewClient, stre *store.MemoryAlloc, provider abstractor.Provider) {

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
		data := master.MasterAgent(input, stre, user, provider)

		dataString := fmt.Sprintf("%v", data)
		user.Conn.Write([]byte("\nroxai> " + dataString + "\n"))
	}

}

func ChatServer(stre *store.MemoryAlloc, provider abstractor.Provider) error {
	listner, err := net.Listen("tcp", ChatTCPAddr)

	if err != nil {
		return fmt.Errorf("start RoxAI TCP server on %s: %w", ChatTCPAddr, err)
	}
	fmt.Println("Listening RoxAI Connections at localhost" + ChatTCPAddr)
	startupLog("server ports", "AI TCP "+ChatTCPAddr)

	go func() {
		defer listner.Close()

		for {
			conn, err := listner.Accept()

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				fmt.Println("Connection could not be established:", err)
				logger.ErrorLog("Connection failed could not be established")
				continue
			}

			client := *utils.CreateClient(conn, "system")

			cmutex.Lock()
			if len(utils.TotalConnecntions) > MaxConnections {
				client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
				client.Conn.Close()
				cmutex.Unlock()
				continue
			}
			utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
			cmutex.Unlock()

			fmt.Println("Connected: ", client.ID)
			go handleChatConnection(&client, stre, provider)
		}
	}()

	return nil
}
