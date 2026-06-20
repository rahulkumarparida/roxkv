package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)



func handleConnection(conn net.Conn, store *store.MemoryAlloc) {
	reader := bufio.NewReader(conn)


	for {
		worker.ExpiryWorker(store)
		input, err := reader.ReadString('\n')

		if input == "q" || input == "exit" {

			break
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client Disconnected")
			} else {
				fmt.Println("err:", err)
				logger.ErrorLog("Error while reading data")
			}
			return
		}
		// Gets the data from type interface{}/any to string and then writes to byte
		data := commands.ParseCommands(store,strings.Fields(input))
		datastr := fmt.Sprintf("%v", data)
		_, werr := conn.Write([]byte("roxkv> " + datastr + " \n" ))

		if werr != nil {
			fmt.Println("err:", werr)

		}

	}
	defer conn.Close()

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

		go handleConnection(conn , store)

	}

}
