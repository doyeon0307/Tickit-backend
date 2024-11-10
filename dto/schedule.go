package dto

type ScheduleCalendarPreviewDTO struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Date  string `json:"date"`
}

type ScheduleTicketPreviewDTO struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type ScheduleDTO struct {
	Date      string `json:"date" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Number    int    `json:"number"`
	Image     string `json:"image"`
	Thumbnail bool   `json:"thumbmail"`
	Location  string `json:"location"`
	Time      string `json:"time"`
	Seat      string `json:"seat"`
	Casting   string `json:"casting"`
	Company   string `json:"company"`
	Link      string `json:"link"`
	Memo      string `json:"memo"`
}

type ScheduleResponseDTO struct {
	Id        string `json:"id"`
	Date      string `json:"date"`
	Title     string `json:"title"`
	Number    int    `json:"number"`
	Image     string `json:"image"`
	Thumbnail bool   `json:"thumbmail"`
	Location  string `json:"location"`
	Time      string `json:"time"`
	Seat      string `json:"seat"`
	Casting   string `json:"casting"`
	Company   string `json:"company"`
	Link      string `json:"link"`
	Memo      string `json:"memo"`
}
