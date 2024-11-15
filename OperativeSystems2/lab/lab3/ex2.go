//  !! GIVEN EXAMPLE FROM UNIVERSITY - ONLY SLIGHT CHANGES !!
/*
1. Goroutines:
   - Lightweight threads managed by the Go runtime. 
   - Used to implement concurrency in this program. The `client` and `server` functions are run as goroutines to execute concurrently.
 2. Channels:
   - Provide a way for goroutines to communicate and synchronize by passing data.
3. Select Statement:
   - Used to wait on multiple channel operations simultaneously.
   - The `select` in the `server` function monitors the following:
     - Resource release (`rilascio`).
     - Resource requests (`richiesta`), conditioned by availability (`when` function).
     - Termination signal (`termina`).
4. Logical Guards (`when` function):
   - Implements conditional channel enabling. If the condition (availability of resources) is false, the channel becomes `nil` and is ignored in the `select` statement.
5. Concurrency Synchronization:
   - Channels (`done` and `termina`) are used to synchronize the completion of goroutines and the termination of the server.
*/
package main

import (
	"fmt"
	"time"
)

const MAXPROC = 100
const MAXRES = 5

var richiesta = make(chan int)    //Channel for clients to request resources
var rilascio = make(chan int)     //Channel for clients to release resources
var risorsa [MAXPROC]chan int     //Per-client channels to receive allocated resources
var done = make(chan int)         //Channel to notify the main thread when a client is done
var termina = make(chan int)      //Channel to signal the server to terminate

func when(b bool, c chan int) chan int {
	if !b {
		return nil                   //Disables the channel if the condition is false
	}
	return c                       //Returns the original channel if the condition is true
}

/* The <- operator is used for channel communication.
   Sending: channel <- value sends a value into a channel.
   Receiving: value := <-channel receives a value from a channel.
   Channels synchronize goroutines and can transfer data between them.
*/
func client(i int) {
	richiesta <- i                //Request a resource by sending the client ID to the server
	r := <-risorsa[i]             //Wait to receive the allocated resource from the server
	fmt.Printf("\n [client %d] uso della risorsa %d\n", i, r)
	rilascio <- r                 //Release the resource back to the server
	done <- i                     //Notify the main thread of completion
}

func server(nris int) {
	var disponibili int = nris     //Number of available resources
	var res, p, i int              //Temporary variables for resource index, client ID, etc.
	var libera [MAXRES]bool        //Tracks whether each resource is available
	for i := 0; i < nris; i++ {
		libera[i] = true             //Initialize all resources as available
	}
	for {
		time.Sleep(time.Second * 1) 
		fmt.Println("nuovo ciclo server")
		select {
    		case res = <-rilascio:                                             //Resource release
    			disponibili++
    			libera[res] = true
    			fmt.Printf("[server]  restituita risorsa: %d  \n", res)
    		case p = <-when(disponibili > 0, richiesta):                       //Handle resource request if available
    			for i = 0; i < nris && !libera[i]; i++ {
    			}
    			libera[i] = false
    			disponibili--
    			risorsa[p] <- i                                                 //Allocate resource to client
    			fmt.Printf("[server]  allocata risorsa %d a cliente %d \n", i, p)
    		case <-termina:                                                   //Terminate when signaled
    			fmt.Println("FINE !!!!!!")
    			done <- 1                                                       //Notify main thread of server completion
    			return
		}
	}
}


func main() {
	var cli, res int
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)
	fmt.Printf("\n quante risorse (max %d)? ", MAXRES)
	fmt.Scanf("%d", &res)
	fmt.Println("risorse da gestire:", res)
  // Initialize client-specific channels
	for i := 0; i < MAXPROC; i++ {
		risorsa[i] = make(chan int)
	}
	// Start client and server goroutines
	for i := 0; i < cli; i++ {
		go client(i)
	}
	go server(res)
	// Wait for all clients to complete
	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1 // Signal server to terminate
	<-done
}
