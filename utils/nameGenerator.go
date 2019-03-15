package utils

import (
	"math/rand"
	"strings"
	"bytes"
	"time"
)

type Element struct {
	Content		string
	Weight		float32
}

func randomWeightedChoice(elements []Element) (Element) {
    rand.Seed(time.Now().UnixNano())
    n := float32(rand.Float64())
    for _, element := range elements {
        n -= element.Weight
        if n <= 0 {
            return element
        }
    }
    return Element{Content: "", Weight: 0}
}

func stringInSlice(str string, list []string) bool {
 	for _, v := range list {
 		if v == str {
 			return true
 		}
 	}
 	return false
}

func removeDuplicates(source string, letters []string) (string) {
	destination := source
	for _, character := range source {
		scharacter := string(character)
			
		pi := strings.Index(source, scharacter) - 1
		ni := strings.Index(source, scharacter) + 1

		if stringInSlice(scharacter, letters) && pi >= 0 && stringInSlice(string(source[pi]), letters) {
			destination = destination[:pi] + destination[pi+1:]
		} 
		if stringInSlice(scharacter, letters) && ni < len(source) && stringInSlice(string(source[ni]), letters) {
			destination = destination[:ni] + destination[ni+1:]
		}
	}
	return destination
}

func generateFrequencies(names []string) ([]Element) {
	var elements []Element

	for _, name := range names {
		name = strings.ToLower(strings.Split(name, " ")[0])
		last_character_read := ""

		for _, character := range name {
			// Convert character (rune) to string
			scharacter := string(character)
			
			// Get previous character in name
			index := strings.Index(name, scharacter)
			previous_character := ""
			if index == 0 {
				previous_character = string(name[len(name) - 1])
			} else {
				previous_character = string(name[index - 1])
			}

			// Update frequencies (or initialize)
			if previous_character == last_character_read {
				// If element found, update frequencies
				is_found := false
				for i := 0; i < len(elements); i++ {
					element := &elements[i]
					if element.Content == last_character_read + scharacter {
						element.Weight += 1
						is_found = true
						break
					}
				}

				// Else, initialize
				if !is_found {
					elements = append(elements, Element{
						Content: last_character_read + scharacter,
						Weight: 1,
					})
				}
			} else {
				last_character_read = scharacter
			} 
		}
	}
	
	// Make total sum of frequencies
	sum := float32(0)
	for _, element := range elements {
		sum += element.Weight
	}

	// Create REAL frequencies (weight / sum)
	for i := 0; i < len(elements); i++ {
		element := &elements[i]
		element.Weight = element.Weight / sum
	}

	return elements
}

func generatePlanetName(elements []Element) (string) {
	var nameBuffer bytes.Buffer
	for iterator := 0; iterator < rand.Intn(8 - 3) + 3; iterator++ {
		nameBuffer.WriteString(randomWeightedChoice(elements).Content)
	}

	name := strings.Title(nameBuffer.String())
	// name = removeDuplicates(name, []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "v", "w", "x", "z"})
	// name = removeDuplicates(name, []string{"a", "e", "i", "o", "u", "y"})

	return name
}