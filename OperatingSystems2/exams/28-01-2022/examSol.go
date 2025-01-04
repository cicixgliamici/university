package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Constants for the simulation
const (
	MAXWORKERS    =25   //Maximum number of workers (AR)
	MAXCYCLES     =3    //Maximum number of work cycles per worker
	SSM           =10   //Shelf capacity for surgical masks
	SFFP2         =10   //Shelf capacity for FFP2 masks
	
        BSM           =3    //Batch size for surgical masks
	BFFP2         =4    //Batch size for FFP2 masks
	BMM           =2    //Batch size for mixed masks

	T_MIX         =0    //Mixed mask type
	T_FFP2        =1    //FFP2 mask type
	T_CHIR        =2    //Surgical mask type

	S_FFP2        =0    //Supplier for FFP2 masks
	S_SM          =1    //Supplier for surgical masks
)

//Channels for synchronization and communication
//General syntax:                             reader <- sender
//Waiting ack (no interest for the data):     <- sender
//Saving the information:                     variable (may be another channel) <- sender
var (
	doneTask            =make(chan bool)           //Signals task completion
	doneWarehouse       =make(chan bool)           //Signals warehouse shutdown
	closeWarehouse      =make(chan bool)           //Closes the warehouse

	//Channels for workers (AR)
	startwithdrawal     [3]chan request               //Start request channels by type
	endwithdrawal       =make(chan request, 100)      //End of pickup notifications

	//Channels for suppliers
	startDelivery       [2]chan request               //Start request channels by supplier
	endDelivery         =make(chan request, 100)      //End of supply notifications
)

//Struct to represent a request
type request struct {
	id   int         //ID of the worker or supplier
	tipo int         //Type of request (mask type or supplier type)
	ack  chan int    //Acknowledgment channel
}

//Simulates random sleep time
func sleep(t int) {
	time.Sleep(time.Duration(rand.Intn(t)+1)*time.Second)
}

//Conditional channel helper
// Implements a logical guard for channels.
// If the provided condition is true, the function returns the input channel `c`.
// Otherwise, it returns `nil`, effectively disabling the channel in the `select` statement.
// This mechanism is used to dynamically enable or disable specific channels
// in the `select` block based on runtime conditions.
// The `nil` channel is non-selectable, meaning it will not block or trigger in the `select` block.
func when(condition bool, c chan request) chan request {
	if !condition {
		return nil
	}
	return c
}

//Worker function
// Represents a worker responsible for withdrawing mask batches from the warehouse.
// The worker randomly selects a batch type (mixed, FFP2, or surgical masks) and attempts to withdraw it.
// If the requested type is unavailable, the worker waits for the replenishment.
// The function handles the worker's withdrawal process over multiple cycles until completion.
// Each withdrawal includes synchronization with the warehouse to ensure constraints are met.
func AR(id int) { 
	tipo := rand.Intn(3)                                                        
	r := request{id, tipo, make(chan int)}
	sleep(10)
	cycles := rand.Intn(MAXCYCLES) + 1
	for i := 0; i < cycles; i++ {                                               
		tipo := rand.Intn(3)                                                      
		r.tipo = tipo
		if tipo == T_MIX {
			fmt.Printf("[Worker %d] requesting a mixed batch\n", id)
		} else if tipo == T_FFP2 {
			fmt.Printf("[Worker %d] requesting an FFP2 batch\n", id)
		} else {
			fmt.Printf("[Worker %d] requesting a surgical mask batch\n", id)
		}
		startwithdrawal[tipo] <- r                                                  //Start a request and wait for the synchronization
		<-r.ack                                                                     //Waiting for the acknowledgement 
		sleep(10)                                                                   //Time to process the withdrawal
		endwithdrawal <- r                                                          //Notify end of withdrawal
		<-r.ack
	}
	fmt.Printf("[Worker %d] finished\n", id)
	doneTask <- true
}

//Supplier function
// Represents a supplier responsible for restocking the warehouse shelves.
// There are two suppliers: one for surgical masks and another for FFP2 masks.
// The supplier randomly initiates restocking cycles and replenishes the respective shelf fully when empty.
// Restocking can only begin if no workers are withdrawing from the same shelf.
// The function also handles termination signals to stop the supplier's activity.
func supplier(tipo int) { 
	r := request{tipo, tipo, make(chan int)}
	for {                                                                                
		sleep(5)
		fmt.Printf("[Supplier %d] requesting restock for mask type %d\n", tipo, tipo)
		startDelivery[tipo] <- r                                                            //No need for synchronization
		flag := <-r.ack                                                                     //But need to know if he have to stop
		sleep(20)                                                                           //Time to restock the shelf
		if flag == 0 {                                                                      //Supplier needs to terminate
			doneTask <- true
			fmt.Printf("[Supplier %d] finished!\n", tipo)
			return
		}
		endDelivery <- r
		<-r.ack
		fmt.Printf("[Supplier %d] finished restocking the shelf for mask type %d\n", tipo, tipo)
		sleep(3)
	}
}

