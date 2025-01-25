package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ============================================================
//                     COSTANTI / PARAMETRI
// ============================================================
const (
	MAXBUFFER    = 100  // Dimensione massima dei buffer per i canali
	MAX_CLIENTS  = 20   // Numero massimo di Client/Worker gestibili

	TYPE_A    = 0       // Primo tipo di risorsa
	TYPE_B    = 1       // Secondo tipo di risorsa
	TYPE_MIX  = 2       // Tipo "misto"

	MAX_A     = 4000    // Capacità massima risorsa di tipo A
	MAX_B     = 3000    // Capacità massima risorsa di tipo B

	LOT_A     = 700     // Lotto per la risorsa di tipo A
	LOT_B     = 300     // Lotto per la risorsa di tipo B
	LOT_MIX   = 500     // Lotto per la risorsa "mista"
)

// ============================================================
//                     STRUTTURA DATI
// ============================================================
// Request rappresenta la richiesta che un Client/Worker o un Supplier
// può fare al sistema. `ack` è il canale di risposta (acknowledgment).
type Request struct {
	id   int         // Identificativo di chi fa la richiesta
	tipo int         // Indica il tipo di risorsa coinvolta
	ack  chan int    // Canale di ack per segnalare eventi completati
}

// ============================================================
//                       CANALI
// ============================================================
// requestChan: canali usati dai Client/Worker per richiedere risorse.
//   - requestChan[TYPE_A]
//   - requestChan[TYPE_B]
//   - requestChan[TYPE_MIX]
// E' un array di canali: ognuno gestisce le richieste di uno specifico tipo.
// Esempio: requestChan[0] = canale per TYPE_A, requestChan[1] = canale per TYPE_B, requestChan[2] = canale per TYPE_MIX.
//
// Uso:
//   r := Request{id: 123, tipo: TYPE_A, ack: make(chan int)}
//   requestChan[TYPE_A] <- r  // Invia la richiesta sul canale per TYPE_A.
//   ...

var requestChan [3]chan Request

// restockChan: canali usati dai Supplier per rifornire il magazzino.
//   - restockChan[TYPE_A]
//   - restockChan[TYPE_B]
var restockChan [2]chan Request

// endRequest e endRestock: canali per segnalare la conclusione
// di un’operazione di prelievo/rifornimento.
var endRequest = make(chan Request, MAXBUFFER)
var endRestock = make(chan Request)

// Canali per la terminazione dei processi.
var done = make(chan bool)
var stopWarehouse = make(chan bool)
var stopSupplier = make(chan bool)

// ============================================================
//                   FUNZIONI DI SUPPORTO
// ============================================================

// Gestisce una guardia logica: restituisce il canale c se la condizione b è vera, nil altrimenti.
// Utile per lo switch select condizionato.
func when(b bool, c chan Request) chan Request {
	if !b {
		return nil
	}
	return c
}

// Attende un tempo randomico (in secondi) entro un range [0, max).
func sleepRandTime(max int) {
	if max > 0 {
		time.Sleep(time.Duration(rand.Intn(max)+1) * time.Second)
	}
}

// Attende un tempo randomico (in secondi) entro un range [min, max).
func sleepRandTimeRange(min, max int) {
	if min >= 0 && max > 0 && min < max {
		time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Second)
	}
}

// Restituisce una stringa in base al tipo di risorsa; i case sono presenti nelle costanti.
func getResourceName(t int) string {
	switch t {
	case TYPE_A:
		return "tipo A"
	case TYPE_B:
		return "tipo B"
	case TYPE_MIX:
		return "tipo MISTO"
	default:
		return "sconosciuto"
	}
}

// Stampa di debug sullo stato di prelievi, code, ecc.
func debug(activeA, activeB, lenA, lenB, lenMix int) {
	fmt.Printf("Stato: prelievi in corso (A: %d, B: %d) | Code(A: %d, B: %d, Mix: %d)\n",
		activeA, activeB, lenA, lenB, lenMix)
}

// ============================================================
//                       GOROUTINE
// ============================================================

