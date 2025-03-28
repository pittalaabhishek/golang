package main

import (
	"errors"
	"fmt"
	"sync"
)

// Struct to represent a Car
type Car struct {
	ID              int
	Make            string
	Model           string
	Year            int
	LicensePlate    string
	RentalPricePerDay float64
	IsAvailable     bool
}

// Struct to represent a Customer
type Customer struct {
	Name            string
	ContactDetails  string
	DriversLicense  string
}

// Struct to represent a Reservation
type Reservation struct {
	ID         int
	Customer   Customer
	CarID      int
	StartDate  string
	EndDate    string
	TotalPrice float64
}

// RentalSystem struct to manage the entire system
type RentalSystem struct {
	Cars          []Car
	Reservations  []Reservation
	mu            sync.Mutex
	reservationID int
}

// AddCar adds a new car to the system
func (rs *RentalSystem) AddCar(car Car) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.Cars = append(rs.Cars, car)
}

// SearchCars searches for cars based on criteria
func (rs *RentalSystem) SearchCars(carType string, maxPrice float64) []Car {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var results []Car
	for _, car := range rs.Cars {
		if car.Make == carType && car.RentalPricePerDay <= maxPrice && car.IsAvailable {
			results = append(results, car)
		}
	}
	return results
}

// CreateReservation creates a new reservation
func (rs *RentalSystem) CreateReservation(customer Customer, carID int, startDate, endDate string) (Reservation, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	for i, car := range rs.Cars {
		if car.ID == carID && car.IsAvailable {
			rs.reservationID++
			totalPrice := car.RentalPricePerDay // Assuming only one day for simplicity
			reservation := Reservation{
				ID:         rs.reservationID,
				Customer:   customer,
				CarID:      carID,
				StartDate:  startDate,
				EndDate:    endDate,
				TotalPrice: totalPrice,
			}
			rs.Reservations = append(rs.Reservations, reservation)
			rs.Cars[i].IsAvailable = false
			return reservation, nil
		}
	}
	return Reservation{}, errors.New("car not available")
}

// CancelReservation cancels an existing reservation
func (rs *RentalSystem) CancelReservation(reservationID int) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	for i, res := range rs.Reservations {
		if res.ID == reservationID {
			// Make car available again
			for j, car := range rs.Cars {
				if car.ID == res.CarID {
					rs.Cars[j].IsAvailable = true
					break
				}
			}
			// Remove reservation
			rs.Reservations = append(rs.Reservations[:i], rs.Reservations[i+1:]...)
			return nil
		}
	}
	return errors.New("reservation not found")
}

// Main function
func main() {
	rentalSystem := &RentalSystem{}

	// Adding cars to the system
	rentalSystem.AddCar(Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2020, LicensePlate: "ABC123", RentalPricePerDay: 50, IsAvailable: true})
	rentalSystem.AddCar(Car{ID: 2, Make: "Honda", Model: "Civic", Year: 2021, LicensePlate: "XYZ789", RentalPricePerDay: 60, IsAvailable: true})

	// Searching for available cars
	availableCars := rentalSystem.SearchCars("Toyota", 100)
	fmt.Println("Available Cars:", availableCars)

	// Creating a reservation
	customer := Customer{Name: "John Doe", ContactDetails: "john.doe@example.com", DriversLicense: "D123456"}
	reservation, err := rentalSystem.CreateReservation(customer, 1, "2025-03-29", "2025-03-30")
	if err != nil {
		fmt.Println("Error creating reservation:", err)
	} else {
		fmt.Println("Reservation created:", reservation)
	}

	// Canceling a reservation
	err = rentalSystem.CancelReservation(reservation.ID)
	if err != nil {
		fmt.Println("Error canceling reservation:", err)
	} else {
		fmt.Println("Reservation canceled.")
	}
}
