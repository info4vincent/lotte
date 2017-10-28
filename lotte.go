package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

func checkIfSpeechIsalreadyAvailable(rubensTextToSay string) (bool, string) {
	hasher := sha1.New()
	hasher.Write([]byte(rubensTextToSay))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	fileName := fmt.Sprintf("./%s.mp3", sha)
	fmt.Printf("mp3 file stored as:'%s'\n", fileName)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false, ""
	}
	return true, fileName
}

func main() {
	rubensTextToSay := "He'! Wie kietelt mij daar? Oke goed je hebt mijn aandacht. Ik zal zeggen wie ik ben. Ik heet Ruben. Hoe heet jij?"

	fileExist, fileName := checkIfSpeechIsalreadyAvailable(rubensTextToSay)
	if fileExist {
		fmt.Println("mp3 is already available...")
		return
	}
	sess := session.Must(session.NewSession())

	svc := polly.New(sess, aws.NewConfig().WithRegion("eu-west-1"))
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		SampleRate:   aws.String("16000"),
		Text:         aws.String(rubensTextToSay),
		TextType:     aws.String("text"),
		VoiceId:      aws.String("Ruben"),
	}

	result, err := svc.SynthesizeSpeech(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case polly.ErrCodeTextLengthExceededException:
				fmt.Println(polly.ErrCodeTextLengthExceededException, aerr.Error())
			case polly.ErrCodeInvalidSampleRateException:
				fmt.Println(polly.ErrCodeInvalidSampleRateException, aerr.Error())
			case polly.ErrCodeInvalidSsmlException:
				fmt.Println(polly.ErrCodeInvalidSsmlException, aerr.Error())
			case polly.ErrCodeLexiconNotFoundException:
				fmt.Println(polly.ErrCodeLexiconNotFoundException, aerr.Error())
			case polly.ErrCodeServiceFailureException:
				fmt.Println(polly.ErrCodeServiceFailureException, aerr.Error())
			case polly.ErrCodeMarksNotSupportedForFormatException:
				fmt.Println(polly.ErrCodeMarksNotSupportedForFormatException, aerr.Error())
			case polly.ErrCodeSsmlMarksNotSupportedForTextTypeException:
				fmt.Println(polly.ErrCodeSsmlMarksNotSupportedForTextTypeException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
	data, err := ioutil.ReadAll(result.AudioStream)
	ioutil.WriteFile(fileName, data, 0644)
}
