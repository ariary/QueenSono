package message

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Function used to transfor a data chunked in the QueenSono format (ie each chunk has its index added behind)
func QueenSonoMarshall(dataSlice []string) (nDataSlice []string) {
	nDataSlice = make([]string, len(dataSlice))
	for i := 0; i < len(dataSlice); i++ {
		nDataSlice[i] = strconv.Itoa(i) + "," + dataSlice[i]
	}
	return nDataSlice
}

//Unmarshall a QueenSono message (ie parsing it to get the index of the packet and the content)
func QueenSonoUnmarshall(msg string) (nMsg string, index int) {
	qsMsg := strings.SplitN(msg, ",", 2)
	index, err := strconv.Atoi(qsMsg[0])

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "QueenSonoUnmarshall: failed to convert %s into int", qsMsg[0])
	}
	nMsg = qsMsg[1]
	return nMsg, index
}
