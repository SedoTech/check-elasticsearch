package structs

type ServiceTweet struct {
	Timestamp	string `json:"@timestamp"`
	Level       string `json:"level"`
	Environment string `json:"environment"`
	Application string `json:"application"`
	Facility    string `json:"facility"`
	Host        string `json:"host"`
	Severity    int    `json:"severity"`
	Message     string `json:"message"`
}

type EcommerceTweet struct {
	Timestamp	string `json:"timestamp"`
	Level       int    `json:"level"`
	Environment string `json:"environment"`
	Application string `json:"application"`
	Facility    string `json:"facility"`
	Host        string `json:"host"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
}