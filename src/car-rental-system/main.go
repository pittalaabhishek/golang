package main

import (
	"car-rental-system/models"
	"car-rental-system/services"
	"fmt"
)

func main() {
	// Initialize Rental System
	rentalSystem := services.NewRentalSystem()

	// Adding cars
	rentalSystem.AddCar(models.Car{ID: 1, Make: "Toyota", Model: "Corolla", Year: 2020, LicensePlate: "ABC123", RentalPricePerDay: 50, IsAvailable: true})
	rentalSystem.AddCar(models.Car{ID: 2, Make: "Honda", Model: "Civic", Year: 2021, LicensePlate: "XYZ789", RentalPricePerDay: 60, IsAvailable: true})

	// Searching cars
	fmt.Println("Available Cars:", rentalSystem.SearchCars("Toyota", 100))

	// Creating reservation
	customer := models.Customer{Name: "John Doe", ContactDetails: "john.doe@example.com", DriversLicense: "D123456"}
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