package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	INITIAL_MONEY = 550
	INITIAL_MILK = 540
	INITIAL_WATER = 400
	INITIAL_COFFEE = 120
	INITIAL_CHOCOLATE = 200
	INITIAL_CUP = 9
	CUP_USAGE = 1
)

var recipe map[string]Recipe

type Price = int

type Supply struct {
	Water, Milk, Coffee, Chocolate, Cups int
}

func (s Supply) CanPrepare(r Recipe) (bool, error) {
	if s.Cups < 1 {
		return false, errors.New("Sorry, not enough cups!")
	}
	
	if r.Water > s.Water {
		return false, errors.New("Sorry, not enough water!")
	}

	if r.Coffee > s.Coffee {
		return false, errors.New("Sorry, not enough coffee!")
	}

	if r.Chocolate > s.Chocolate {
		return false, errors.New("Sorry, not enough chocolate!")
	}

	if r.Milk > s.Milk {
		return false, errors.New("Sorry, not enough milk!")
	}

	return true, nil
}

func (s *Supply) Consume(r Recipe) (bool, error) {
	if _, err := s.CanPrepare(r); err != nil {
		return false, err
	}

	s.Water -= r.Water
	s.Milk -= r.Milk
	s.Coffee -= r.Coffee
	s.Chocolate -= r.Chocolate
	s.Cups -= CUP_USAGE
	return true, nil
}

func NewSupply(water, milk, coffee, cups, chocolate int) *Supply {
	return &Supply{
		Water: water,
		Milk: milk,
		Coffee: coffee,
		Cups: cups,
		Chocolate: chocolate,
	}
}

type Recipe struct {
	Water, Milk, Coffee, Chocolate int
}

func NewRecipe(water, milk, coffee, chocolate int) *Recipe {
	return &Recipe{
		Water: water,
		Milk: milk,
		Coffee: coffee,
		Chocolate: chocolate,
	}
}

type Menu struct {
	Items map[string]Price
}

func NewMenu() *Menu {
	return &Menu{
		Items: map[string]Price{
			"espresso": 4,
			"latte": 7,
			"cappuccino": 6,
			"hot_chocolate": 11,
		},
	}
}

func InitRecipes() map[string]Recipe {
	recipe := make(map[string]Recipe)

	recipe["espresso"] = *NewRecipe(250, 0, 16, 0)
	recipe["latte"] = *NewRecipe(350, 75, 20, 0)
	recipe["cappuccino"] = *NewRecipe(200, 100, 12, 0)
	recipe["hot_chocolate"] = *NewRecipe(0, 100, 0, 50)

	return recipe
}

type Machine struct {
	Money int
	Supply
	Menu
}

func NewMachine() *Machine {
	return &Machine{
		Money: INITIAL_MONEY,
		Supply: *NewSupply(INITIAL_WATER, INITIAL_MILK, INITIAL_COFFEE, INITIAL_CUP, INITIAL_CHOCOLATE),
		Menu: *NewMenu(),
	}
}

func (m Machine) ShowAvailableIngredients() {
	fmt.Printf("The coffee machine has:\n" + 
			   "%d ml of water\n" +
			   "%d ml of milk\n" +
			   "%d g of coffee beans\n" +
			   "%d g of chocolate\n" +
			   "%d disposable cups\n" +
			   "$%d of money\n\n",
			   m.Water, m.Milk,  m.Coffee, m.Chocolate,m.Cups, m.Money,
			)
}

func (m *Machine) Deposit(amount int) {
	m.Money += amount
}

func (m *Machine) Withdraw() int {
	amount := m.Money
	m.Money = 0
	return amount
}

func (m *Machine) BuyEspresso() {
	recipe := recipe["espresso"]

	success, err := m.Supply.Consume(recipe)
	if !success {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("I have enough resources, making you a coffee!")
	price := m.Menu.Items["espresso"]
	m.Deposit(price)
}

func (m *Machine) BuyLatte() {
	recipe := recipe["latte"]

	success, err := m.Supply.Consume(recipe)
	if !success {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("I have enough resources, making you a coffee!")
	price := m.Menu.Items["latte"]
	m.Deposit(price)
}

func (m *Machine) BuyCappuccino() {
	recipe := recipe["cappuccino"]

	success, err := m.Supply.Consume(recipe)
	if !success {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("I have enough resources, making you a coffee!")
	price := m.Menu.Items["cappuccino"]
	m.Deposit(price)
}

func (m *Machine) BuyHotChocolate() {
	recipe := recipe["hot_chocolate"]

	success, err := m.Supply.Consume(recipe)
	if !success {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("I have enough resources, making you a hot chocolate!")
	price := m.Menu.Items["hot_chocolate"]
	m.Deposit(price)
}

func (m *Machine) FillSupply(water, milk, coffee, cups, chocolate int) {
	m.Supply.Water += water
	m.Supply.Milk += milk
	m.Supply.Coffee += coffee
	m.Supply.Cups += cups
	m.Supply.Chocolate += chocolate
}

func handleBuy(machine *Machine) {
	actions := map[string]func(){
		"1": machine.BuyEspresso,
		"2": machine.BuyLatte,
		"3": machine.BuyCappuccino,
		"4": machine.BuyHotChocolate,
	}
	fmt.Println("What do you want to buy? 1 - espresso, 2 - latte, 3 - cappuccino 4 - hot chocolate, back - to main menu:")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return
	}
	actions[strconv.Itoa(int(num))]()
}

func handleFill(machine *Machine) {
	var water, milk, coffee, chocolate, cups int

	fmt.Println("Write how many ml of water you want to add: ")
	fmt.Scanf("%d", &water)
	fmt.Println("Write how many ml of milk you want to add:")
	fmt.Scanf("%d", &milk)
	fmt.Println("Write how many grams of coffee beans you want to add:")
	fmt.Scanf("%d", &coffee)
	fmt.Println("Write how many grams of chocolate you want to add:")
	fmt.Scanf("%d", &chocolate)
	fmt.Println("Write how many disposable cups you want to add: ")
	fmt.Scanf("%d", &cups)

	machine.FillSupply(water, milk, coffee, cups, chocolate)
	fmt.Println()
}

func handleWithdraw(machine *Machine) {
	amount := machine.Withdraw()
	fmt.Printf("I gave you $%d\n\n", amount)
}

func main() {
	recipe = InitRecipes()
	machine := NewMachine()

	var action string
	for {
		fmt.Println("Write action (buy, fill, take, remaining, exit):")
		fmt.Scan(&action)
	
		switch action {
			case "buy":
				handleBuy(machine)
			case "fill":
				handleFill(machine)
			case "take":
				handleWithdraw(machine)
			case "remaining":
				machine.ShowAvailableIngredients()
			case "exit":
				os.Exit(0)
		}
	}
}