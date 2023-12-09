package dto

// Crypto represents the structure of the cryptocurrency data
type Crypto struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
