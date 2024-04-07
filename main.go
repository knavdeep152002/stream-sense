package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	fmt.Println("Start redis test.")
	redisClient.Publish(ctx, "mychannel", "hello2")
	allFiles, _ := os.ReadDir("uploads/chunks/sf7xev")
	log.Println("All files: ", len(allFiles))
	x := 2
	for i := range allFiles {
		var file1 []byte
		if i == 0 {
			file1, _ = os.ReadFile("uploads/chunks/sf7xev/1")
		} else {
			file1, _ = os.ReadFile("uploads/chunks/test")
		}
		file2, _ := os.ReadFile("uploads/chunks/sf7xev/" + strconv.Itoa(x))
		x++
		file3 := append(file1, file2...)
		os.WriteFile("uploads/chunks/test", file3, 0644)
		log.Println("x: ", x)
	}

	// // iterate backwards to find the first difference
	// // j:=
	// for i, j := len(file1)-1, len(file2)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
	// 	if file1[i] != file2[j] {
	// 		log.Println(i, j)
	// 		break
	// 	} else {
	// 		log.Println("Matched", file1[i], file2[j])
	// 	}
	// }

	// // for i := range file1 {
	// // 	if file1[i] != file2[i] {
	// // 		log.Println(i)
	// // 		break
	// // 	} else {
	// // 		// log.Println("Matched", file1[i], file2[i])
	// // 	}
	// // }

	// log.Println("File1: ", file1[:26], len(file1))
	// log.Println("File2: ", file2[:26], len(file2))
	// log.Println("EOF file 1: ", file1[len(file1)-26:])
	// log.Println("EOF file 1: ", file2[len(file2)-26:])
}
