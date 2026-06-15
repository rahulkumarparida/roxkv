package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func CreateFile(data []store.Item, absoluteFilePath string)  bool{
	

	datafile , err:= os.Create(absoluteFilePath)
	if utils.HandleError("Error while creating file: ", err) {
		return  false
	}
	
	defer datafile.Close()	
	
	databytes, err := json.MarshalIndent(data, "", "  ")
	if utils.HandleError("Error while Marshaling file: ", err){
		return false
	}
	
	_ , err = datafile.Write(databytes)
	if utils.HandleError("Error while writing file: ", err){
		return  false
	}

	return true
	
}

//Create a func previous data where i can add previous data form the json so that it can be added accorindly
// such that it should seregate by the time like difffrent time diffrent positoin of the Item in the araray
// the oldest will be added first and then later will bw addded later so we cna get the latest values

//Func to retireve the current dates data
func retreiveData(fileName string , data []store.Item) []store.Item{

	_, err := os.Stat(fileName)

	if err == nil {
		databytes , err := os.ReadFile(fileName)
		var olddata []store.Item
		if utils.HandleError("Error while creating the folder", err) {
			logger.ErrorLog("Error while retrieving old data")
			return  data
		}

		errr := json.Unmarshal(databytes, &olddata)
		if utils.HandleError("Error while creating the folder", errr) {
			logger.ErrorLog("Error while retrieving old data")
			return  data
		}


		func ()  {
			
			seen := make(map[string]bool)
	
			for _, val := range data {
				seen[val.Key] = true
			}
			
			for _, val := range olddata {
				if !seen[val.Key] {
					data = append(data, val)
					seen[val.Key] = true 
				}
			}


		}()


		return  data
	}
	if errors.Is(err , os.ErrNotExist) {
		return data
	}

	return  data

}

func StoreToJson(data []store.Item) bool{
	dbFolder := utils.DbFolder()
	if len(dbFolder) == 0 {	
		return false
	}
	err := os.MkdirAll(dbFolder, 0755)

	if err != nil {
		utils.HandleError("Error while creating the folder", err)
		return  false
	}

	filename := time.Now().Format("2006-01-02")+".json"
	KeyfolderName := filepath.Join(dbFolder,filename)

	allData := retreiveData(KeyfolderName,data)
			
	val := CreateFile(allData,KeyfolderName)

	if !val {
		fmt.Println("File not created")
		return	false
	}
	
	return val
	
}


// oldest to newest
func sortFile(allfiles []os.DirEntry) []os.DirEntry{
	
	 sort.Slice(allfiles, func(i, j int) bool {
		
		infoI, err := allfiles[i].Info()
		
		if utils.HandleError("error reading file info",err) {
			return false
		}

		infoJ , err := allfiles[j].Info()
		
		if utils.HandleError("error reading file info",err) {
			return false
		}
		
		return infoI.ModTime().Before(infoJ.ModTime())

	})

	return  allfiles

}



func LoadJsons(ms *store.MemoryAlloc) int{
	dbFolder := utils.DbFolder()
	count := 1
	if len(dbFolder) == 0 {	
		fmt.Println("Home directory not found")
		return count
	}

	// Retieves all file form the path stores in an slice
	allFiles , err := os.ReadDir(dbFolder)
	
	if utils.HandleError("Error while reading the files", err) {
		return count
	}

	var dataArr []store.Item  


	var sortedFiles = sortFile(allFiles)
	
	// Loops through all the files to get the data 
	for _, file := range sortedFiles {
		if file.IsDir() {
			continue
		}

		fmt.Println("Executing: ",file.Name())
		filePath  := filepath.Join(dbFolder ,file.Name())

		databytes , err := os.ReadFile(filePath)
		if utils.HandleError("Error while reading the folder", err) {
			
			continue
		}
		var data []store.Item
		err = json.Unmarshal(databytes , &data)
		dataArr = append(dataArr, data...)

		count++
	}


	for _, val := range dataArr {
		store.SetKv(ms , &val , []string{"",""})
		fmt.Println("Val Set: ", val.Key)
	}

	return count
}