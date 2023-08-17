package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/pkg/browser"
)

func main() {

	mode := ""

	prompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Please select one",
		Options:  []string{"make a link", "decrypt a link"},
		Help:     "literally just 2 options",
		PageSize: 0,
	}

	survey.AskOne(prompt, &mode, survey.WithValidator(survey.Required))

	switch mode {
	case "make a link":
		encrypt()
	case "decrypt a link":
		decrypt()
	}

}

func encrypt() {

	raw_url := ""

	url_prompt := &survey.Input{
		Message: "gib link",
		Help:    "input the link to encrypt",
	}

	survey.AskOne(url_prompt, &raw_url, survey.WithValidator(survey.Required))

	encryption_key := ""

	key_prompt := &survey.Input{
		Message: "gib key, this is gonna also be the url (yes this is not secure)",
		Help:    "the string that is gonna be used as both the url and the encryption key",
		Default: "kociumba",
	}

	survey.AskOne(key_prompt, &encryption_key)

	mumbo_jumbo(raw_url, encryption_key)

}

func decrypt() {

	raw_url := ""
	processed_url := ""

	url_prompt := &survey.Input{
		Message: "gib link",
		Help:    "input the link to decrypt",
	}

	survey.AskOne(url_prompt, &raw_url, survey.WithValidator(survey.Required))

	encryption_key := ""

	dotIndex := strings.Index(raw_url, ".")
	if dotIndex != -1 {
		encryption_key = raw_url[:dotIndex]
	}

	slashIndex := strings.Index(raw_url, "/")
	if slashIndex != -1 {
		processed_url = strings.TrimPrefix(raw_url[slashIndex:], "/")
	}

	// this is gucci
	// works fine

	mumbo_jumbo_reverse(processed_url, encryption_key)

}
func mumbo_jumbo(raw_url string, encryption_key string) error { // bro fuck this func name

	encryption_key_corrected := ""

	if encryption_key == "" {
		encryption_key_corrected = "kociumba"
	} else {
		encryption_key_corrected = encryption_key
	}

	runes := []rune(encryption_key_corrected)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	encryption_key_reversed := string(runes)

	encryption_hash := make([]int, len(encryption_key_reversed))

	for i, char := range encryption_key_reversed {
		encryption_hash[i] = int(char)
	}

	encryption_hash_full := 1

	for _, num := range encryption_hash {
		encryption_hash_full *= num
	}

	encryptor(raw_url, encryption_hash_full, encryption_key)

	return nil
}

func mumbo_jumbo_reverse(processed_url string, encryption_key string) error {

	encryption_key_corrected := ""
	if encryption_key == "" {
		encryption_key_corrected = "kociumba"
	} else {
		encryption_key_corrected = encryption_key
	}

	runes := []rune(encryption_key_corrected)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	encryption_key_reversed := string(runes)

	encryption_hash := make([]int, len(encryption_key_reversed))
	for i, char := range encryption_key_reversed {
		encryption_hash[i] = int(char)
	}

	encryption_hash_full := 1

	for _, num := range encryption_hash {
		encryption_hash_full *= num
	}

	decryptor(processed_url, encryption_hash_full, encryption_key)

	// println(processed_url)
	// println(encryption_hash_full)
	// println(encryption_key)

	return nil
}

func encryptor(raw_url string, encryption_hash_full int, encryption_key string) error {

	//println(raw_url)

	inted_url := make([]int, len(raw_url))

	for i, char := range raw_url {
		inted_url[i] = int(char)
	}

	encrypted_url := ""

	for _, num := range inted_url {
		encrypted_url += strconv.Itoa(encryption_hash_full * num)
		encrypted_url += "_" // thi is 😬
	}

	// println(encrypted_url)

	shortened_url := fmt.Sprintf(encryption_key+".dev/%s", encrypted_url)

	println("<====================>")
	println(shortened_url)
	println("this is encrypted (poorly), this should be usable only back throught this script")
	println("<====================>")

	return nil
}

func decryptor(processed_url string, encryption_hash_full int, encryption_key string) error {

	// println("processed_url: ", processed_url)
	// println("encryption_key: ", encryption_key)
	// println("encryption_hash_full: ", encryption_hash_full)

	encrypted_url := strings.TrimPrefix(processed_url, encryption_key+".dev/")
	encrypted_url = strings.TrimSuffix(encrypted_url, "/")

	encrypted_numbers := strings.Split(strings.Trim(encrypted_url, "_"), "_") // this is also 😬

	// println("encrypted_numbers: ", encrypted_numbers)

	decrypted_url := ""

	for _, numStr := range encrypted_numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return err
		}

		decrypted_url += string(rune(num / encryption_hash_full))
	}

	println("<====================>")
	println(decrypted_url)
	println("this is your decrypted URL")
	println("<====================>")

	err := opener(decrypted_url)
	if err != nil {
		return err
	}

	return nil
}

func opener(decrypted_url string) error {

	if !strings.HasPrefix(decrypted_url, "https://") {
		decrypted_url = "https://" + decrypted_url
	}

	matched, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`, decrypted_url)

	if !matched {
		return nil
	}

	browser.OpenURL(decrypted_url)

	return nil
}
