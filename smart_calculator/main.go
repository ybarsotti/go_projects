package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "/exit":
			fmt.Println("Bye!")
			return
		case "/help":
			printHelp()
			continue
		case "":
			continue
		default:
			if strings.HasPrefix(input, "/") {
				fmt.Println("Unknown command")
				continue
			}
			result, err := calculateExpression(input)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		}
	}
}

func printHelp() {
	fmt.Println("The program calculates expressions involving addition and subtraction.")
	fmt.Println("You can enter expressions like '4 + 6 - 8' or '2 - 3 - 4'.")
	fmt.Println("To exit the program, type '/exit'.")
}

func calculateExpression(expression string) (int, error) {
	// Replace multiple consecutive minus signs with a single minus or plus if necessary
	expression = handleConsecutiveSigns(expression)

	tokens := strings.Fields(expression)
	if len(tokens) == 0 {
		return 0, nil
	}

	total, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid expression")
	}

	for i := 1; i < len(tokens); i += 2 {
		operator := tokens[i]
		value, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return 0, fmt.Errorf("Invalid expression")
		}

		switch operator {
		case "+":
			total += value
		case "-":
			total -= value
		default:
			return 0, fmt.Errorf("Invalid operator '%s'. Please use '+' or '-'.", operator)
		}
	}

	return total, nil
}

func handleConsecutiveSigns(expression string) string {
	// Convert multiple consecutive '-' signs to a single '-' if odd, or '+' if even
	for strings.Contains(expression, "--") || strings.Contains(expression, "-+") || strings.Contains(expression, "+-") || strings.Contains(expression, "++") {
		expression = strings.ReplaceAll(expression, "--", "+")
		expression = strings.ReplaceAll(expression, "-+", "-")
		expression = strings.ReplaceAll(expression, "+-", "-")
		expression = strings.ReplaceAll(expression, "++", "+")
	}
	return expression
}
