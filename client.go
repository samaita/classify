package identifier

type Client struct {
}

type Classification struct {
	Likely bool
	Score  float64
}

func InitWithFile() (Client, error) {
	// Check DB connection
	// Check filepath for model & library
	// Generate Default if none exist
	return Client{}, nil
}
