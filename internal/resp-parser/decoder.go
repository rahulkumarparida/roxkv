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

var ErrIncompleteRESP = errors.New("ErrIncompleteRESP3")

func simpleReadError(data []byte) (string,int,error){
	pos := 1

	for; data[pos] != '\r'; pos++{
		if pos == len(data)-1 && data[pos] != '\r' {
			return "" , 0 , ErrIncompleteRESP
		}
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
		if pos == len(data)-1 && data[pos] != '\r' {
			return "", 0 , ErrIncompleteRESP
		}
	}

	return string(data[1:pos]) , pos+2 , nil
}


func readInteger(data []byte)  (int64,int,error){
	pos := 1

	for; data[pos] != byte('\r'); pos++{
		
	}

	i64 , err := strconv.ParseInt(string(data[1:pos]), 10 , 64)

	if err != nil {
		return  0 , 0 , err
	}


	

	return i64 , pos+2 , nil
}

func ReadLength(data []byte) (int,int){
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

func readBulkString(data []byte)   ([]byte,int,error){
	
	pos := 1

	
	// returns the length of the string incoming and the delta value which is the number of charachters read after the first char
	length , delta := ReadLength(data[pos:])
	pos = pos + delta
	fmt.Println("Bulk length:", length)
	fmt.Println("Remaining:", len(data)-pos)
	if pos+length+2 > len(data) {
   		 return nil, 0,  ErrIncompleteRESP
	}
	// starts slicing aftre the \r\n and ends on the exact length provided
	return  data[pos:(pos+length)] , pos+ length + 2 , nil
}


func readArray(data []byte)  (any,int,error){


	pos := 1

	totalElems , delta := ReadLength(data[pos:])

	pos += delta

	var elements []any = make([]any, totalElems);

	for i := range totalElems{

		val , posi , err := DecodeOne(data[pos:])
	
		if err != nil {
			if err == ErrIncompleteRESP {
				return  nil , 0 , ErrIncompleteRESP	
			}
			return nil , 0 , errors.New("Error while parsing elements")
		}

		elements[i] = val
		pos += posi
	}

	return  elements , pos , nil
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
		return 0.0 , 0 , errors.New("Error Parsing the decimal")
	}


	return f64 , pos+2 , nil
}


func DecodeOne(data []byte) (any,int,error){

	if len(data) == 0 {
		return  nil, 0 , errors.New("No Data Found")
	}

	switch data[0]{
		case SIMPLESTRING:
			return readSimpleString(data)
		case INTERGER:
			return readInteger(data)
		case BULKSTRING:
			return readBulkString(data)
		case ARRAY:
			return readArray(data)
		case BOOLEAN:
			return readBoolean(data)
		case DOUBLES:
			return readDoubles(data)
		case ERROR:
			return simpleReadError(data)
		default:
			return 	nil , 0 , errors.New("Data Type Didn't match")


	}
	


}


func Decode(data []byte) (any,int,error){
	if len(data) == 0 {
		return  nil , 0 , errors.New("No Data Found")
	}

	value , idx , err := DecodeOne(data)

	return  value, idx , err
}


func DecodeArrayString(data []byte, n int) ([][]byte, int, error) {

    value, idx, err := Decode(data[:n])

    if err != nil {
        return nil, 0, err
    }

    array, ok := value.([]any)
    if !ok {
        return nil, 0, errors.New("expected RESP array")
    }

    tokens := make([][]byte, len(array))

    for i := range array {

        s, ok := array[i].([]byte)
        if !ok {
            return nil, 0, errors.New("expected bulk string")
        }

        tokens[i] = s
    }

    return tokens, idx, nil
}