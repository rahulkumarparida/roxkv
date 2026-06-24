package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)



const MaxConnections = 5
type NewClient struct {
	ID any
	Conn net.Conn
	LastUsed time.Time
	ConnectedAt time.Time
	Mu sync.RWMutex
}

func CreateClient(conn net.Conn) *NewClient{


	return  &NewClient{
		ID:conn.RemoteAddr().String(),
		Conn: conn,
		LastUsed: time.Now(),
		ConnectedAt: time.Now(),
		Mu: sync.RWMutex{},
	}
}

var TotalConnecntions []*NewClient

func DeleteClientListing(target *NewClient){
	TotalConnecntions = slices.DeleteFunc(TotalConnecntions, func(n *NewClient) bool{

		return n.ID == target.ID
		 
	})	

}

func handleConnection(client *NewClient, store *store.MemoryAlloc ) {
	reader := bufio.NewReader(client.Conn)	
	ctx , cancel := context.WithCancel(context.Background())	

	msg := fmt.Sprintln("Total Users Connected: ", len(TotalConnecntions) )	
	client.Conn.Write([]byte(msg))
	
	

	for {
		input, err := reader.ReadString('\n')
		print(input)
		go DeadOrAliveConnections(ctx,client)
		go worker.ExpiryWorker(ctx,store)
		

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client Disconnected:",err)
				fmt.Println("Client name: ", client.ID )
			} else {
				fmt.Println("err:", err)
				logger.ErrorLog("Error while reading data")
			}
			DeleteClientListing(client)
			break
		}

		// Gets the data from type interface{}/any to string and then writes to byte
		data := commands.ParseCommands(store,strings.Fields(input))
		datastr := fmt.Sprintf("%v", data)
		_, werr := client.Conn.Write([]byte("roxkv> " + datastr + " \n" ))
		client.LastUsed = time.Now()

		if werr != nil {
			fmt.Println("err:", werr)

		}
		
	}
	
	defer cancel()

	DeleteClientListing(client)	
	defer client.Conn.Close()
}

func ClearConnections(t time.Time,client *NewClient) {
		if time.Since(client.LastUsed) > (5*time.Minute) {
			fmt.Println("Client died: ", client.ID)
			client.Conn.Write([]byte("Client was Inactive for too long \n" ))
			client.Conn.Close()
			DeleteClientListing(client)	
			logger.InfoLog("Ticker worked at : "+ t.Format("2006-01-02 15:04:05")) 
			
		}
}

func DeadOrAliveConnections(ctx context.Context,client *NewClient){
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	
	
	
		fmt.Println("Ticked")	
		for{
			select{
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				ClearConnections(t,client)
				


			}
		}
	


}

func Server() {
	store := store.StoreInMemory()


	listner, err := net.Listen("tcp", ":6969")

	if err != nil {
		fmt.Println("Server is busy and not listening at port 6969:", err)
		logger.ErrorLog("error while connecting to the TCP server, port 6969 is busy")
		return
	}
	defer listner.Close()
	fmt.Println("Listening at localhost:6969")
	for {

		conn, err := listner.Accept()





		if err != nil {
			fmt.Println("Connection could not be established:", err)
			logger.ErrorLog("Connection failed could not be established")
			continue
		}

		client := *CreateClient(conn)

		if len(TotalConnecntions) > MaxConnections {
			client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
			client.Conn.Close()
			continue
		}


		TotalConnecntions = append(TotalConnecntions, &client)
		fmt.Println("Connected: ", client.ID)
		go handleConnection(&client , store )
		


	}

}
