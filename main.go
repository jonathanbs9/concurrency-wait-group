package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Wait Group
func main() {
	// Creamos el waitgroup. Se pueden agrupar una cant de routinas go. Le decimos cuantas. La vamos agregando
	// Cuando termina una, el wg queda listo para volver a recibir u otro grupo. (puede ser 1, 2, 100, no tiene limite)
	var (
		wg       sync.WaitGroup
		budgets  []Budget
		ors      []OfficeRevenue
		users    []User
		response Response
	)
	// Creamos un channel
	done := make(chan bool)
	errChan := make(chan error)

	// Le asigno 1 goRoutine
	wg.Add(1)
	// Defino la goRoutine
	go func() {
		budgets = getBudgets()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		officeRevenues, err := getOfficeRevenues()
		ors = officeRevenues
		if err != nil {
			go func() {
				errChan <- err
			}()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		users = getUsers()
		wg.Done()
	}()

	// Para poder construir la respuesta, necesitamos si o si que las 3 goRoutines finalicen
	// Definimos
	go func() {
		// Espera que todas terminen
		wg.Wait()
		// Le enviamos true al channel. (Le avisamos que ya terminó)
		done <- true //o close(done)
	}()

	select {
	case <-done:
		response = Response{
			Users:         users,
			OfficeRevenue: ors,
			Budget:        budgets,
		}
		close(done)
		// Si sucede un error en el errCh => 2 opciones.
		// 1- Hacer un break
		// 2-
	case err := <-errChan:
		fmt.Println(err.Error())
		close(errChan)
		break
	}

	fmt.Println(response)
}

func getOfficeRevenues() ([]OfficeRevenue, error) {
	var err error

	if time.Now().Unix()%2 == 0 {
		return nil, errors.New("Algun error en el código")
	}

	return []OfficeRevenue{
		{City: "Mar del Plata", Revenue: 9000},
		{City: "Rosario", Revenue: 8},
		{City: "San Juan", Revenue: 152},
		{City: "Santa Fe", Revenue: 91},
		{City: "San Luis", Revenue: 15125},
		{City: "Cordoba", Revenue: 440},
		{City: "Salta", Revenue: 900},
	}, err
}

func getUsers() []User {
	return []User{
		{ID: 1, Name: "Jonathan"},
		{ID: 2, Name: "Andrés"},
		{ID: 3, Name: "Gustavo"},
		{ID: 4, Name: "Pedro"},
		{ID: 3, Name: "Jesús"},
	}
}

func getBudgets() []Budget {
	return []Budget{
		{City: "Tierra del Fuego", EstimatedBudget: 150000, Variation: 1000},
		{City: "Santa Cruz", EstimatedBudget: 300000, Variation: 15000},
		{City: "Rawson", EstimatedBudget: 202000, Variation: 10000},
		{City: "Santa Rosa", EstimatedBudget: 10000, Variation: 500},
		{City: "Comodoro Rivadavia", EstimatedBudget: 5500, Variation: 5000},
		{City: "Bariloche", EstimatedBudget: 3000, Variation: 200},
	}
}

// Response struct =>
type Response struct {
	Users         []User
	OfficeRevenue []OfficeRevenue
	Budget        []Budget
}

// User struct =>
type User struct {
	ID   int
	Name string
}

// OfficeRevenue struct =>
type OfficeRevenue struct {
	City    string
	Revenue int
}

// Budget struct =>
type Budget struct {
	City            string
	EstimatedBudget int
	Variation       int
}
