// pool di risorse equivalenti senza guardie logiche
package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXPROC = 100 //massimo numero di processi
const MAXRES = 5    //massimo numero di risorse nel pool

var richiesta = make(chan int)
var rilascio = make(chan int)
var risorsa [MAXPROC]chan int
var done = make(chan int)
var termina = make(chan int)

func client(i int) {
	richiesta <- i
	r := <-risorsa[i]
	fmt.Printf("\n [client %d] uso della risorsa %d\n", i, r)
	timeout := rand.Intn(3)
	time.Sleep(time.Duration(timeout) * time.Second)
	rilascio <- r
	done <- i //comunico al main la terminazione
}

func server(nris int, nproc int) {

	var disponibili int = nris
	var res, p, i int
	var libera [MAXRES]bool
	var sospesi [MAXPROC]bool
	var nsosp int = 0

	for i := 0; i < nris; i++ {
		libera[i] = true
	}
	for i := 0; i < nproc; i++ {
		sospesi[i] = false
	}

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("nuovo ciclo server")
		select {
		case res = <-rilascio:
			if nsosp == 0 {
				disponibili++
				libera[res] = true
				fmt.Printf("[server]  restituita risorsa: %d  \n", res)
			} else {
				for i = 0; i < nproc && !sospesi[i]; i++ {
				}
				sospesi[i] = false
				nsosp--
				risorsa[i] <- res
			}
		case p = <-richiesta:
			if disponibili > 0 { //allocazione della risorsa
				for i = 0; i < nris && !libera[i]; i++ {
				}
				libera[i] = false
				disponibili--
				risorsa[p] <- i
				fmt.Printf("[server]  allocata risorsa %d a cliente %d \n", i, p)
			} else { // attesa
				nsosp++
				sospesi[p] = true
				fmt.Printf("[server]  il cliente %d attende..\n", i)
			}
		case <-termina: // quando tutti i processi clienti hanno finito
			fmt.Println("FINE !!!!!!")
			done <- 1
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

	//inizializzazione canali
	for i := 0; i < cli; i++ {
		risorsa[i] = make(chan int)
	}

	for i := 0; i < cli; i++ {
		go client(i)
	}
	go server(res, cli)

	//attesa della terminazione dei clienti:
	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1 //terminazione server
	<-done
}
