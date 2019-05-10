package main

import "flag"
import "fmt"
import "os"
import "time"
import "math/rand"
import "encoding/csv"

type quizProblem struct {
	question string
	answer   string
}

func checkErrors(e error, message string) {
	if e != nil {
		fmt.Println(message, e)
		os.Exit(1)
	}
}

func shuffleQuiz(quiz []quizProblem) []quizProblem {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(quiz), func(i, j int) {
		quiz[i], quiz[j] = quiz[j], quiz[i] 
	})

	return quiz
}

func parseCSVFileToStruct(csvLines [][]string) []quizProblem {
	parsedCSVSlice := make([]quizProblem, len(csvLines))

	for i := 0; i < len(csvLines); i++ {
		parsedCSVSlice[i] = quizProblem{
			question: csvLines[i][0],
			answer:   csvLines[i][1],
		}
	}
	return parsedCSVSlice
}

func startQuiz(quiz []quizProblem) int {
	var score int

	for i := 0; i < len(quiz); i++ {
		var answer string
		fmt.Println(quiz[i].question)
		fmt.Scanf("%s\n", &answer)
		if answer == quiz[i].answer {
			score++
		}
	}

	return score
}

func main() {
	// flag for user input csv file
	csvFileNameInput := flag.String("csv", "problems.csv", "Name of the quiz source csv file")
	flag.Parse()

	csvFile, errorInOpeningFile := os.Open(*csvFileNameInput)
	if errorInOpeningFile != nil {
		checkErrors(errorInOpeningFile, "Cannot open file")
	}

	fmt.Println("Quiz loaded from file :", *csvFileNameInput)

	csvReader := csv.NewReader(csvFile)
	csvFileLines, errInParsingCSV := csvReader.ReadAll()
	if errInParsingCSV != nil {
		checkErrors(errInParsingCSV, "Error in parsing CSV file")
	}

	quiz := parseCSVFileToStruct(csvFileLines)
	quiz = shuffleQuiz(quiz)

	score := startQuiz(quiz)

	fmt.Printf("You scored %d out of %d\n", score, len(quiz))

}