// client simula un generico “Worker” o “Addetto” che ciclicamente
// richiede e preleva risorse dal magazzino.
func client(id int) {
    // è possibile dichiarare parzialmente le struct
    //r := Request{id: id, ack: make(chan int)}
    tipo := -1
	r := Request{id: id, tipo: tipo, ack: make(chan int)}

	fmt.Printf("[CLIENT %d] Avviato\n", id)
	for i := 0; i < 5; i++ {
		// Scelta casuale del tipo di risorsa (TYPE_A, TYPE_B o TYPE_MIX).
		tipoRand := rand.Intn(100)
		if tipoRand >= 80 {
			r.tipo = TYPE_MIX
		} else {
			r.tipo = tipoRand % 2 // 0 o 1
		}

		fmt.Printf("[CLIENT %d] Richiedo risorsa %s\n", id, strings.ToUpper(getResourceName(r.tipo)))
		requestChan[r.tipo] <- r       // invio richiesta
		<-r.ack                        // attendo ack di inizio operazione

		fmt.Printf("[CLIENT %d] Sto prelevando risorsa %s...\n", id, strings.ToUpper(getResourceName(r.tipo)))
		sleepRandTime(3)              // simulazione prelievo

		endRequest <- r               // segnalazione di fine prelievo
		<-r.ack

	}

    // > Inviamo un segnale sul canale `done` per notificare che questo client
    //   ha terminato i suoi prelievi. In altre parole, una volta completate
    //   tutte le iterazioni (il for), il client non servira' piu' e segnala
    //   al main (o a chi di competenza) che il suo lavoro e' concluso.
	done <- true
	fmt.Printf("[CLIENT %d] Termino\n", id)
}

// supplier simula un generico fornitore che ciclicamente rifornisce
// il magazzino di un certo tipo di risorsa.
func supplier(resourceType int) {
    //Non è necessario l'id perché non serve a questo fine
	r := Request{tipo: resourceType, ack: make(chan int)}

	fmt.Printf("[SUPPLIER %s] Avviato\n", strings.ToUpper(getResourceName(resourceType)))
	for {
		sleepRandTimeRange(5, 10)

		fmt.Printf("[SUPPLIER %s] Voglio rifornire il magazzino\n", strings.ToUpper(getResourceName(resourceType)))
		restockChan[resourceType] <- r   // invio richiesta rifornimento
		<-r.ack                          // attendo ack di inizio rifornimento

		fmt.Printf("[SUPPLIER %s] Rifornimento in corso...\n", strings.ToUpper(getResourceName(resourceType)))
		sleepRandTimeRange(3, 5)        // simulazione rifornimento

		endRestock <- r                 // segnalo il termine
		<-r.ack                          // attendo che il magazzino completi l’operazione
		fmt.Printf("[SUPPLIER %s] Rifornimento completato\n", strings.ToUpper(getResourceName(resourceType)))

        // > Questo select ci permette di verificare se e' stato richiesto
        //   di fermare il fornitore (ad esempio, quando il main vuole
        //   terminare la simulazione). Se non riceviamo nulla dal canale
        //   stopSupplier, continuiamo il ciclo. Se invece riceviamo un segnale,
        //   usciamo dal ciclo e terminiamo.
		select {
		case <-stopSupplier:
			fmt.Printf("[SUPPLIER %s] Termino\n", strings.ToUpper(getResourceName(resourceType)))
			done <- true
			return
		default:
			continue
		}
	}
}

