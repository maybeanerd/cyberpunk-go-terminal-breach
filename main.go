package main

import (
	"fmt"
	"log"
	"main/breachModel"
	"main/breachUI"

	"github.com/pkg/term"
)

// UI interaction
// Raw input keycodes
var up byte = 65    // A
var down byte = 66  // B
var right byte = 67 // C
var left byte = 68  // D
var escape byte = 27
var enter byte = 13
var keys = map[byte]bool{
	up:    true,
	down:  true,
	right: true,
	left:  true,
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

func main() {
	var breachSurfaceSize = 6
	var breachSurface = breachModel.GenerateBreachSurface(breachSurfaceSize)
	var breachSequence1 = breachModel.GenerateBreachSequence(3)
	var breachSequence2 = breachModel.GenerateBreachSequence(3)
	var breachSequence3 = breachModel.GenerateBreachSequence(3)

	var hoverRowIndex = 0
	var hoverColumnIndex = 0
	var currentSelectionModeRow = true

	var breachBuffer = breachModel.GenerateBreachBuffer(6)
	var currentBufferIndex = 0
	var printedLinesCount = 0

	breachUI.PrintLogo()
	breachUI.PrintInstructions()

	for {
		breachBufferWidth := breachUI.PrintBreachBuffer(breachBuffer, &printedLinesCount)
		fmt.Println()
		printedLinesCount++
		breachUI.PrintBreachSequenceTitle(&printedLinesCount)
		breach1Result := breachUI.PrintBreachSequence(breachSequence1, breachBuffer, breachBufferWidth, &printedLinesCount)
		breach2Result := breachUI.PrintBreachSequence(breachSequence2, breachBuffer, breachBufferWidth, &printedLinesCount)
		breach3Result := breachUI.PrintBreachSequence(breachSequence3, breachBuffer, breachBufferWidth, &printedLinesCount)
		fmt.Println()
		printedLinesCount++
		breachUI.PrintBreachSurface(breachSurface, hoverRowIndex, hoverColumnIndex, currentSelectionModeRow, &printedLinesCount)

		if breach1Result != 0 && breach2Result != 0 && breach3Result != 0 {
			break
		}
		if currentBufferIndex >= len(breachBuffer) {
			break
		}
		keyCode := getInput()
		if keyCode == escape {
			return
		} else if keyCode == enter {
			var currentBreachedHole = breachSurface[hoverRowIndex][hoverColumnIndex]
			if currentBreachedHole.IsFree {
				currentBreachedHole.IsFree = false
				currentSelectionModeRow = !currentSelectionModeRow
				breachBuffer[currentBufferIndex] = currentBreachedHole.Address
				currentBufferIndex++
			}
		} else if keyCode == up && !currentSelectionModeRow {
			hoverRowIndex = (hoverRowIndex + len(breachSurface) - 1) % len(breachSurface)
		} else if keyCode == down && !currentSelectionModeRow {
			hoverRowIndex = (hoverRowIndex + 1) % len(breachSurface)
		} else if keyCode == right && currentSelectionModeRow {
			hoverColumnIndex = (hoverColumnIndex + 1) % len(breachSurface)
		} else if keyCode == left && currentSelectionModeRow {
			hoverColumnIndex = (hoverColumnIndex + len(breachSurface) - 1) % len(breachSurface)
		}
		// If we're gonna redraw, we need to move the cursor back up the number of lines that need redrawing
		fmt.Printf("\033[%dA", printedLinesCount)
		printedLinesCount = 0
	}
}
