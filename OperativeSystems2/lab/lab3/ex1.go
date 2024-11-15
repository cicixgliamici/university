//  !! GIVEN EXAMPLE FROM UNIVERSITY - ONLY SLIGHT CHANGES !!
/* Purpose: This program simulates a pool of equivalent resources managed by a server. 
            Clients (goroutines) request resources, use them, and release them. 
	    The server handles allocation and release of resources.
   Channels: Used for communication between the server and clients, avoiding explicit locking mechanisms like mutexes.
   Concurrency: Clients and the server operate concurrently, using goroutines to manage parallelism.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXPROC = 100 
const MAXRES = 5    

/* Channels are a powerful concurrency primitive used for communication between goroutines. 
   They allow one goroutine to send data to another, enabling synchronization and data sharing in a safe and efficient way. 
   Channels are part of Go's CSP (Communicating Sequential Processes) model.
   Synchronization: Channels automatically synchronize goroutines, preventing race conditions.
   Data Sharing: They allow safe and structured data sharing between goroutines.
   Concurrency Simplification: Channels reduce the complexity of managing locks or shared memory.
*/
var richiesta = make(chan int)        //Channel to request a resource
var rilascio = make(chan int)         //Channel to release a resource
var risorsa [MAXPROC]chan int         //Array of channels for each client to receive resource allocation
var done = make(chan int)             //Channel to signal when a client finishes
var termina = make(chan int)          //Channel to signal server termination

//Function executed by each client
func client(i int) {
	richiesta <- i                      //Client requests a resource by sending its ID
	r := <-risorsa[i]                   //Wait for the server to allocate a resource
	fmt.Printf("\n [client %d] using resource %d\n", i, r)
	timeout := rand.Intn(3)             //Simulate resource usage time (0-2 seconds)
	time.Sleep(time.Duration(timeout) * time.Second)
	rilascio <- r                       //Release the resource after usage
	done <- i                           //Notify main that the client has finished
}

// Server function to manage resource allocation and release
func server(nris int, nproc int) {
	var disponibili int = nris           
	var res, p, i int
	var libera [MAXRES]bool              //Tracks whether each resource is free
	var sospesi [MAXPROC]bool            //Tracks whether each client is waiting for a resource
	var nsosp int = 0                    //Number of clients waiting for resources
	for i := 0; i < nris; i++ {          //Initialize all resources as free
		libera[i] = true
	}
	for i := 0; i < nproc; i++ {         //Initialize all clients as not waiting
		sospesi[i] = false
	}
	for {
		time.Sleep(time.Second * 1)      
		fmt.Println("new server cycle")
		select {
		// Handle resource release
		case res = <-rilascio:
			if nsosp == 0 {                   //No clients are waiting
				disponibili++
				libera[res] = true        //Mark the resource as free
				fmt.Printf("[server] resource %d returned\n", res)
			} else {                          //Allocate resource to a waiting client
				for i = 0; i < nproc && !sospesi[i]; i++ {
				}
				sospesi[i] = false        //Remove client from waiting list
				nsosp--
				risorsa[i] <- res         //Assign resource to the client
			}
		// Handle resource requests
		case p = <-richiesta:
			if disponibili > 0 {             //Resources are available
				for i = 0; i < nris && !libera[i]; i++ {
				}
				libera[i] = false        //Mark the resource as allocated
				disponibili--
				risorsa[p] <- i          //Send the resource to the client
				fmt.Printf("[server] allocated resource %d to client %d\n", i, p)
			} else {                         //No resources available; client waits
				nsosp++
				sospesi[p] = true
				fmt.Printf("[server] client %d is waiting..\n", p)
			}
		// Handle server termination
		case <-termina:
			fmt.Println("FINISHED !!!")
			done <- 1                    // Notify the main function of termination
			return
		}
	}
}

func main() {
	var cli, res int
	rand.Seed(time.Now().Unix())
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)
	fmt.Printf("\n quante risorse (max %d)? ", MAXRES)
	fmt.Scanf("%d", &res)
	fmt.Println("risorse da gestire:", res)
	for i := 0; i < cli; i++ {                            //Initialize channels for each client
		risorsa[i] = make(chan int)
	}
	for i := 0; i < cli; i++ {                            //Launch client processes as goroutines
		go client(i)
	}
	go server(res, cli)                                   //Launch the server goroutine                
	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1                           		      // Signal the server to terminate
	<-done                                 		      // Wait for server termination confirmation
}
