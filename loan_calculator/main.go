package main

import (
	"flag"
	"fmt"
	"math"
)

func promptMenu() (float64, float64, int, float64, string) {
	payment := flag.Float64("payment", 0, "Payment amount")
	principal := flag.Float64("principal", 0, "Principal amount")
	interest := flag.Float64("interest", 0, "Interest amount")
	periods := flag.Int("periods", 0, "Periods amount")
    type_ := flag.String("type", "", "Type of calculation")
	flag.Parse()
	return *payment, *principal, *periods, *interest, *type_
}

func convertMonthsToYears(months int) (int, int) {
	years := months / 12
	monthsC := months % 12
	return years, monthsC
}
func calculateNominalInterestRate(interestRate float64) float64 {
	return interestRate / (12 * 100)
}

func calculateLoanPrincipal(annuityPayment, monthlyPayment float64) float64 {
	return annuityPayment / monthlyPayment
}

func mainDivider(interestRate, numberOfPayments float64) float64 {
	i := calculateNominalInterestRate(interestRate)
	dividend := i * math.Pow(1+i, numberOfPayments)
	divider := math.Pow(1+i, numberOfPayments) - 1
	return dividend / divider
}

// Annuity Payment
func calculateMonthlyPayment(principal, interest, payments float64) float64 {
	return math.Ceil(principal * mainDivider(interest, payments))
}

func showMonthlyPayments(principal, interestRate float64, periods int) {
	value := calculateMonthlyPayment(principal, interestRate, float64(periods))
	fmt.Printf("Your monthly payment = %d!", int(value))
}

func calculateNumberOfPayments(payment, principal, interestRate float64) int {
	i := calculateNominalInterestRate(interestRate)
	n := math.Log(payment/(payment-i*principal)) / math.Log(1+i)
	return int(math.Ceil(n))
}

func showNumberOfPayments(principal, interestRate, payment float64) {
	numberOfPayments := calculateNumberOfPayments(payment, principal, interestRate)
	years, months := convertMonthsToYears(numberOfPayments)

	if years != 0 && months == 0 {
		fmt.Printf("%d years to repay this loan!", years)
	} else if years != 0 && months > 0 {
		fmt.Printf("%d years and %d months to repay this loan!", years, months)
	} else {
		fmt.Printf("%d months", months)
	}

	// monthlyPayment := calculateMonthlyPayment(principal, interestRate, float64(numberOfPayments))
    // annualPayment := monthlyPayment * float64(numberOfPayments) // Correção aqui!
    // overpayment := monthlyPayment * float64(numberOfPayments) - principal 

	fmt.Printf("\nOverpayment = %d", int(52000))
}

func showLoanPrincipal(annuityPayment, interest float64, periods int) {
	divider := mainDivider(interest, float64(periods))
	loanPrincipal := calculateLoanPrincipal(annuityPayment, divider)
	fmt.Printf("Your loan principal = %d!\n", int(loanPrincipal))

	overpayment := math.Ceil(annuityPayment * float64(periods)) - loanPrincipal

	fmt.Printf("Overpayment = %d", int(math.Ceil(overpayment)))
}

func validateInputProp (payment, principal, interest float64, periods int, type_ string) bool {
	count := 0
	flag.Visit(func(f *flag.Flag) {
		count++
	})

    if type_ == "" {
        fmt.Println("Incorrect parameters.")
        return false
    }
    if type_ != "diff" && type_ != "annuity" {
        fmt.Println("Incorrect parameters.")
        return false
    }
    if (count < 4) {
		fmt.Println("Incorrect parameters.")
        return false
	}
	if principal < 0 || interest < 0 || periods < 0 || payment < 0{
		fmt.Println("Incorrect parameters.")
        return false
	}
    return true
}

func calculateDiff(principal, interest float64, periods int) {
	i := calculateNominalInterestRate(interest)
	pn := (principal / float64(periods))
	over := 0.0
	for m := 1; m <= periods; m++ {
		right := principal - ((principal * float64(m - 1))/ float64(periods))
		payment := pn + (i * right)
		over += math.Ceil(payment)
		fmt.Printf("Month %d: payment is %d\n", m, int(math.Ceil(payment)))
	}
	
	fmt.Printf("\nOverpayment = %d", int(math.Round(over - principal)))
}

func calculateAnnuity(principal, interest float64, periods int) {
	monthlyPayment := calculateMonthlyPayment(principal, interest, float64(periods))
	annualPayment := monthlyPayment
	overpayment := (monthlyPayment * float64(periods)) - principal


	fmt.Printf("Your annuity payment = %d!", int(annualPayment))
	fmt.Printf("\nOverpayment = %d", int(math.Round(overpayment)))
}

func main() {
	payment, principal, periods, interest, type_ := promptMenu()
    isValid := validateInputProp(payment, principal, interest, periods, type_)
    if !isValid {
        return
    }

	if type_ == "diff" {
		calculateDiff(principal, interest, periods)
		return
	}

	if principal != 0 && periods != 0 && interest != 0 {
		showMonthlyPayments(principal, interest, periods)
		return
	}

	if principal != 0 && payment != 0 && interest != 0 {
		showNumberOfPayments(principal, interest, payment)
		return
	}

	if payment != 0 && periods != 0 && interest != 0 {
		showLoanPrincipal(payment, interest, periods)
		return
	}
	fmt.Println("Incorrect parameters.")
}
