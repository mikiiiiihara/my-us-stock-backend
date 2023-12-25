package crypto

// ApiResponse represents the structure of the API response
type ApiResponse struct {
    Success   int `json:"success"`
    Data struct {
        Sell      string `json:"sell"`
        Buy       string `json:"buy"`
        Open      string `json:"open"`
        High      string `json:"high"`
        Low       string `json:"low"`
        Last      string `json:"last"`
        Vol       string `json:"vol"`
        Timestamp int64  `json:"timestamp"`
    } `json:"data"`
}