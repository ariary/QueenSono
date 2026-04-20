package message

import (
	"fmt"
	"strconv"
	"strings"
)

// QueenSonoMarshall prepends each chunk with its index: "i,data".
func QueenSonoMarshall(dataSlice []string) []string {
	nDataSlice := make([]string, len(dataSlice))
	for i, s := range dataSlice {
		nDataSlice[i] = strconv.Itoa(i) + "," + s
	}
	return nDataSlice
}

// QueenSonoUnmarshall parses a "index,data" message.
// Returns an error if the message has no comma or an invalid index.
func QueenSonoUnmarshall(msg string) (nMsg string, index int, err error) {
	parts := strings.SplitN(msg, ",", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("QueenSonoUnmarshall: missing comma in %q", msg)
	}
	index, err = strconv.Atoi(parts[0])
	if err != nil {
		return "", 0, fmt.Errorf("QueenSonoUnmarshall: invalid index %q: %w", parts[0], err)
	}
	return parts[1], index, nil
}
