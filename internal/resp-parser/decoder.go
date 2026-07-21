package redisparser

import (
	// "bufio"
	"errors"
	"fmt"
	// "io"
	"strconv"
	// "strings"
)




const (
	BOOLEAN = '#'
	DOUBLES = ','
	SIMPLESTRING = '+'
	BULKSTRING = '$'
	INTERGER = ':'
	ARRAY = '*'
	ERROR = '-'
)

// func main(){
// 	rawStream := io.Reader(strings.NewReader("-Error Reading a data\r\n\t"))
// 	reader := bufio.NewReader(rawStream)
// 	line, err := reader.ReadBytes('\t')
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}


// 	val , err := Decode(line)
// 	fmt.Println("Val:", val)

// }


func simpleReadError(data []byte) (string,int,error){
	pos := 1

	for; data[pos] != '\r'; pos++{
	}

	errorval := string(data[1:pos])

	return errorval , pos+2 , nil
}

// #<t|f>\r\n
func readBoolean(data []byte) (bool,int,error){
	pos :=1
	var value bool

	if string(data[pos]) == "t"  {
		value = true
	}else if string(data[pos]) == "f"{
		value = false
	}else{
		value = false
		return value , pos , errors.New("Courrpted Value")
	}


	return value , pos+2 , nil

}


func readSimpleString(data []byte)  (string,int,error){
	pos := 1

	for; data[pos] != '\r'; pos++{
	}

	return string(data[1:pos]) , pos+2 , nil
}


func readInteger(data []byte)  (int64,int,error){
	pos := 1

	for; data[pos] != byte('\r'); pos++{
	}

	i64 , err := strconv.ParseInt(string(data[1:pos]), 10 , 64)

	if err != nil {
		return  0 , pos , err
	}


	

	return i64 , pos+2 , nil
}

func readLength(data []byte) (int,int){
	pos , length := 0,0

	for pos = range data {
		b := data[pos]

		if !(b >= '0' && b <= '9') {
			return  length , pos +2
		}
		length = length*10 + int(b - '0')
	}

	return  0,0

}

func readBulkString(data []byte)   (string,int,error){
	
	pos := 1

	
	// returns the length of the string incoming and the delta value which is the number of charachters read after the first char
	length , delta := readLength(data[pos:])
	pos = pos + delta
			// starts slicing aftre the \r\n and ends on the exact length provided
	return  string(data[pos:(pos+length)]) , pos+ length + 2 , nil
}


func readArray(data []byte)  (any,int,error){


	pos := 1

	totalElems , delta := readLength(data[pos:])

	pos += delta

	var elements []any = make([]any, totalElems);

	for i := range totalElems{

		val , posi , err := DecodeOne(data[pos:])
	
		if err != nil {
			return nil , 0 , errors.New("Error while parsing elements")
		}

		elements[i] = val
		pos += posi
	}

	return  elements , pos+2 , nil
}


//  ,[<+|->]<integral>[.<fractional>][<E|e>[sign]<exponent>]\r\n
// ,1.23\r\n  , ,10\r\n

func readDoubles(data []byte) (float64 , int , error){
	pos := 1
	var integerpart int = 0
	var decimalpart int = 0
	var end bool = false


	for; data[pos] != '.'; pos++{
		b := data[pos]
		if b == '\r' {
			end = true
			break
		}
		integerpart = integerpart * 10 + int(b -'0')
	}

	if end {
		return float64(integerpart) , pos+2 , nil
	}else{
		pos +=1
		for; data[pos] != '\r'; pos++{
			b := data[pos]
			decimalpart = decimalpart*10 + int(b - '0')
		}
	}

	f64 , err := strconv.ParseFloat(fmt.Sprintf("%v.%v", integerpart, decimalpart),64)

	if err != nil{
		return 0.0 , pos , errors.New("Error Parsing the decimal")
	}


	return f64 , pos+2 , nil
}


func DecodeOne(data []byte) (any,int,error){

	if len(data) == 0 {
		return  nil, 0 , errors.New("No Data Found")
	}

	switch data[0]{
		case SIMPLESTRING:
			fmt.Println("Simple string")
			return readSimpleString(data)
		case INTERGER:
			fmt.Println("Integer")
			return readInteger(data)
		case BULKSTRING:
			fmt.Println("Bulk String")
			return readBulkString(data)
		case ARRAY:
			fmt.Println("Array")
			return readArray(data)
		case BOOLEAN:
			fmt.Println("Boolean")
			return readBoolean(data)
		case DOUBLES:
			fmt.Println("Double")
			return readDoubles(data)
		case ERROR:
			fmt.Println("Error")
			return simpleReadError(data)
		default:
			return 	nil , 0 , errors.New("Data Type Didn't match")


	}
	


}


func Decode(data []byte) (any,error){
	if len(data) == 0 {
		return  nil , errors.New("No Data Found")
	}

	value , _ , err := DecodeOne(data)

	return  value , err
}


func DecodeArrayString(data []byte) ([]string, error){
	value , err := Decode(data)
	if err != nil {
		return  nil , err
	}


	array := value.([]any)
	tokens := make([]string, len(array))

	for i := range tokens {
		tokens[i] = array[i].(string)
	}

	return tokens , nil
}