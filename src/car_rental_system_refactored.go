package main

import (
	"errors"
	"fmt"
	"sync"
)

// Car represents a rental car.
type Car struct {
	ID              int
	Make            string
	Model           string
	Year            int
	LicensePlate    string
	RentalPricePerDay float64
	IsAvailable     bool
}

// Customer represents a customer.
type Customer struct {
	Name            string
	ContactDetails  string
	DriversLicense  string
}

// Reservation holds booking details.
type Reservation struct {
	ID         int
	Customer   Customer
	CarID      int
	StartDate  string
	EndDate    string
	TotalPrice float64
	Paid       bool
}

// RentalSystem manages cars and reservations.
type RentalSystem struct {
	cars          map[int]*Car
	reservations  map[int]*Reservation
	mu            sync.Mutex
	reservationID int
}

// NewRentalSystem initializes the system.
func NewRentalSystem() *RentalSystem {
	return &RentalSystem{
		cars:         make(map[int]*Car),
		reservations: make(map[int]*Reservation),
	}
}

// AddCar adds a car to the system.
func (rs *RentalSystem) AddCar(car Car) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.cars[car.ID] = &car
}

// SearchCars returns available cars by criteria.
func (rs *RentalSystem) SearchCars(make string, maxPrice float64) []Car {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var results []Car
	for _, car := range rs.cars {
		if car.Make == make && car.RentalPricePerDay <= maxPrice && car.IsAvailable {
			results = append(results, *car)
		}
	}
	return results
}

// IsCarAvailable checks if a car is free on given dates.
func (rs *RentalSystem) IsCarAvailable(carID int, startDate, endDate string) bool {
	for _, res := range rs.reservations {
		if res.CarID == carID && res.StartDate <= endDate && res.EndDate >= startDate {
			return false
		}
	}
	return true
}

// CreateReservation books a car.
func (rs *RentalSystem) CreateReservation(customer Customer, carID int, startDate, endDate string) (*Reservation, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	car, exists := rs.cars[carID]
	if !exists || !car.IsAvailable {
		return nil, errors.New("car not available")
	}

	rs.reservationID++
	reservation := &Reservation{
		ID:         rs.reservationID,
		Customer:   customer,
		CarID:      carID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: car.RentalPricePerDay, // Assuming one-day reservation
	}
	
	rs.reservations[rs.reservationID] = reservation
	car.IsAvailable = false

	return reservation, nil
}

// ModifyReservation updates reservation dates.
func (rs *RentalSystem) ModifyReservation(reservationID int, newStartDate, newEndDate string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	res, exists := rs.reservations[reservationID]
	if !exists {
		return errors.New("reservation not found")
	}

	if !rs.IsCarAvailable(res.CarID, newStartDate, newEndDate) {
		return errors.New("car not available for new dates")
	}

	res.StartDate, res.EndDate = newStartDate, newEndDate
	return nil
}

// CancelReservation removes a reservation.
func (rs *RentalSystem) CancelReservation(reservationID int) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	res, exists := rs.reservations[reservationID]
	if !exists {
		return errors.New("reservation not found")
	}

	if car, exists := rs.cars[res.CarID]; exists {
		car.IsAvailable = true
	}

	delete(rs.reservations, reservationID)
	return nil
}

// ProcessPayment marks a reservation as paid.
func (rs *RentalSystem) ProcessPayment(reservationID int) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	res, exists := rs.reservations[reservationID]
	if !exists {
		return errors.New("reservation not found")
	}

	if res.Paid {
		return errors.New("reservation already paid")
	}

	res.Paid = true
	fmt.Println("Payment processed for reservation ID:", reservationID)
	return nil
}

// Main function
func main() {
	rentalSystem := NewRentalSystem()

	// Adding cars
	rentalSystem.AddCar(Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2020, LicensePlate: "ABC123", RentalPricePerDay: 50, IsAvailable: true})
	rentalSystem.AddCar(Car{ID: 2, Make: "Honda", Model: "Civic", Year: 2021, LicensePlate: "XYZ789", RentalPricePerDay: 60, IsAvailable: true})

	// Searching cars
	fmt.Println("Available Cars:", rentalSystem.SearchCars("Toyota", 100))

	// Creating reservation
	customer := Customer{Name: "John Doe", ContactDetails: "john.doe@example.com", DriversLicense: "D123456"}
	reservation, err := rentalSystem.CreateReservation(customer, 1, "2025-03-29", "2025-03-30")
	if err == nil {
		fmt.Println("Reservation created:", *reservation)
	}

	// Processing payment
	if err := rentalSystem.ProcessPayment(reservation.ID); err == nil {
		fmt.Println("Payment successful")
	}

	// Modifying reservation
	if err := rentalSystem.ModifyReservation(reservation.ID, "2025-04-01", "2025-04-02"); err == nil {
		fmt.Println("Reservation modified successfully")
	}

	// Canceling reservation
	if err := rentalSystem.CancelReservation(reservation.ID); err == nil {
		fmt.Println("Reservation canceled.")
	}
}
