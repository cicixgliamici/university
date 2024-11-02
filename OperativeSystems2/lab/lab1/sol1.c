#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define NUM_THREADS 3
#define K 2

typedef struct{
int voti[K];
char film[K][40];
int pareri;
pthread_mutex_t m;
} sondaggio;

sondaggio S;

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

void esprimi_pareri(sondaggio *s, int th) {
  int i, voto;
  pthread_mutex_lock(&s->m);
  printf("\n\n COMPILAZIONE QUESTIONARIO per lo Spettatore %d:\n", th); 
  for(i=0;i<K; i++)
	{	printf("voto del film  %s [0,.. 10]? ", s->film[i]);
		scanf("%d", &voto);
		s->voti[i]+=voto; //accumulo voti 
	}
   s->pareri++;
   printf("FINE QUESTIONARIO per lo spettatore %d\n RISULTATI PARZIALI SONDAGGIO:\n", th);
   for(i=0;i<K;i++)
	printf("Valutazione media film %s: %f\n", s->film[i], (float)(s->voti[i])/s->pareri);
   pthread_mutex_unlock (&s->m);
 }


void *spettatore(void *t) // codice spettatore
{	long tid, result=0;
	tid = (int)t;
	esprimi_pareri(&S, tid);
    //printf("Spettatore %ld ha compilato i questionari...\n",tid);
	pthread_exit((void*) result);
}


int main (int argc, char *argv[])
{  pthread_t thread[NUM_THREADS];
   int rc, i_max;
   long t;
   float media, max;
   void *status;
   
   init(&S);
  
   for(t=0; t<NUM_THREADS; t++) {
      rc = pthread_create(&thread[t], NULL, spettatore, (void *)t); 
      if (rc) {
         printf("ERRORE: %d\n", rc);
         exit(-1);   }
  }
	for(t=0; t<NUM_THREADS; t++) {
      rc = pthread_join(thread[t], &status);
      if (rc) 
		   printf("ERRORE join thread %ld codice %d\n", t, rc);
   }
  printf("\n\n--- R I S U L T A T I ---\n");
  i_max=0; max=0;
  for(t=0; t<K;t++)
  {		media=(float) S.voti[t]/NUM_THREADS;
		printf("Valutazione media del film n.%ld (%s): %f\n", t+1, S.film[t], media);
		if (media>max)
		{	max=media;
			i_max=t;
		}
  }
  printf("\n\n IL FILM VINCITORE E': %s, con voto %f !\n",  S.film[i_max], max);
  
  return 0;
}
