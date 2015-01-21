package main

type event struct {
	ID         string `json:"id"`
	Client     client `json:"client"`
	Check      check  `json:"check"`
	Occurences int    `json:"occurrences"`
	Action     string `json:"action"`
}

type client struct {
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	Subscriptions []string `json:"subscriptions"`
	Timestamp     int      `json:"timestamp"`
}

type check struct {
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Subscribers []string `json:"subscribers"`
	Interval    int      `json:"interval"`
	Issued      int      `json:"issued"`
	Executed    int      `json:"executed"`
	Output      string   `json:"output"`
	Status      int      `json:"status"`
	Duration    float64  `json:"duration"`
	History     []string `json:"history"`
}

