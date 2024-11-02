//  esercitazione som 1.1

#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define N 20
#define K 4

int V[N];


void *Calcolo(void *t) // codice worker
{   int first, result=0;
    first = (int)t;
    for (int i=first; i<first+K; i++)
	if (V[i]>result)
		result=V[i];
    //printf("Worker ha calcolato il massimo locale: %d\n", result); 
    pthread_exit((void*) result);
}


int main (int argc, char *argv[])
{   pthread_t thread[N/K];
    int rc;
    int M, t, first, status, max=0;
	
    M=N/K;	
    srand(time(0)); 
    printf("inizializzazione vettore V:\n");
	for(int i = 0; i < N; i++){
	   V[i]=1+rand()%200; //numeri casuali tra 1 e 200
	   printf("%d\t", V[i]);
   	}	
    printf("\n");
    
    for(t=0; t<M; t++) {
        printf("Main: creazione thread n.%d\n", t);
	first=t*K; // passo ad ogni thread l'indice del primo elemento da elaborare
        rc = pthread_create(&thread[t], NULL, Calcolo, (void *)first);
        if (rc) {
            printf("ERRORE: %d\n", rc);
            exit(-1);   }
    }
	

    for(t=0; t<M; t++) {
        rc = pthread_join(thread[t], (void *)&status);
	if (rc)
		printf("ERRORE join thread %ld codice %d\n", t, rc);
	else {
		printf("Finito thread %ld con ris. %d\n",t,status);
		if (status>max)
			max=status;
	}
   }
   printf("main-risultato finale: %d\n", max);


    }
