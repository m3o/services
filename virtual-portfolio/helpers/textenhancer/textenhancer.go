package textenhancer

import (
	"fmt"
	"strings"
)

// Service is an implementation of TextEnhancer
type Service struct {
	// Eventually we'll use the users-srv and stocks-srv to provide additonal data
	// which can be displayed on hover. It'll also be used to determine which users
	// to notify.

	// users  users.UsersService
	// stocks stocks.StocksService
}

// Enhance takes a string of text, e.g. a post and adds the relevent metadata
// for the stocks and users who are tagged in that post.
func (s Service) Enhance(text string) (result string) {
	for _, comp := range s.componentsForString(text) {
		result += s.encodeComponent(comp)
	}
	return result
}

// ListTaggedUsers returns the usernames for the users tagged in a text
func (s Service) ListTaggedUsers(text string) (usernames []string) {
	for _, comp := range s.componentsForString(text) {
		if comp.kind == "User" {
			username := strings.Replace(comp.text, "@", "", 1)
			usernames = append(usernames, username)
		}
	}
	return usernames
}

type textComponent struct {
	text, kind string
}

func (s Service) encodeComponent(comp textComponent) string {
	switch comp.kind {
	case "User":
		username := strings.Replace(comp.text, "@", "", 1)
		return fmt.Sprintf("<&User:%v>%v<&/User>", username, comp.text)
	default:
		return comp.text
	}
}

func (s Service) componentsForString(text string) []textComponent {
	var result []textComponent

	var currentResultType string
	for _, char := range text {
		ascii := int(char)

		// Reached an @
		if ascii == 64 {
			result = append(result, textComponent{string(char), "User"})
			currentResultType = "User"
			continue
		}

		// The current tag should be broken if it's not standard text
		isNormalLetter := (ascii >= 65 && ascii <= 90) || (ascii >= 92 && ascii <= 122)
		if !isNormalLetter && currentResultType != "Text" {
			result = append(result, textComponent{string(char), "Text"})
			currentResultType = "Text"
			continue
		}

		// Ensure this isn't the first element in the string
		if currentResultType == "" {
			result = append(result, textComponent{string(char), "Text"})
			currentResultType = "Text"
		} else {
			result[len(result)-1].text += string(char)
		}
	}

	return result
}
