package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)


func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func evaluateExpression(expr string) (float64, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	return result.(float64), nil
}

func generateCombinations(numbers []string, target float64) []string {
	operators := []string{"+", "-", "*", "/"}
	var results []string

	var recurse func(index int, currentExpr string)
	recurse = func(index int, currentExpr string) {
		if index == len(numbers)-1 {
			currentExpr += numbers[index]
			result, err := evaluateExpression(currentExpr)
			if err == nil && result == target {
				results = append(results, currentExpr)
			}
			return
		}
		currentExpr += numbers[index]
		for _, op := range operators {
			recurse(index+1, currentExpr+op)
		}
	}

	recurse(0, "")
	return results
}

func parseInput(input string) []string {
	var numbers []string
	var currentNumber strings.Builder

	for _, char := range input {
		if char >= '0' && char <= '9' {
			currentNumber.WriteRune(char)
		} else {
			if currentNumber.Len() > 0 {
				numbers = append(numbers, currentNumber.String())
				currentNumber.Reset()
			}
		}
	}

	if currentNumber.Len() > 0 {
		numbers = append(numbers, currentNumber.String())
	}

	return numbers
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	numbers := r.URL.Query().Get("numbers")
	target := r.URL.Query().Get("target")

	if numbers == "" || target == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	targetFloat, err := strconv.ParseFloat(target, 64)
	if err != nil {
		http.Error(w, "Invalid target", http.StatusBadRequest)
		return
	}

	numberList := parseInput(numbers)
	if len(numberList) < 3 || len(numberList) > 7 {
		http.Error(w, "Number count must be between 3 and 7", http.StatusBadRequest)
		return
	}

	combinations := generateCombinations(numberList, targetFloat)
	fmt.Fprintf(w, "Valid combinations: %v\n", combinations)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/solve", solveHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
