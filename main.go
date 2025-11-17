package main

import (
	"bufio"
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

	fmt.Println("Press Enter to start the quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	count := 0
	correct := 0

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("You have 30 seconds to complete the quiz. Start now!")

	for _, record := range records {
		count++

		fmt.Print("Question ", count, ": ", record[0], " = ")

		answerCh := make(chan int)

		go func() {
			for {
				var usrInput int
				_, err := fmt.Scan(&usrInput)

				if err != nil {
					fmt.Println("Not a number, try again", err)
					var dump string
					fmt.Scanln(&dump)
					fmt.Print("Question ", count, ": ", record[0], " = ")
					continue
				}
				answerCh <- usrInput
				return
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("\nTime's up!")
			printResults(count-1, correct)
			return

		case usrAnswer := <-answerCh:
			correctAnswer, _ := strconv.Atoi(record[1])
			if err != nil {
				fmt.Println("Invalid answer format. Skipping question.")
				continue
			}

			if usrAnswer == correctAnswer {
				correct++
			}
		}
	}

	fmt.Println("Quiz completed!")
	printResults(count, correct)
}

func printResults(total, correct int) {
	fmt.Printf("You answered %d out of %d questions correctly.\n", correct, total)
	if total > 0 {
		percentage := (float64(correct) / float64(total)) * 100
		fmt.Printf("Your score: %.2f%%\n", percentage)
	} else {
		fmt.Println("No questions were answered.")
	}
}
