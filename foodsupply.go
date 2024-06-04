package main

import (
    "fmt"
    "encoding/json"
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FoodSupply struct {
    contractapi.Contract
}

type Order struct {
    OrderID     string  `json:"OrederID"`
    Owner       string  `json:"Owner"`
    FoodID      string  `json:"FoodID"`
    ConsumerID  string  `json:"ConsumerID"`
    ManufactureID   string  `json:"ManufactureID"`
    ExpirationDate  string  `json:"ExpirationDate"`
    Value       int     `json:"Value"`
    Status      string  `json:Status"`
}

func (f *FoodSupply) OrderExists(ctx contractapi.TransactionContextInterface, orderID string) (bool, error) {
    orderJSON, err := ctx.GetStub().GetState(orderID)
    if err != nil {
        return false, fmt.Errorf("failed to read order: %v", err)
    }
    return orderJSON != nil, nil
}

func (f *FoodSupply) SetupOrder(ctx contractapi.TransactionContextInterface, orderID string, foodID string, value int) error {
    exists, err := f.OrderExists(ctx, orderID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("order %s already exists", orderID)
    }

    order := Order{
        OrderID: orderID,
        Owner: "Company",
        FoodID: foodID,
        Value: value,
        Status: "Order init",
    }

    orderJSON, err := json.Marshal(order)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(orderID, orderJSON)
}

func (f *FoodSupply) ManufactureProcessing(ctx contractapi.TransactionContextInterface, orderID string, manufactureID string) error {
    order, err := f.ReadOrder(ctx, orderID)
    if err != nil {
        return err
    }

    order.ManufactureID = manufactureID
    order.Owner = "Manufacture"
    order.Status = "Manufacture process"

    orderJSON, err := json.Marshal(order)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(orderID, orderJSON)
}

func (f *FoodSupply) DelieverToConsumer(ctx contractapi.TransactionContextInterface, orderID string, consumerID string, expirationDate string) error {
    order, err := f.ReadOrder(ctx, orderID)
    if err != nil {
        return err
    }

    order.ConsumerID = consumerID
    order.ExpirationDate = expirationDate
    order.Owner = "Consumer"
    order.Status = "Consumer received"

    orderJSON, err := json.Marshal(order)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(orderID, orderJSON)
}

func (f *FoodSupply) ReadOrder(ctx contractapi.TransactionContextInterface, orderID string) (*Order, error) {
    orderJSON, err := ctx.GetStub().GetState(orderID)
    if err != nil {
        return nil, fmt.Errorf("failed to read order: %v", err)
    }
    if orderJSON == nil {
        return nil, fmt.Errorf("order %s does not exist", orderID)
    }

    var order Order
    err = json.Unmarshal(orderJSON, &order)
    if err != nil {
        return nil, err
    }

    return &order, nil
}

func (f *FoodSupply) DeleteOrder(ctx contractapi.TransactionContextInterface, orderID string) error {
    exists, err := f.OrderExists(ctx, orderID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("order %s does not exist", orderID)
    }

    return ctx.GetStub().DelState(orderID)
}

func (f *FoodSupply) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {

    resultsIterator, err := ctx.GetStub().GetStateByRange("","")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var orders []*Order
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var order Order
        err = json.Unmarshal(queryResponse.Value, &order)
        if err != nil {
            return nil, err
        }
        orders = append(orders, &order)
    }

    return orders, nil
}


func main() {
    supplyChaincode, err := contractapi.NewChaincode(&FoodSupply{})
    if err != nil {
        log.Panicf("Error creating food-supply chaincode: %v", err)
    }

    if err := supplyChaincode.Start(); err != nil {
        log.Panicf("Error starting food-supply chaincode: %v", err)
    }
}
