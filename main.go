package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	fmt.Println("You will be taking a quick math quiz.")

	count := 0
	correct := 0
	var i int

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("You have 30 seconds to complete the quiz. Start now!")

	for {
		select {

		case <-ctx.Done():
			fmt.Println("\nTime's up!")
			fmt.Printf("You got %d out of %d correct.\n", correct, count)
			fmt.Printf("You scored : %.2f%%\n", (float64(correct)/float64(count))*100)
			return

		default:
			for _, record := range records {
				count++
				fmt.Print("Question ", count, " : ")
				fmt.Print(record[0], " = ")
				fmt.Scan(&i)

				answer, err := strconv.Atoi(record[1])
				if err != nil {
					fmt.Println("Error converting answer to integer:", err)
					continue
				}
				if i == answer {
					correct++
				}
				select {
				case <-ctx.Done():
					fmt.Println("\n⏳ Time's up! Quiz ended.")
					return
				default:
				}
			}
		}
	}
}
