package services

import (
	"car-rental-system/models"
	"errors"
	"fmt"
	"sync"
)

type RentalSystem struct {
	cars          map[int]*models.Car
	reservations  map[int]*models.Reservation
	mu            sync.Mutex
	reservationID int
}

func NewRentalSystem() *RentalSystem {
	return &RentalSystem{
		cars:         make(map[int]*models.Car),
		reservations: make(map[int]*models.Reservation),
	}
}

func (rs *RentalSystem) AddCar(car models.Car) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.cars[car.ID] = &car
}

func (rs *RentalSystem) SearchCars(make string, maxPrice float64) []models.Car {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var results []models.Car
	for _, car := range rs.cars {
		if car.Make == make && car.RentalPricePerDay <= maxPrice && car.IsAvailable {
			results = append(results, *car)
		}
	}
	return results
}

func (rs *RentalSystem) CreateReservation(customer models.Customer, carID int, startDate, endDate string) (*models.Reservation, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	car, exists := rs.cars[carID]
	if !exists || !car.IsAvailable {
		return nil, errors.New("car not available")
	}

	rs.reservationID++
	reservation := &models.Reservation{
		ID:         rs.reservationID,
		Customer:   customer,
		CarID:      carID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: car.RentalPricePerDay,
	}

	rs.reservations[rs.reservationID] = reservation
	car.IsAvailable = false

	return reservation, nil
}

func (rs *RentalSystem) ModifyReservation(reservationID int, newStartDate, newEndDate string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	res, exists := rs.reservations[reservationID]
	if !exists {
		return errors.New("reservation not found")
	}

	res.StartDate = newStartDate
	res.EndDate = newEndDate
	return nil
}

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