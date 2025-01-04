package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX_BUFFER = 100
const MAX_CLIENTS = 100

const SmallBottle = 0 		//0.5 liters, price 0.10
const LargeBottle = 1 		//1.5 liters, price 0.20

const CapacitySmall = 0.5 	//Capacity of a small bottle
const CapacityLarge = 1.5 	//Capacity of a large bottle

const TankCapacity = 50.0 	//Tank capacity (liters)

const SmallCoinsBox = 0 	//Box for 10-cent coins
const LargeCoinsBox = 1 	//Box for 20-cent coins

const MaxSmallCoins = 15	//Max number of 10-cent coins
const MaxLargeCoins = 20 	//Max number of 20-cent coins

//Array and buffered channels for clients
var start_request   [2]chan request
var end_request     =make(chan request, MAX_BUFFER)

//Buffered channels for the operator
var start_refill    =make(chan int, MAX_BUFFER)
var end_refill      =make(chan int, MAX_BUFFER)
var ack_operator    =make(chan int, MAX_BUFFER)

//Unbuffered termination channels
var done              =make(chan bool)
var terminate         =make(chan bool)
var terminateOperator =make(chan bool)

//Unbuffered = no synchronization
//Buffered   = synchronization

type request struct {
	index int
	kind  int
	ack   chan int
}

func whenRequest(b bool, c chan request) chan request {
	if !b {
		return nil
	}
	return c
}

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func sleepRandomTime(limit int) {
	if limit > 0 {
		time.Sleep(time.Duration(rand.Intn(limit)+1) * time.Second)
	}
}

//Goroutine for clients
func client(index int) {
	kind := rand.Intn(2)         //0 for small bottle, 1 for large bottle
	r := request{index, kind, make(chan int)}
	sleepRandomTime(2)           //Simulating payment
	if r.kind == SmallBottle {
		fmt.Printf("[client %d] requested a small bottle\n", index)
	} else {
		fmt.Printf("[client %d] requested a large bottle\n", index)
	}
	start_request[kind] <- r
	<-r.ack
	sleepRandomTime(3)           //Simulating bottle filling
	end_request <- r
	<-r.ack
	fmt.Printf("[client %d] finished filling my bottle, exiting!\n", index)
	done <- true
}

//Goroutine for operator
func operator() {
	var response int
	sleepRandomTime(4)
	for {
		start_refill <- 1
		response = <-ack_operator
		if response == -1 {
			fmt.Printf("[operator] exiting...\n")
			done <- true
			return
		}
		fmt.Printf("[operator] starting the refill process...\n")
		sleepRandomTime(3)                                                                //Time required for refilling
		end_refill <- 1
		<-ack_operator
		fmt.Printf("[operator] Refill complete, water station is operational again...\n")
		sleepRandomTime(5)
	}
}

//Server goroutine
func waterStation() {
	var currentWater = TankCapacity  //Current tank water level
	var smallCoinCount = 0           //Total 10-cent coins in the box
	var largeCoinCount = 0           //Total 20-cent coins in the box
	var busy = false
	var stop = false
	fmt.Printf("[waterStation] Water station is operational!\n")
	for {
		select {
			
		//Client requesting a small bottle
		case x := <-whenRequest(!busy && currentWater >= CapacitySmall && smallCoinCount < MaxSmallCoins &&
			               (largeCoinCount < MaxLargeCoins || 
				       (largeCoinCount == MaxLargeCoins && len(start_refill) == 0)), start_request[SmallBottle]):
			busy = true
			smallCoinCount++
			currentWater -= CapacitySmall
			fmt.Printf("[waterStation] Client %d started filling a bottle of type %d\n", x.index, x.kind)
			x.ack <- 1
			
		//Client requesting a large bottle
		case x := <-whenRequest(!busy && currentWater >= CapacityLarge && largeCoinCount < MaxLargeCoins && len(start_request[SmallBottle]) == 0 &&
			               (smallCoinCount < MaxSmallCoins || 
				       (smallCoinCount == MaxSmallCoins && len(start_refill) == 0)), start_request[LargeBottle]):
			busy = true
			largeCoinCount++
			currentWater -= CapacityLarge
			fmt.Printf("[waterStation] Client %d started filling a bottle of type %d\n", x.index, x.kind)
			x.ack <- 1
			
		//Operator refill process
		case <-when(!stop && !busy &&
			   ((smallCoinCount == MaxSmallCoins || largeCoinCount == MaxLargeCoins || currentWater == 0) ||
			   (len(start_request[SmallBottle])+len(start_request[LargeBottle]) == 0)), start_refill):
			busy = true
			currentWater = TankCapacity
			smallCoinCount = 0
			largeCoinCount = 0
			fmt.Printf("[waterStation] Operator started refilling the tank and emptying coin boxes\n")
			ack_operator <- 1
			
		//End of client request
		case x := <-end_request:
			busy = false
			x.ack <- 1
			
		//End of operator refill
		case <-end_refill:
			busy = false
			ack_operator <- 1
			
		//Termination of the operator
		case <-terminateOperator:
			stop = true
			fmt.Printf("[waterStation] All clients served, notifying operator to terminate\n")

		//Termination of the refill
		case <-when(stop, start_refill):
			ack_operator <- -1

	        //General termination
		case <-terminate:
			fmt.Printf("[waterStation] Shutting down!\n")
			done <- true
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2; i++ {
		start_request[i] = make(chan request, MAX_BUFFER)
	}
	for i := 0; i < MAX_CLIENTS; i++ {
		go client(i)
	}
	go operator()
	go waterStation()
	fmt.Printf("\n[MAIN] Water station is open.\n")
	for i := 0; i < MAX_CLIENTS; i++ {
		<-done
	}
	terminateOperator <- true
	<-done
	terminate <- true
	<-done
	fmt.Printf("\n[MAIN] Water station is closed.\n")
}
