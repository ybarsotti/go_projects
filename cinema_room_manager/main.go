package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type Cinema struct {
	SeatsStruct [][]bool
	Rows, SeatsPerRow, CurrentIncome int
}

const (
	AVAILABLE_SEAT = "S"
	RESERVED_SEAT = "B"
	MAX_TICKET_PRICE_SEAT_THRESHOLD = 60
)

func (c *Cinema) initCinema(rows, seatsPerRow int) {
	c.SeatsStruct = make([][]bool, rows)
	c.Rows = rows
	c.SeatsPerRow = seatsPerRow

	for i := range c.SeatsStruct {
		c.SeatsStruct[i] = make([]bool, seatsPerRow)
	}
}

func MakeCinema(rows, seatsPerRow int) *Cinema { 
	cinema := new(Cinema)
	cinema.initCinema(rows, seatsPerRow)
	return cinema
}

func (c *Cinema) ShowSeats() {
	fmt.Println("Cinema:")

	// Header
	for i := 1; i <= c.SeatsPerRow; i++ {
		if i == 1 {
			fmt.Print("  " + strconv.Itoa(i) + " ")
		} else {
			fmt.Print(strconv.Itoa(i) + " ")
		}
	}
	fmt.Println()
	for indexRow, row := range c.SeatsStruct {
		fmt.Print(indexRow + 1, " ")

		for _, seat := range row {
			if seat {
				fmt.Print(RESERVED_SEAT, " ")
			} else {
				fmt.Print(AVAILABLE_SEAT, " ")
			}
		}

		fmt.Println()
	}
	fmt.Println()
}

func (c *Cinema) ReserveSeat() (bool, error) {
	var row, seat int
	for {
		row, seat = promptReserve()

		if row < 1 || seat < 1 || len(c.SeatsStruct) <= row - 1 || len(c.SeatsStruct[row-1]) <= seat - 1 {
			fmt.Printf("Wrong input!\n\n")
			continue
		}

		if c.SeatsStruct[row - 1][seat - 1] {
			fmt.Printf("That ticket has already been purchased!\n\n")
			continue
		} else {
			break
		}
	}
	
	c.SeatsStruct[row - 1][seat - 1] = true
	printTicketPrice(c, row)
	return true, nil
}

func (c Cinema) GetReservedCount() int {
	var reservedSeats int

	for _, row := range c.SeatsStruct {
		for _, seat := range row {
			if seat {
				reservedSeats++
			}
		}
	}

	return reservedSeats
}

func (c *Cinema) GetTotalSeats() int {
	return c.Rows * c.SeatsPerRow
}

func (c *Cinema) AddCurrentIncome(price int) {
	c.CurrentIncome += price
}

func (c *Cinema) CalculateTicketPrice(row int) int {
	var price int
	if c.GetTotalSeats() < MAX_TICKET_PRICE_SEAT_THRESHOLD {
		price := 10
		c.AddCurrentIncome(price)
		return price
	}

	frontHalf := math.Floor(float64(c.Rows) / float64(2))

	if row <= int(frontHalf) {
		price := 10
		c.AddCurrentIncome(price)
		return price
	}
	price = 8
	c.AddCurrentIncome(price)
	return price
}

func (c Cinema) CalculatePossibleTotalIncome() {
	totalSeats := c.GetTotalSeats()

	fmt.Print("Total income: ")

	if totalSeats < 60 {
		fmt.Printf("$%d\n\n", totalSeats * 10)
		return 
	} 

	half := float64(c.Rows) / float64(2)
	total := (math.Floor(half) * float64(c.SeatsPerRow) * 10) + (math.Ceil(half) * float64(c.SeatsPerRow) * 8)
	fmt.Printf("$%d\n\n", int(total))
}

func (c Cinema) GetBoughtPercentage() float32 {
	return float32(c.GetReservedCount()) / float32(c.GetTotalSeats()) * 100
}

func (c *Cinema) ShowStatistics() {
	fmt.Printf("Number of purchased tickets: %d\n", c.GetReservedCount())
	fmt.Printf("Percentage: %.2f%%\n", c.GetBoughtPercentage())
	fmt.Printf("Current income: $%d\n", c.CurrentIncome)
	c.CalculatePossibleTotalIncome()
}

func printTicketPrice(cinema *Cinema, row int) {
	price := cinema.CalculateTicketPrice(row)
	fmt.Printf("Ticket price: $%d\n\n", price)
}

func promptReserve() (row, seat int) {
	fmt.Println("Enter a row number:")
	fmt.Scanf("%d", &row)

	fmt.Println("Enter a seat number in that row:")
	fmt.Scanf("%d", &seat)
	fmt.Println()
	return row, seat
}


func setupCinema() *Cinema {
	var rows, seatPerRow int
	fmt.Println("Enter the number of rows:")
	fmt.Scanf("%d", &rows)

	fmt.Println("Enter the number of seats in each row:")
	fmt.Scanf("%d", &seatPerRow)
	fmt.Println()

	return MakeCinema(rows, seatPerRow)
}

func promptMenu() (option int) {
	fmt.Println("1. Show the seats")
	fmt.Println("2. Buy a ticket")
	fmt.Println("3. Statistics")
	fmt.Println("0. Exit")
	fmt.Scanf("%d", &option)
	return option
}

func main() {
	cinema := setupCinema()

	for {
		option := promptMenu()
		switch option {
		case 1:
			cinema.ShowSeats()
		case 2:
			cinema.ReserveSeat()
		case 3:
			cinema.ShowStatistics()
		default:
			fmt.Print("Bye bye!")
			os.Exit(0)
		}
	}
}

