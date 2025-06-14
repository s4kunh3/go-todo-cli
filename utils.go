package main

func wrapText (text string, width int) []string {
	var lines []string
	for len(text) > width {
		lines = append(lines, text[:width])
		text = text[width:]
	}
	lines = append(lines, text)
	return lines
}

func getID() int {
	nextID := 0
	for _, task := range tasks {
		if task.ID > nextID {
			nextID = task.ID
		}
	}
	return nextID + 1
}