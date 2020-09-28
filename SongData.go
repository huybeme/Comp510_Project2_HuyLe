package main

type SongData struct {
	Error		bool		`json:"error"`
	Message		string		`json:"message,omitempty"`
	Response 	Response	`json:"response"`
}

type Response struct {
	Result []Result	`json:"results"`
}

type Result struct {
	ID 			string	`json:"id"`
	Name 		string	`json:"name"`
}