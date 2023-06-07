package utils

func enoughBudget(numbers []AnonNumbers, availableMoney int) bool {
	for _, number := range numbers {
		availableMoney -= number.Value
	}

	return availableMoney > 0
}