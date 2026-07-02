package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func CreateFile(data any, absoluteFilePath string)  bool{
	

	datafile , err:= os.Create(absoluteFilePath)
	if utils.HandleError("Error while creating file: ", err) {
		return  false
	}
	
	defer datafile.Close()	
	fmt.Println("DataCreations: ", data)
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

//Func to retireve the current dates data
func retreiveData(fileName string , data any) any{

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
			
			if data , ok := data.([]store.Item); ok {
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
			}
		}()


		return  data
	}
	if errors.Is(err , os.ErrNotExist) {
		return data
	}

	return  data

}

func StoreToJson(dbFolder string,filename string,data any) bool{
	if len(dbFolder) == 0 {	
		return false
	}
	err := os.MkdirAll(dbFolder, 0755)

	if err != nil {
		utils.HandleError("Error while creating the folder", err)
		return  false
	}

	KeyfolderName := filepath.Join(dbFolder,filename)
	

	var allData any

	if stdata , ok := data.([]store.Item); ok {
		
		allData = retreiveData(KeyfolderName,stdata)	

	}else{
		allData = data
	}
	fmt.Println("DataRecieved: ", allData)
	val := CreateFile(allData,KeyfolderName)		


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


// loads to the memory not needeee when retrieveing snapshots
func LoadJsons(dbFolder string,ms *store.MemoryAlloc,namespace *store.NameSpace) int{
	
	count := 1
	if len(dbFolder) == 0 {	
		fmt.Println("Home directory not found")
		return count
	}

	// Retieves all file form the path stores in an slice
	allFiles , err := os.ReadDir(dbFolder)
	
	if utils.HandleError("Error while reading the files:", err) {
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
		store.SetKv(ms ,namespace, &val)
		fmt.Println("Val Set: ", val.Key)
	}

	return count
}