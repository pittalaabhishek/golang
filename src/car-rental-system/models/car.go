package models

type Car struct {
	ID              int
	Make            string
	Model           string
	Year            int
	LicensePlate    string
	RentalPricePerDay float64
	IsAvailable     bool
}

type Customer struct {
	Name           string
	ContactDetails string
	DriversLicense string
}

type Reservation struct {
	ID         int
	Customer   Customer
	CarID      int
	StartDate  string
	EndDate    string
	TotalPrice float64
	Paid       bool
}