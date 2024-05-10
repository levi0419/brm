package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"encoding/json"
    "net/http"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(7))
}

// NewUuid generates a uuid
func NewUuid() uuid.UUID {
	id := uuid.New()
	return id
}


func Generate() string {
    // Define the characters for generating the code
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    code := make([]byte, 6)

    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    for i := range code {
        code[i] = charset[rand.Intn(len(charset))]
    }
    return string(code)
}


// IPInfoResponse represents the structure of the response from ipinfo.io
type IPInfoResponse struct {
    IP      string `json:"ip"`
    Country string `json:"country"`
    Region  string `json:"region"`
    // Add more fields as needed
}

// GetCountryAndStateFromIP fetches country and state information based on the given IP address
func GetCountryAndStateFromIP(ip string) (country, state string, err error) {
    url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
    fmt.Println("url", url)
    resp, err := http.Get(url)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    var ipInfo IPInfoResponse
    if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
        return "", "", err
    }

    fmt.Println("response", ipInfo.Country)

    return ipInfo.Country, ipInfo.Region, nil
}

// List of third-world countries
var ThirdWorldCountries = map[string]bool{
    "AF": true, // Afghanistan
    "BD": true, // Bangladesh
    "BJ": true, // Benin
    "BF": true, // Burkina Faso
    "BI": true, // Burundi
    "KH": true, // Cambodia
    "CM": true, // Cameroon
    "CF": true, // Central African Republic
    "TD": true, // Chad
    "KM": true, // Comoros
    "CD": true, // Democratic Republic of the Congo
    "ER": true, // Eritrea
    "ET": true, // Ethiopia
    "GM": true, // Gambia
    "GN": true, // Guinea
    "GW": true, // Guinea-Bissau
    "HT": true, // Haiti
    "KP": true, // Korea (Democratic People’s Republic of)
    "LR": true, // Liberia
    "MG": true, // Madagascar
    "MW": true, // Malawi
    "ML": true, // Mali
    "MZ": true, // Mozambique
    "NE": true, // Niger
    "NG": true, // Nigeria
    "RW": true, // Rwanda
    "SL": true, // Sierra Leone
    "SO": true, // Somalia
    "SS": true, // South Sudan
    "SD": true, // Sudan
    "SY": true, // Syrian Arab Republic
    "TJ": true, // Tajikistan
    "TZ": true, // Tanzania
    "TG": true, // Togo
    "UG": true, // Uganda
    "YE": true, // Yemen (Republic of)
    "ZM": true, // Zambia
    "ZW": true, // Zimbabwe
}

// IsThirdWorldCountry checks if a given country is considered a third-world country
func IsThirdWorldCountry(country string) bool {
    // Check if the country is in the list of third-world countries
    _, ok := ThirdWorldCountries[country]
    return ok
}