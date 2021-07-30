package main

type PwBody struct {
	Password       string `json:"pw" form:"password"`
	DaysLimit      int    `json:"days_limit" form:"days_limit"`
	ViewsRemaining int    `json:"views_remaining" form:"views_remaining"`
}