// warehouse gestisce la logica di accesso e aggiornamento di un
// magazzino con due tipi di risorse (TYPE_A e TYPE_B) e un tipo “misto” (TYPE_MIX).
// Qui inserisci le condizioni di precedenza/priorità come da traccia d'esame.
func warehouse() {
	// Quante risorse disponibili inizialmente (puoi cambiare la logica di init)
	resources := [2]int{MAX_A, MAX_B}

	// Per sapere quanti prelievi e rifornimenti sono in corso 
  
    // > `activePrel` tiene traccia di quanti client stanno prelevando risorse contemporaneamente (diviso per tipo A e B). 
    //   `activeRestock` indica se è in corso un rifornimento per il tipo corrispondente.
	activePrel := [2]int{0, 0}
	activeRestock := [2]bool{false, false}

	fmt.Printf("[WAREHOUSE] Avviato. Stato iniziale: A: %d/%d, B: %d/%d\n", resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)

	for {
		select {
		//---------------------------------------------------
		//                  PRELIEVO (INIZIO)
		//---------------------------------------------------
		case req := <-when(
			// Condizioni per permettere prelievo TYPE_A
			( (LOT_A*(activePrel[TYPE_A]+1)) <= resources[TYPE_A] ) &&
				(!activeRestock[TYPE_A]) &&
				(len(requestChan[TYPE_MIX]) == 0),   // es. priorità a TYPE_MIX
			requestChan[TYPE_A]):

			activePrel[TYPE_A]++
			fmt.Printf("[WAREHOUSE] Client %d inizia prelievo di %d (tipo A)\n",
				req.id, LOT_A)
			req.ack <- 1 // sblocco client

		case req := <-when(
			// Condizioni per permettere prelievo TYPE_B
			( (LOT_B*(activePrel[TYPE_B]+1)) <= resources[TYPE_B] ) &&
				(!activeRestock[TYPE_B]) &&
				(len(requestChan[TYPE_MIX]) == 0 && len(requestChan[TYPE_A]) == 0),
			requestChan[TYPE_B]):

			activePrel[TYPE_B]++
			fmt.Printf("[WAREHOUSE] Client %d inizia prelievo di %d (tipo B)\n",
				req.id, LOT_B)
			req.ack <- 1 // sblocco client

		case req := <-when(
			// Condizioni per permettere prelievo TYPE_MIX (risorse A e B insieme)
			( (LOT_MIX*(activePrel[TYPE_A]+1)) <= resources[TYPE_A] &&
			  (LOT_MIX*(activePrel[TYPE_B]+1)) <= resources[TYPE_B] ) &&
				(!activeRestock[TYPE_A] && !activeRestock[TYPE_B]),
			requestChan[TYPE_MIX]):

			activePrel[TYPE_A]++
			activePrel[TYPE_B]++
			fmt.Printf("[WAREHOUSE] Client %d inizia prelievo misto di %d (A) e %d (B)\n",
				req.id, LOT_MIX, LOT_MIX)
			req.ack <- 1 // sblocco client

		//---------------------------------------------------
		//                  PRELIEVO (FINE)
		//---------------------------------------------------
		case req := <-endRequest:
			switch req.tipo {
			case TYPE_A:
				resources[TYPE_A] -= LOT_A
				activePrel[TYPE_A]--
			case TYPE_B:
				resources[TYPE_B] -= LOT_B
				activePrel[TYPE_B]--
			case TYPE_MIX:
				resources[TYPE_A] -= LOT_MIX
				resources[TYPE_B] -= LOT_MIX
				activePrel[TYPE_A]--
				activePrel[TYPE_B]--
			default:
				fmt.Println("[WAREHOUSE] ERRORE: tipo risorsa non valido.")
			}
			fmt.Printf("[WAREHOUSE] Client %d ha terminato. Stato: A: %d/%d, B: %d/%d\n",
				req.id, resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)

		//---------------------------------------------------
		//                 RIFORNIMENTO (INIZIO)
		//---------------------------------------------------
		case req := <-when(
			// Condizione per rifornire TYPE_A
			(activePrel[TYPE_A] == 0) &&
				(resources[TYPE_A] <= resources[TYPE_B] || len(restockChan[TYPE_B]) == 0),
			restockChan[TYPE_A]):
			activeRestock[TYPE_A] = true
			fmt.Printf("[WAREHOUSE] Inizio rifornimento A...\n")
			req.ack <- 1

		case req := <-when(
			// Condizione per rifornire TYPE_B
			(activePrel[TYPE_B] == 0) &&
				(resources[TYPE_B] < resources[TYPE_A] || len(restockChan[TYPE_A]) == 0),
			restockChan[TYPE_B]):
			activeRestock[TYPE_B] = true
			fmt.Printf("[WAREHOUSE] Inizio rifornimento B...\n")
			req.ack <- 1

		//---------------------------------------------------
		//                 RIFORNIMENTO (FINE)
		//---------------------------------------------------
		case req := <-endRestock:
			switch req.tipo {
			case TYPE_A:
				resources[TYPE_A] = MAX_A
				activeRestock[TYPE_A] = false
				fmt.Printf("[WAREHOUSE] Fine rifornimento A. A: %d/%d, B: %d/%d\n",
					resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)
				req.ack <- 1
			case TYPE_B:
				resources[TYPE_B] = MAX_B
				activeRestock[TYPE_B] = false
				fmt.Printf("[WAREHOUSE] Fine rifornimento B. A: %d/%d, B: %d/%d\n",
					resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)
				req.ack <- 1
			default:
				fmt.Println("[WAREHOUSE] ERRORE: tipo risorsa non valido.")
				req.ack <- -1
			}

		//---------------------------------------------------
		//                 TERMINAZIONE
		//---------------------------------------------------
		case <-stopWarehouse:
			fmt.Printf("[WAREHOUSE] Termino\n")
			done <- true
			return
		}
	}
}

