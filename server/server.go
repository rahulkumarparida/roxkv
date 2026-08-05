package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)

const MaxConnections = 10

var mutex = sync.RWMutex{}

func DeleteClientListing(target *utils.NewClient) {
	mutex.Lock()
	defer mutex.Unlock()
	utils.TotalConnecntions = slices.DeleteFunc(utils.TotalConnecntions, func(n *utils.NewClient) bool {

		return n.ID == target.ID

	})

}

func handleConnection(client *utils.NewClient, store *store.MemoryAlloc) {
	reader := bufio.NewReader(client.Conn)
	ctx, cancel := context.WithCancel(context.Background())

	msg := fmt.Sprintln("Total Users Connected: ", len(utils.TotalConnecntions))
	client.Conn.Write([]byte(msg))

	defer cancel()
	defer client.Conn.Close()
	defer DeleteClientListing(client)

	go DeadOrAliveConnections(ctx, client)
	go worker.ExpiryWorker(ctx, store)

	for {

		input, err := reader.ReadString('\n')
		print(input)

		mutex.Lock()
		utils.TotalInputs += 1
		mutex.Unlock()

		client.Mu.Lock()
		client.Interactions += 1
		client.LastUsed = time.Now()
		client.Mu.Unlock()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client Disconnected:", err)
				fmt.Println("Client name: ", client.ID)
			} else {
				fmt.Println("err:", err)
				logger.ErrorLog("Error while reading data")
			}
			DeleteClientListing(client)
			break
		}

		// Gets the data from type interface{}/any to string and then writes to byte
		data := commands.ParseCommands(store, strings.Fields(input), client)
		datastr := fmt.Sprintf("%v", data)

		_, werr := client.Conn.Write([]byte("roxkv> " + datastr + "\n"))
		// _, werr := client.Conn.Write([]byte("+" + string(input) + "\r\n"))

		if werr != nil {
			fmt.Println("err:", werr)

		}

	}

}

func Server() {
	config, err := abstractor.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load AI configuration: ", err)
	}
	provider, err := abstractor.NewProvider(*config)
	if err != nil {
		log.Fatal("Failed to initialize AI provider: ", err)
	}

	stre, _ := store.StoreInMemory()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listner, err := net.Listen("tcp", ":6969")

	if err != nil {
		fmt.Println("Server is busy and not listening at port 6969:", err)
		logger.ErrorLog("error while connecting to the TCP server, port 6969 is busy")
		return
	}
	defer listner.Close()
	utils.ServerStarted = time.Now()
	store.StoreHelper = stre
	go worker.SnapshotWorker(ctx, stre, &mutex, utils.TotalConnecntions)
	fmt.Println("Listening CLI Connection at localhost:6969")
	go ChatServer(stre, provider)
	go WebServer()
	go WebChatServer(stre, provider)
	go RespServer(stre, provider)

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
			client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
			client.Conn.Close()
			continue
		}
		utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
		mutex.Unlock()
		fmt.Println("Connected: ", client.ID)
		go handleConnection(&client, stre)

	}

}

// Background worker to remove inactive clients

func ClearConnections(t time.Time, client *utils.NewClient) {
	if time.Since(client.LastUsed) > (10 * time.Minute) {
		fmt.Println("Client died: ", client.ID)
		client.Conn.Write([]byte("Client was Inactive for too long \n"))
		client.Conn.Close()
		DeleteClientListing(client)
		logger.InfoLog("Ticker worked at : " + t.Format("2006-01-02 15:04:05"))

	}
}

func DeadOrAliveConnections(ctx context.Context, client *utils.NewClient) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			ClearConnections(t, client)

		}
	}
}
