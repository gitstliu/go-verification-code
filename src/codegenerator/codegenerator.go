package codegenerator

import (
	//	"bufio"
	//	"bytes"
	//	"errors"
	//	"fmt"
	"imagegenerator"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/gitstliu/log4go"
)

var wordsList = []string{}
var wordsListLength = 0

func InitWords(wordFileName string) error {
	wordFile, openFileErr := ioutil.ReadFile(wordFileName)
	if openFileErr != nil {
		return openFileErr
	}

	//	wordsMeta, readWordsErr := ioutil.ReadAll(wordFile)

	//	if readWordsErr != nil {
	//		return readWordsErr
	//	}

	wordsList = strings.Split(string(wordFile), "\n")
	wordsListLength = len(wordsList)
	//	for index := 0; index < wordsListLength; i++ {
	//		wordsList[index] = strings.TrimSpace()wordsList[index].
	//	}

	//	reader := bufio.NewReader(bytes.NewReader(wordFile))
	//	for line, _, readLineErr := reader.ReadLine(); readLineErr == nil; {
	//		wordsList = append(wordsList, string(line))
	//	}

	return nil
}

func GetRandWords() (string, string) {
	for index := 0; index < 5; index++ {
		keyIndex := rand.Intn(wordsListLength)
		noiceIndex := rand.Intn(wordsListLength)

		log4go.Debug("keyIndex = %v, noiceIndex = %v", keyIndex, noiceIndex)

		key := wordsList[keyIndex]
		noice := wordsList[noiceIndex]

		log4go.Debug("key = %v, noice = %v", key, noice)

		noiceKey := key + noice
		noiceKeyMeta := []rune(noiceKey)

		checkOK := true

		for _, currMetaR1 := range noiceKeyMeta {
			startIndex := strings.Index(noiceKey, string(currMetaR1))
			endIndex := strings.LastIndex(noiceKey, string(currMetaR1))

			log4go.Debug("startIndex = %v, endIndex = %v", startIndex, endIndex)
			log4go.Debug("noiceKey = %v, currMetaR1 = %v", noiceKey, string(currMetaR1))
			log4go.Debug("noiceKey = %v, currMetaR1 = %v", noiceKey, currMetaR1)

			if startIndex != endIndex {
				checkOK = false
				break
			}
			//			if !checkOK {
			//				break
			//			}
		}

		if checkOK {
			return key, noice
		}
	}

	return "", ""
}

func CreateNoiceConfusionCodeImage(keyWord string, noiceWord string, lineCount int, columnCount int, fg string, bg string, scale float64) (*imagegenerator.NoiceRunePosition, []byte, error) {

	position, noiceErr := noiceConfusionWord(keyWord, noiceWord, lineCount, columnCount)
	log4go.Debug("position.ConfusionWords = %v", position.ConfusionWords)
	position.KeyWord = keyWord
	if noiceErr != nil {
		return nil, nil, noiceErr
	}

	// CreateCodeImageByPosition(xCount int, yCount int, stepCell int, dpi float64, fontSize float64, hinting string, spacing float64, position *NoiceRunePosition, fgValue string, bgValue string)
	currImage, imageErr := imagegenerator.CreateCodeImageByPosition(
		columnCount,   //xCount
		lineCount,     //yCount
		int(50*scale), //stepCell
		72,            //dpi
		20*scale,      //fontSize
		"none",        //hinting
		1.5,           //spacing
		position,      //position
		fg,            //fgValue
		bg)            //bgValue

	if imageErr != nil {
		return nil, nil, imageErr
	}

	buffer, bufferErr := imagegenerator.ImageToPngBytes(currImage)
	return position, buffer, bufferErr
}

func noiceConfusionWord(keyWord string, noiceWord string, lineCount int, columnCount int) (*imagegenerator.NoiceRunePosition, error) {
	result := &imagegenerator.NoiceRunePosition{Position: []imagegenerator.Position{}, WordIndex: []int{}, LineCount: lineCount, ColumnCount: columnCount}
	result.ConfusionWords = confusionWord(keyWord + noiceWord)
	log4go.Debug("result.ConfusionWords = %v", string(result.ConfusionWords))
	//	newWords := string(result.ConfusionWords)
	keyWordMeta := []rune(keyWord)
	for _, currWord := range keyWordMeta {
		for index, currNewWord := range result.ConfusionWords {
			if currNewWord == currWord {
				log4go.Debug("CurrWordIndex = %v", index)
				result.WordIndex = append(result.WordIndex, index)
			}
		}
		//		currIndex := strings.IndexRune(newWords, currWord)
		//		log4go.Debug("newWords = %v, currWord=%v", newWords, string(currWord))
		//		log4go.Debug("currIndex = %v", currIndex)
		//		if currIndex == -1 {
		//			return nil, errors.New(fmt.Sprintf("Char [%v] is not in the string [%v]", currWord, newWords))
		//		}
		//		result.WordIndex = append(result.WordIndex, currIndex)
	}

	totalCount := lineCount * columnCount

	log4go.Debug("columnCount = %v", columnCount)

	for i := 0; i < totalCount; i++ {
		wordPos := imagegenerator.Position{i / columnCount, i % columnCount}
		log4go.Debug("wordPos = %v", wordPos)
		result.Position = append(result.Position, wordPos)
	}

	return result, nil
}

func confusionWord(word string) []rune {
	result := []rune{}
	wordRunes := []rune(word)
	wordRunesLength := len(wordRunes)
	pos := rand.Perm(wordRunesLength)

	for index := 0; index < wordRunesLength; index++ {
		result = append(result, wordRunes[pos[index]])
	}
	return result
}