// ============================================================
//                         MAIN
// ============================================================
func main() {
	fmt.Println("[MAIN] Inizio")
	rand.Seed(time.Now().UnixNano())

    //Aggiungi opzione di selezione da input utente
    // Esempio: l'utente potrebbe selezionare una "modalita'" di simulazione.
    // Puoi aggiungere piu' opzioni e logica in base al valore inserito.
    var mode int
    fmt.Print("Inserisci la modalita' di simulazione (1=base, 2=avanzata): ")
    fmt.Scanf("%d\n", &mode)
    fmt.Printf("Hai scelto la modalita' %d. Personalizza il comportamento in base al valore...\n", mode)

	// Esempio: specifica quanti client e quanti supplier avviare
	nClients := 5
	nSuppliers := 2

	fmt.Printf("[MAIN] Quanti Client vuoi avviare (max %d)? ", MAX_CLIENTS)
	fmt.Scanf("%d\n", &nClients)
	if nClients < 2 {
		fmt.Printf("[MAIN] Numero di client troppo basso. Uso valore di default: 4.\n")
		nClients = 4
	}

	// Inizializza i canali principali
	for i := 0; i < len(requestChan); i++ {
		requestChan[i] = make(chan Request, MAXBUFFER)
	}
	for i := 0; i < len(restockChan); i++ {
		restockChan[i] = make(chan Request, MAXBUFFER)
	}

	// Avvio goroutine
	go warehouse()                  // gestore delle risorse

	// Avvio supplier (uno per TYPE_A e uno per TYPE_B, se nSuppliers=2)
	for i := 0; i < nSuppliers; i++ {
		go supplier(i)              // i corrisponde a TYPE_A o TYPE_B
	}

	// Avvio client
	for i := 0; i < nClients; i++ {
		go client(i)
	}

    // > Il canale `done` viene usato come segnale di completamento da parte
    //   di client e supplier. In sostanza, quando un client finisce, invia `true`
    //   sul canale `done`. Qui nel main attendiamo quei segnali per sapere quando
    //   tutti i processi hanno concluso il loro lavoro, e poi possiamo chiudere
    //   l'intero programma in modo ordinato.

	// Attendo la terminazione dei client
	for i := 0; i < nClients; i++ {
		<-done
	}

	// Segnalo ai supplier di terminare (lo faccio tante volte quanto nSuppliers)
	for i := 0; i < nSuppliers; i++ {
		stopSupplier <- true
	}
	// Attendo la terminazione di tutti i supplier
	for i := 0; i < nSuppliers; i++ {
		<-done
	}

	// Segnalo al warehouse di terminare
	stopWarehouse <- true
	<-done

	fmt.Println("[MAIN] Fine")
}

// ============================================================
//                SEMANTICA DI GO (BASICS)
// ============================================================
// Di seguito, in breve, alcuni concetti chiave:
//
// 1. Dichiarazione di variabili:
//    - Forme più comuni:
//      var x int           // dichiaro x di tipo int
//      x := 10             // dichiarazione e assegnazione implicita del tipo
//      var y, z = 10, 20   // dichiarazione multipla
//
// 2. Canali:
//    - Dichiarazione e creazione:
//      var c chan int      // dichiaro un canale di int, non ancora creato
//      c = make(chan int)  // creazione del canale
//    - Invio (send):
//      c <- valore
//    - Ricezione (receive):
//      <- c         // se non mi interessa il contenuto, tipicamente gli ack
//      v := <- c  // se mi interessa il contenuto 
//      oppure: v, ok := <- c   // dove ok indica se il canale è ancora aperto
//
// 3. Cicli (for):
//    - Go non ha il while, si usa for:
//      for i := 0; i < 10; i++ {
//          ...
//      }
//    - for condizionale:
//      for condition {
//          ...
//      }
//    - for infinito:
//      for {
//          ...
//      }
//
// Questi concetti base coprono l'uso di variabili, canali e cicli in Go.
