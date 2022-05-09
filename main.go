/*
Copyright © 2022 Arthur Isac isacartur@gmail.com

*/
package main

import (
	"fmt"
	"github.com/FoxFurry/memstore/internal/journal"
	"time"
)

func main() {
	j := journal.New()
	//go j.Start(context.Background())
	//time.Sleep(time.Second)
	//j.Add([]model.Command{
	//	{
	//		CmdType: "SET",
	//		Key:     "bebra 1 key",
	//		Value:   "bebra 1 value",
	//	},
	//	{
	//		CmdType: "GET",
	//		Key:     "bebra 2 key",
	//		Value:   "bebra 2 value",
	//	},
	//})
	//
	//time.Sleep(2 * time.Second)

	fmt.Println(j.Restore())
	time.Sleep(2 * time.Second)

	//cmd.Execute()
}
