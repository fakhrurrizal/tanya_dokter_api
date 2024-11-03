package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"tanya_dokter_app/config"
	"time"
	"unicode"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/sqids/sqids-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func StripTagsFromStruct(input interface{}) {
	structValue := reflect.ValueOf(input).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)

		if fieldValue.Kind() == reflect.String {
			originalValue := fieldValue.String()
			strippedValue := strip.StripTags(originalValue)
			fieldValue.SetString(strippedValue)
		} else if fieldValue.Kind() == reflect.Struct {
			StripTagsFromStruct(fieldValue.Addr().Interface())
		}
	}
}

func ObjectToString(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// time to string
func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Respond(code int, data interface{}, message string) (response Response) {
	return Response{
		Status:  code,
		Message: message,
		Data:    data,
	}
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZqwertyuiopasdfghjklzxcvbnm0123456789!@#$%^&*?"
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

func ConvertToKebabCase(input string) string {
	var result []rune

	input = strings.ToLower(input)

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
		} else if len(result) > 0 && result[len(result)-1] != '-' {
			result = append(result, '-')
		}
	}

	return string(result)
}

func ConvertToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = TitleCase(parts[i])
	}
	return strings.Join(parts, "")
}

func TitleCase(str string) string {
	tc := cases.Title(language.Indonesian)
	return tc.String(str)
}

func RemoveDuplicates(str string) (data string) {
	s := strings.Split(str, ",")
	uniqueStrings := make([]string, 0, len(s))
	seen := make(map[string]bool)
	for _, str := range s {
		if !seen[str] {
			if str != "" {
				seen[str] = true
				uniqueStrings = append(uniqueStrings, str)
			}
		}
	}
	for _, v := range uniqueStrings {
		data += v + ","
	}
	data = strings.TrimRight(data, ",")
	return
}

func GetNumberFromStr(str string) (num int) {
	numStr := ""
	for _, char := range str {
		if char >= '0' && char <= '9' {
			numStr += string(char)
		}
	}
	num, _ = strconv.Atoi(numStr)

	return
}

func Average(numbers []float64) float64 {
	var sum float64
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

func StripTags(str string) string {
	str = strip.StripTags(str)
	return str
}

func LastId(table string, columnId ...string) (id int64) {
	type OnlyId struct {
		ID int64
	}
	var last OnlyId

	if len(columnId) > 0 {
		col := columnId[0]
		config.DB.Table(table).Select(fmt.Sprintf("%s as id", col)).Order(fmt.Sprintf("%s desc", col)).Limit(1).Scan(&last)
	} else {
		config.DB.Table(table).Order("id desc").Limit(1).Scan(&last)
	}
	id = last.ID + 1
	return
}

func GenerateKeyStruct(data interface{}) string {
	value := reflect.ValueOf(data)
	key := ""
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i).Interface()
		key += fmt.Sprintf(":%v", fieldValue)
	}
	return key
}

func GenerateRandomPIN() string {
	rand.Seed(time.Now().UnixNano())
	charset := "0123456789"
	randomBytes := make([]byte, 6)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(randomBytes)
	return randomString
}

func MakeKey(values ...interface{}) string {
	var keyParts []string

	for _, value := range values {
		switch v := value.(type) {
		case string:
			keyParts = append(keyParts, v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			keyParts = append(keyParts, fmt.Sprintf("%v", v))
		case bool:
			keyParts = append(keyParts, fmt.Sprintf("%t", v))
		default:
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Struct {
				for i := 0; i < val.NumField(); i++ {
					field := val.Field(i)
					keyParts = append(keyParts, fmt.Sprintf("%v", field.Interface()))
				}
			}
		}
	}

	return strings.Join(keyParts, "_")
}

func ContainsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func EndcodeID(id int) (response string) {
	s, _ := sqids.New(sqids.Options{
		MinLength: 64,
	})
	response, _ = s.Encode([]uint64{uint64(id)})

	return
}

func DecodeID(encodedID string) (response int) {
	s, _ := sqids.New(sqids.Options{
		MinLength: 64,
	})

	responseData := s.Decode(encodedID)
	if len(responseData) > 0 {
		response = int(responseData[0])
	}

	return
}