//Warehouse function
// Simulates the warehouse as a shared resource for workers and suppliers.
// Handles synchronization of withdrawals and restocking while enforcing constraints:
// - Workers can only withdraw if the requested batch type is available.
// - Suppliers can only restock when their respective shelf is empty and not in use by workers.
// Implements priority rules:
// - Mixed batches have the highest priority.
// - FFP2 batches have priority over surgical batches.
// - The supplier for the shelf with the fewest masks has priority; in case of ties, the FFP2 supplier is prioritized.
// Manages the termination of operations when all workers and suppliers finish their tasks.
func warehouse() {               
	var surgicalMasks = SSM        //Number of surgical masks available on the shelf
	var ffp2Masks = SFFP2          //Number of FFP2 masks available on the shelf
	var suppliersInFFP2 = 0        //Number of suppliers restocking the FFP2 shelf
	var suppliersInSurgical = 0    //Number of suppliers restocking the surgical shelf
	var workersInFFP2 = 0          //Number of workers withdrawing from the FFP2 shelf
	var workersInSurgical = 0      //Number of workers withdrawing from the surgical shelf
	var end = false                //Becomes true when all workers have finished

	for {
		select {
		//It works like: if the condition (first argument in the when) it's true, look at the channel and see if there is something, like an ack or a valure
		case x := <-when(surgicalMasks >= BMM && ffp2Masks >= BMM && suppliersInSurgical == 0 && suppliersInFFP2 == 0, startwithdrawal[T_MIX]):
			workersInSurgical++
			workersInFFP2++
			surgicalMasks -= BMM
			ffp2Masks -= BMM
			fmt.Printf("[Warehouse] Worker %d begins to withdraw a mixed batch\n", x.id)
			x.ack <- 1

		case x := <-when(ffp2Masks >= BFFP2 && suppliersInFFP2 == 0 && len(startwithdrawal[T_MIX]) == 0, startwithdrawal[T_FFP2]):
			workersInFFP2++
			ffp2Masks -= BFFP2
			fmt.Printf("[Warehouse] Worker %d begins to withdraw an FFP2 batch\n", x.id)
			x.ack <- 1

		case x := <-when(surgicalMasks >= BSM && suppliersInSurgical == 0 && len(startwithdrawal[T_MIX]) == 0 && len(startwithdrawal[T_FFP2]) == 0, startwithdrawal[T_CHIR]):
			workersInSurgical++
			surgicalMasks -= BSM
			fmt.Printf("[Warehouse] Worker %d begins to withdraw a surgical mask batch\n", x.id)
			x.ack <- 1

		case x := <-endwithdrawal:
			if x.tipo == T_MIX {
				workersInSurgical--
				workersInFFP2--
			} else if x.tipo == T_FFP2 {
				workersInFFP2--
			} else {
				workersInSurgical--
			}
			fmt.Printf("[Warehouse] Worker %d has finished the withdrawal\n", x.id)
			x.ack <- 1

		case x := <-when(!end && surgicalMasks < SSM && suppliersInSurgical == 0 && workersInSurgical == 0 && ((surgicalMasks >= ffp2Masks && len(startDelivery[S_FFP2]) == 0) || surgicalMasks < ffp2Masks), startDelivery[S_SM]):
			surgicalMasks = SSM
			suppliersInSurgical++
			x.ack <- 1

		case x := <-when(!end && ffp2Masks < SFFP2 && suppliersInFFP2 == 0 && workersInFFP2 == 0 && ((surgicalMasks < ffp2Masks && len(startDelivery[S_SM]) == 0) || surgicalMasks >= ffp2Masks), startDelivery[S_FFP2]):
			ffp2Masks = SFFP2
			suppliersInFFP2++
			fmt.Printf("[Warehouse] Supplier %d has started restocking the shelf for type %d\n", x.id, x.tipo)
			x.ack <- 1

		case x := <-endDelivery:
			if x.tipo == S_FFP2 {
				suppliersInFFP2--
			} else {
				suppliersInSurgical--
			}
			fmt.Printf("[Warehouse] Supplier %d has finished restocking the shelf for type %d\n", x.id, x.tipo)
			x.ack <- 1

		case <-doneTask:
			end = true
			fmt.Printf("[Warehouse] The warehouse is about to close...\n")

		case x := <-when(end, startDelivery[0]):
			x.ack <- 0
		case x := <-when(end, startDelivery[1]):
			x.ack <- 0

		case <-closeWarehouse:
			end = true
			fmt.Printf("[Warehouse] The warehouse is CLOSED!\n")
			doneWarehouse <- true
			return
		}
	}
}

//Orchestrates the simulation by launching concurrent processes for the warehouse, workers, and suppliers.
//Initializes the necessary channels for communication and synchronization.
//Creates a random number of workers and starts their withdrawal tasks.
//Waits for all workers and suppliers to complete their operations before signaling the warehouse to terminate.
//Ensures proper termination of the entire system.
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2; i++ {
		startDelivery[i] = make(chan request, 100)
	}
	for i := 0; i < 3; i++ {
		startwithdrawal[i] = make(chan request, 100)
	}
	numWorkers := rand.Intn(MAXWORKERS) + 2              //Ensure at least 2 workers
	fmt.Printf("Number of workers: %d\n", numWorkers)

	// Launch warehouse and suppliers
	go warehouse()
	go supplier(S_SM)
	go supplier(S_FFP2)

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		go AR(i)
	}

	// Wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		<-doneTask
	}
	fmt.Printf("[MAIN] All workers have finished!\n")
	doneTask <- true                                     //Notify the warehouse that all workers are finished

	// Wait for both suppliers
	for i := 0; i < 2; i++ {
		<-doneTask
	}
	closeWarehouse <- true                               //Command the warehouse to terminate
	<-doneWarehouse                                      //Wait for the warehouse to finish
}
