package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/config"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)

const MaxConnections = 10

const (
	NativeTCPAddr     = ":6969"
	ChatTCPAddr       = ":6970"
	EventsHTTPAddr    = ":6971"
	ChatHTTPAddr      = ":6972"
	RespTCPAddr       = ":6973"
	DashboardHTTPAddr = ":8080"
)

type StartOptions struct {
	Command        string
	ExecutablePath string
}

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
		logger.InfoLog("Received input from client " + fmt.Sprintf("%v", client.ID))

		mutex.Lock()
		utils.TotalInputs += 1
		mutex.Unlock()

		client.Mu.Lock()
		client.Interactions += 1
		client.LastUsed = time.Now()
		client.Mu.Unlock()

		if err != nil {
			if err == io.EOF {
				logger.InfoLog("Client disconnected: " + fmt.Sprintf("%v", client.ID))
			} else {
				logger.ErrorLog("Connection read error for client " + fmt.Sprintf("%v", client.ID) + ": " + err.Error())
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
			logger.ErrorLog("Connection write error for client " + fmt.Sprintf("%v", client.ID) + ": " + werr.Error())
		}

	}

}

func Run(options StartOptions) error {
	startedAt := time.Now()
	commandName := strings.TrimSpace(options.Command)
	if commandName == "" {
		commandName = "roxkv tcp"
	}

	executablePath := strings.TrimSpace(options.ExecutablePath)
	if executablePath == "" {
		if currentExecutable, err := os.Executable(); err == nil {
			executablePath = currentExecutable
		}
	}

	startupLog("command detected", commandName)
	if executablePath != "" {
		startupLog("executable location", executablePath)
	}

	if err := config.EnsureConfigDirectory(); err != nil {
		return fmt.Errorf("initialize configuration directory: %w", err)
	}
	startupLog("configuration loaded", fmt.Sprintf("provider config directory ready at %s", utils.ConfigFolder()))

	providerConfig, err := abstractor.LoadConfig()
	if err != nil {
		return fmt.Errorf("load AI configuration from %s: %w", abstractor.ConfigPath(), err)
	}
	startupLog("configuration loaded", fmt.Sprintf("AI provider=%s model=%s path=%s", providerConfig.Provider, providerConfig.Model, abstractor.ConfigPath()))

	provider, err := abstractor.NewProvider(*providerConfig)
	if err != nil {
		return fmt.Errorf("initialize AI provider %s/%s: %w", providerConfig.Provider, providerConfig.Model, err)
	}
	startupLog("agents initialized", fmt.Sprintf("provider %s model %s", providerConfig.Provider, providerConfig.Model))

	stre, _ := store.StoreInMemory()
	store.StoreHelper = stre
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startupLog("storage initialized", fmt.Sprintf("db=%s snapshots=%s logs=%s", utils.DbFolder(), utils.SnapshotFolder(), utils.LogFolder()))

	listner, err := net.Listen("tcp", NativeTCPAddr)
	if err != nil {
		return fmt.Errorf("start TCP server on %s: %w", NativeTCPAddr, err)
	}
	defer listner.Close()

	utils.ServerStarted = time.Now()
	go worker.SnapshotWorker(ctx, stre, &mutex, utils.TotalConnecntions)
	startupLog("server ports", "native TCP "+NativeTCPAddr)

	if err := ChatServer(stre, provider); err != nil {
		return err
	}
	if err := WebServer(); err != nil {
		return err
	}
	if err := WebChatServer(stre, provider); err != nil {
		return err
	}
	if err := RespServer(stre, provider); err != nil {
		return err
	}
	if err := DashboardServer(); err != nil {
		return err
	}

	startupLog("initialization time", time.Since(startedAt).String())
	return serveCLIConnections(listner, stre)
}

func Server() {
	if err := Run(StartOptions{}); err != nil {
		log.Fatal(err)
	}
}

func serveCLIConnections(listener net.Listener, stre *store.MemoryAlloc) error {
	fmt.Println("Listening CLI Connection at localhost" + NativeTCPAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			logger.ErrorLog("CLI connection accept error: " + err.Error())
			continue
		}

		client := *utils.CreateClient(conn, "client")

		mutex.Lock()
		if len(utils.TotalConnecntions) > MaxConnections {
			client.Conn.Write([]byte("\nMax connections from the TCP server exceeded\n"))
			client.Conn.Close()
			mutex.Unlock()
			continue
		}
		utils.TotalConnecntions = append(utils.TotalConnecntions, &client)
		mutex.Unlock()
		logger.InfoLog("New CLI client connected: " + fmt.Sprintf("%v", client.ID))
		go handleConnection(&client, stre)
	}
}

// Background worker to remove inactive clients

func ClearConnections(t time.Time, client *utils.NewClient) {
	if time.Since(client.LastUsed) > (10 * time.Minute) {
		logger.InfoLog("Inactive client removed: " + fmt.Sprintf("%v", client.ID))
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

func startupLog(step string, message string) {
	formatted := fmt.Sprintf("%s: %s", step, message)
	log.Printf("[startup] %s", formatted)
	logger.InfoLog("[startup] " + formatted)
}
