package main

import (
	"errors"
	"log"
)

const (
	startCharLower   = 'a'
	startCharUpper   = 'A'
	startNumber      = '0'
	charNumber       = 26
	doubleCharNumber = 52
	numCharNumber    = 62
	maxCharNumber    = 14776336
)

func intToChar(idx int) rune {
	if idx < charNumber {
		return rune(startCharLower + idx)
	} else if idx < doubleCharNumber {
		return rune(startCharUpper + (idx - charNumber))
	} else if idx < numCharNumber {
		return rune(startNumber + (idx - doubleCharNumber))
	} else {
		return '-'
	}

}

// translateNumber translates input into set of 52-based numbers
func translateNumber(idx int) []int {
	result := make([]int, 0)

	var idxDiv int
	var idxMod int

	idxMod = idx
	for {
		idxDiv = idxMod % numCharNumber
		idxMod = idxMod / numCharNumber
		result = append(result, idxDiv)
		if idxMod == 0 {
			break
		}
	}

	if resLen := len(result); resLen < 4 {
		for idxRes := 0; idxRes < 4-resLen; idxRes++ {
			result = append(result, 0)
		}
	}
	return result
}

func getKeyByID(ID int) (string, error) {

	if ID < 0 || ID >= maxCharNumber {
		return "", errors.New("out of range")
	}

	keySet := translateNumber(ID)

	if len(keySet) < 4 {
		return "", errors.New("wrong key translation")
	}

	key := string([]rune{
		intToChar(keySet[0]),
		intToChar(keySet[1]),
		intToChar(keySet[2]),
		intToChar(keySet[3]),
	})

	log.Println("key generated:", key)
	return key, nil
}

func getNextKey(idx int) string {
	key, err := getKeyByID(idx)
	if err != nil {
		log.Println(err)
	}
	return key
}

func getMaxKeyNumber() int {
	return maxCharNumber
}
