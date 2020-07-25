package main

import (
	"encoding/json"
	"os"

	"github.com/Ahsan-J/HexParser/teltonika"
)

func main() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	encoder.Encode(teltonika.GenerateHex("getver", "0E", 352093081452251))
}
