package redisparser

import (
	"fmt"
	"strconv"
)



func EncodeSimpleError(stmt string) (string , error) {
	
	encode := fmt.Sprintf("-%s\r\n", stmt)

	return encode , nil
}

func EncodeNullValues() (string){

	return "$-1\r\n"

}

func EncodeBoolean(data bool) (string,error){
	var encode string
	if data {
		encode = "#t\r\n"
	}else{
		encode = "#f\r\n"
	}
	return encode , nil
}


func EncodeSimpleString(data string) (string , error){

	encode := fmt.Sprintf("+%s\r\n",data)

	return  encode , nil
}


func EncodeInteger(data int64) (string,error){
	
	strvalue := strconv.FormatInt(data,10)
	
	encode := fmt.Sprintf(":%s\r\n",strvalue)

	return encode , nil
}


func EncodeBulkString(data string) (string,error){

	lengthOfData := len(data)
	encode := fmt.Sprintf("$%v\r\n%s\r\n",lengthOfData,data)

	return  encode, nil
}


// *<no-of-elemnets>\r\n$4\r\nrahu\r\n
func EncodeArray[T any](data []T) (string , error){
   
	length := len(data)
	var NewData string 
	fmt.Println("Writinf to user:", length)
	for idx := range length {
		b := data[idx]
		var encode string	
		switch any(b).(type) {
			case nil:
				encode = EncodeNullValues()
				
			case string:
				encode , _ = EncodeBulkString(fmt.Sprintf("%v",b)) 			
				
			case []string:
				var err error
				encode , err = EncodeStringArray(b) 	
				if err != nil {
					encode , _ = EncodeSimpleError("ERR parsing the argument")
				}

			case int,int64:
				fmt.Println("Converstion int ")
				val , err := strconv.ParseInt(fmt.Sprintf("%v",b),10,64)
				if err != nil {
					encode , _ = EncodeSimpleError("ERR parsing the argument")
				}else{
					encode , _ = EncodeInteger(val)		
				}
				
			default:
				encode , _ = EncodeBulkString(fmt.Sprintf("%v",b)) 			
			}
			
		NewData += encode 	
		fmt.Println("New Data:",NewData)
		encode = ""
	}

	encode := fmt.Sprintf("*%v\r\n%s",length,NewData)
	

	return encode , nil
}

func EncodeStringArray(data any)  (string,error){

	
	var encoded string
	for _, v := range data.([]string) {

		encode, err := EncodeBulkString(v)
		if err != nil{
			encode , err = EncodeSimpleError("ERR parsinr string array")
		}
		encoded += encode

	}

	return fmt.Sprintf("*%v\r\n%s",len(data.([]string)),encoded),nil
	
}