package main

import "flag"
import "fmt"
import "os"
import "time"
import "encoding/csv"
import "math/rand"
import "strings"

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

func startQuiz(quiz []quizProblem, timerInput int) int {
	var score int

	fmt.Printf("Quiz is starting, time limit : %ds \n ", timerInput)
	timer := time.NewTimer(time.Duration(timerInput) *  time.Second)

	
	quizloop:
	for i := 0; i < len(quiz); i++ {
		fmt.Printf("QNo.%d: %s = ", i+1, quiz[i].question)
		
		answeringChannel := make(chan string)
		go func(){
			var answer string
			fmt.Scanf("%s\n", &answer)
			answeringChannel <- strings.TrimSpace(answer)		
		}()

		select{
		case answerInput := <-answeringChannel:
			if answerInput == quiz[i].answer {
				score++
			}
		case <- timer.C :
			fmt.Println("Times up!")
			break quizloop
		}
	}

	return score
}

func main() {
	csvFileNameInput := flag.String("csv", "problems.csv", "Name of the quiz source csv file")
	timerInput := flag.Int("timer", 10, "Timer in seconds")
	flag.Parse()

	csvFile, errorInOpeningFile := os.Open(*csvFileNameInput)
	if errorInOpeningFile != nil {
		checkErrors(errorInOpeningFile, "Cannot open file")
	}
	defer csvFile.Close()

	fmt.Println("Loading file :", *csvFileNameInput)

	csvReader := csv.NewReader(csvFile)
	csvFileLines, errInParsingCSV := csvReader.ReadAll()
	if errInParsingCSV != nil {
		checkErrors(errInParsingCSV, "Error in parsing CSV file")
	}

	quiz := parseCSVFileToStruct(csvFileLines)
	quiz = shuffleQuiz(quiz)

	score := startQuiz(quiz, *timerInput)

	fmt.Printf("You scored %d out of %d\n", score, len(quiz))

}
