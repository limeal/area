package utils

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Name, Count
var Petname = map[string]int{
	"cat":     0,
	"dog":     0,
	"bird":    0,
	"fish":    0,
	"rabbit":  0,
	"hamster": 0,
	"snake":   0,
	"lizard":  0,
	"tiger":   0,
	"lion":    0,
}

// It generates a random number, then iterates through a map of strings and integers, and if the random
// number matches the current iteration, it increments the integer value of the map and returns the
// string key and the new integer value as a string
func GenerateNameForUser() string {

	seed := rand.NewSource(time.Now().UnixNano())
	base := rand.New(seed)
	// Generate a random number
	random := base.Intn(len(Petname))
	i := 0
	name := ""
	for key, value := range Petname {
		if i == random {
			Petname[key]++
			name = (key + strconv.Itoa(value))
		}
		i++
	}
	return name
}

// It generates a random string of a given length
func GenerateRandomString(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// It replaces all the keys in the data map with their values in the content string
func GenerateFinalComponent(content string, data map[string]interface{}, allowed []string) string {

	if len(allowed) == 0 {
		for key, value := range data {
			switch value.(type) {
			case string:
				content = strings.Replace(content, "{{"+key+"}}", value.(string), -1)
				break
			case int:
				content = strings.Replace(content, "{{"+key+"}}", strconv.Itoa(value.(int)), -1)
				break
			case float64:
				content = strings.Replace(content, "{{"+key+"}}", strconv.FormatFloat(value.(float64), 'f', -1, 64), -1)
				break
			case bool:
				content = strings.Replace(content, "{{"+key+"}}", strconv.FormatBool(value.(bool)), -1)
				break
			}
		}
		return content
	}

	for key, value := range data {
		for _, allowedKey := range allowed {
			if regexp, err := regexp.MatchString(allowedKey, key); err == nil && regexp {
				ok := false
				switch value.(type) {
				case string:
					content = strings.Replace(content, "{{"+key+"}}", value.(string), -1)
					ok = true
					break
				case int:
					content = strings.Replace(content, "{{"+key+"}}", strconv.Itoa(value.(int)), -1)
					ok = true
					break
				case float64:
					content = strings.Replace(content, "{{"+key+"}}", strconv.FormatFloat(value.(float64), 'f', -1, 64), -1)
					ok = true
					break
				case bool:
					content = strings.Replace(content, "{{"+key+"}}", strconv.FormatBool(value.(bool)), -1)
					ok = true
					break
				}
				if ok {
					break
				}
			}
		}
	}

	return content
}

// It takes a Fiber context and an auth service name, and returns an authorization object and an error
func VerifyRoute(c *fiber.Ctx, authService string) (*models.Authorization, error) {

	account := c.Locals("account").(models.Account)

	var authorization models.Authorization
	if result := postgres.DB.Where(&models.Authorization{AccountUUID: account.UUID, AuthService: authService}).First(&authorization); result.Error != nil {
		return nil, result.Error
	}

	return &authorization, nil
}
