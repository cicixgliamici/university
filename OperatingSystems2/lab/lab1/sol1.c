//  !! GIVEN EXAMPLE FROM UNIVERSITY - ONLY SLIGHT CHANGES !!

#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define NUM_THREADS 3
#define K 2

//Ho bisogno di una SD organizzata così: Ogni thread corrisponde a una persona
//                                       ogni persona esprime K pareri su K film assegnando un voto 1<=x<=10
//                                       (il fatto che il nuemro di pareri coincida con il numero di film è puramente incidimentale)

typedef struct{
int voti[K];
char film[K][40];
int pareri;
pthread_mutex_t m;
} sondaggio;

sondaggio S;

//Inizializzazione della "base di dati"=db dei film
void init(sondaggio *s) {
	int i;
	s->pareri=0;
	for(i=0;i<K; i++)
	{	printf("Qual è il nome del film numero %d ? ", i+1);
		scanf("%s",s->film[i]);
		s->voti[i]=0;
	}
	pthread_mutex_init(&s->m, NULL); 
}

//Funzione sottoposta ad ogni "utente"(thread)
void esprimi_pareri(sondaggio *s, int th) {         
  int i, voto;
  //Chiudi momentaneamente il db per evitare che altri scrivano nel mentre
  pthread_mutex_lock(&s->m);
  printf("\n\n COMPILAZIONE QUESTIONARIO per lo Spettatore %d:\n", th); 
  for(i=0;i<K; i++) {
		printf("voto del film  %s [0,.. 10]? ", s->film[i]);
		scanf("%d", &voto);
		s->voti[i]+=voto;
	}
   s->pareri++;
   printf("FINE QUESTIONARIO per lo spettatore %d\n RISULTATI PARZIALI SONDAGGIO:\n", th);
   for(i=0;i<K;i++)
	printf("Valutazione media film %s: %f\n", s->film[i], (float)(s->voti[i])/s->pareri);
   //Riapri altrimenti facciamo deadlock
   pthread_mutex_unlock (&s->m);
 }

//Spettatore = thread
void *spettatore(void *t) {
	long tid, result=0;
	tid = (int)t;
	esprimi_pareri(&S, tid);
        printf("Spettatore %ld ha compilato i questionari...\n",tid);
	pthread_exit((void*) result);
}

//Crea e sincronizza tutti i thread
int main (int argc, char *argv[]) {
   pthread_t thread[NUM_THREADS];
   int rc, i_max;
   long t;
   float media, max;
   void *status;
   init(&S);
   //Dopo aver inizializzato il db crea tutti i thread e loro proseguono autonomamente,	
   for(t=0; t<NUM_THREADS; t++) {
      rc = pthread_create(&thread[t], NULL, spettatore, (void *)t); 
      if (rc) {
         printf("ERRORE: %d\n", rc);
         exit(-1);  
      }
   }
   //Pertanto va aspettato che tutti terminino e si sincronizzino, altimenti si prosegue senza avere tutti i sondaggi	
   for(t=0; t<NUM_THREADS; t++) {
      rc = pthread_join(thread[t], &status);
      if (rc) 
          printf("ERRORE join thread %ld codice %d\n", t, rc);
   }
   printf("\n\n--- RISULTATI ---\n");
   i_max=0; max=0;
   for(t=0; t<K;t++) {
  		media=(float) S.voti[t]/NUM_THREADS;
		printf("Valutazione media del film n.%ld (%s): %f\n", t+1, S.film[t], media);
		if (media>max) {
			max=media;
			i_max=t;
		}
  }
  printf("\n\n IL FILM VINCITORE E': %s, con voto %f !\n",  S.film[i_max], max);
  return 0;
}
