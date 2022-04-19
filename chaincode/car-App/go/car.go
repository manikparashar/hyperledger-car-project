package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//Car Manfaturing Smart Contract
type CarManufactureSmartContract struct {
	contractapi.Contract
}

//Car object
type Car struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:owner`
	State  string `json:"state"`
}

// type State struct {
// 	CurrentState string `json:"currentState"`
// 	Delivered    bool   `json:"delivered"`
// 	ReadyForSale bool   `json:"readyForSale"`
// }

// InitLedger adds a base set of assets to the ledger
func (s *CarManufactureSmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	car := []Car{
		{ID: "car1", Name: "Mercedes", Model: "M123", Colour: "Black", Owner: "Manufaturer", State: "Under Production"},
		{ID: "car2", Name: "Santro", Model: "S123", Colour: "Metallic", Owner: "Manufaturer", State: "Under Production "},
		{ID: "car3", Name: "Balleno", Model: "B123", Colour: "Blue", Owner: "Manufaturer", State: "Under Production"},
		{ID: "car4", Name: "Dzire", Model: "D123", Colour: "White", Owner: "Manufaturer", State: "Under Production"},
		{ID: "car5", Name: "Maruti", Model: "Ma123", Colour: "Green", Owner: "Manufaturer", State: "Under Production"},
		{ID: "car6", Name: "BMW", Model: "BM123", Colour: "White", Owner: "Manufaturer", State: "Under Production"},
	}

	for _, c := range car {
		carJSON, err := json.Marshal(c)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(c.ID, carJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// This function Manufatures a car and puts its state into db
func (c *CarManufactureSmartContract) CreateCar(ctx contractapi.TransactionContextInterface,
	id, name, model, colour string) error {

	// checking if the car already exists in the ledger
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read the data from world state %s", err)
	}
	if carJSON != nil {
		return fmt.Errorf("car with %s already exists", id)
	}

	if err != nil {
		return err
	}

	car := Car{
		ID:     id,
		Name:   name,
		Model:  model,
		Colour: colour,
		Owner:  "Manufaturer",
		State:  "CREATED",
	}

	carbytes, err := json.Marshal(car)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, carbytes)
}

// This function gets car by id
func (c *CarManufactureSmartContract) GetCarById(ctx contractapi.TransactionContextInterface, id string) (*Car, error) {
	carJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from the world state Err:%s", err)
	}

	if carJson == nil {
		return nil, fmt.Errorf("the car with id %s does not exist", id)
	}
	var car *Car
	err = json.Unmarshal(carJson, &car)
	if err != nil {
		return nil, err
	}
	return car, nil
}

// This function delivers the car to dealer and also changes the state of the Car
func (c *CarManufactureSmartContract) DeliverToDealer(ctx contractapi.TransactionContextInterface, id string) error {

	car, err := c.GetCarById(ctx, id)
	if err != nil {
		return err
	}
	car.Owner = "Dealer"
	car.State = "READY_FOR_SALE"

	carJson, err := json.Marshal(car)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, carJson)
}

// list all the cars dealer has
func (c *CarManufactureSmartContract) ListCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {

	carIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer carIterator.Close()

	var cars []*Car
	for carIterator.HasNext() {
		carResponse, err := carIterator.Next()
		if err != nil {
			return nil, err
		}

		var car *Car
		err = json.Unmarshal(carResponse.Value, &car)
		if err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}
	return cars, nil
}

// This function sells the car
func (c *CarManufactureSmartContract) SellCar(ctx contractapi.TransactionContextInterface, id string, newBuyer string) error {

	car, err := c.GetCarById(ctx, id)
	if err != nil {
		return err
	}
	if car.State == "READY_FOR_SALE" {
		car.Owner = newBuyer
		car.State = "Sold"
	}

	carJson, err := json.Marshal(car)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, carJson)
}

func main() {
	carSmartContract := new(CarManufactureSmartContract)
	c, err := contractapi.NewChaincode(carSmartContract)
	if err != nil {
		panic(err.Error())
	}

	if err := c.Start(); err != nil {
		panic("Error:" + err.Error())
	}
}
