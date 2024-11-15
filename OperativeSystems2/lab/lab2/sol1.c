//  !! GIVEN EXAMPLE FROM UNIVERSITY - ONLY SLIGHT CHANGES !!
/* Mutex: Protects shared access to the park's data. 
* Semaphore: Coordinates threads that need to wait for specific conditions (e.g., available spots or vehicles).
* Threading: Each visitor is represented by a thread that operates independently but is coordinated through synchronization mechanisms.
*/

#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>

#define NUM_THREADS 100
#define MaxP 20
#define MaxB 5
#define MaxM 8
#define MONOP 0
#define BICI 1

typedef struct {
	int posti_liberi;
	int bici_libere; //bici disponibili
	int mono_liberi; //monopattini disponibili
	int sospesi;
	sem_t S;
	pthread_mutex_t m;
} parco;

parco P;

void init(parco *p);
void entra(parco *P, int mezzo);
void esci(parco *p, int mezzo);
void *visitatore(void *t);

void init(parco *p) {
	int i;
	p->posti_liberi=MaxP;
	p->bici_libere=MaxB;
	p->mono_liberi=MaxM;
	sem_init(&(p->S),0,0);        //Inizializza il semaforo con valore iniziale 0
	p->sospesi=0;
	pthread_mutex_init(&p->m, NULL);
}

void entra(parco *p, int mezzo) {
	pthread_mutex_lock(&p->m);     //Blocca il mutex per accedere ai dati condivisi
	if (mezzo==BICI) {
		while(p->posti_liberi==0 || p->bici_libere==0 ) {
			p->sospesi++;
			pthread_mutex_unlock(&p->m);
			sem_wait(&(p->S));                          //Aspetta sul semaforo
			pthread_mutex_lock(&p->m);
			p->sospesi--;
		}
		p->posti_liberi--;
		p->bici_libere--;
	}
	else {
		while(p->posti_liberi==0 || p->mono_liberi==0 ) {
			p->sospesi++;
			pthread_mutex_unlock(&p->m);
			sem_wait(&(p->S));                         //Aspetta sul semaforo
			pthread_mutex_lock(&p->m);
			p->sospesi--;
		}
		p->posti_liberi--;
		p->mono_liberi--;
	}
	pthread_mutex_unlock (&p->m);                             //Rilascia il mutex dopo aver aggiornato lo stato
}

void esci(parco *p, int mezzo)
{	int i,k;
	pthread_mutex_lock(&p->m);
	p->posti_liberi++;
	if (mezzo==BICI)
		p->bici_libere++;
	else
		p->mono_liberi++;
	for (i=0; i<p->sospesi; i++)               //Risveglia tutti i visitatori in attesa
		sem_post(&p->S); 
	pthread_mutex_unlock (&p->m);
}

//Funzione eseguita da ciascun thread (visitatore)
void *visitatore(void *t) { 
	int mezzo;
	int TH=(int)t;
	mezzo=rand()%2;
	entra(&P, mezzo);    //Scelta casuale del mezzo (0 = bici, 1 = monopattino)
	printf("entrato il visitatore n. %d con mezzo %d (0 bici, 1 monopattino)\n\n", TH, mezzo);
	sleep(rand()%3);
	esci(&P,mezzo);
	printf("uscito il visitatore n. %d con mezzo %d (0 bici, 1 monopattino)\n\n", TH, mezzo);
}


int main (int argc, char *argv[]) {
	pthread_t thread[NUM_THREADS];
	int rc, num;
	long t;
	float media, max;
	void *status;
	srand(time(NULL));
	init(&P);
	for(t=0; t<NUM_THREADS; t++) {
		rc = pthread_create(&thread[t], NULL, visitatore, (void *)t);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}
	}
	for(t=0; t<NUM_THREADS; t++) {
		rc = pthread_join(thread[t], &status);
		if (rc)
			printf("ERRORE join thread %ld codice %d\n", t, rc);
	}
	return 0;
}
