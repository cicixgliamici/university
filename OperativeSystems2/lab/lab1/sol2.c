#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define NUM_THREADS 9
#define K 2

typedef struct{
int voti[K];
char film[K][40];
int pareri;
int indiceVincitore;
pthread_mutex_t m;
} sondaggio;

typedef struct {
int completati;
sem_t mb;  // valore iniziale 1
sem_t barriera; // valore iniziale 0
}barriera_sincro;

sondaggio S;
barriera_sincro B;

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

void initBarriera(barriera_sincro *b)
{	sem_init(&b->mb,0, 1);
	sem_init(&b->barriera,0, 0);
	b->completati=0;
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

void passBarriera(barriera_sincro *b, sondaggio *s)
{   sem_wait(&b->mb);
    b->completati++;
    if (b->completati==NUM_THREADS)
    {  //RICERCA VINCITORE:
        int t, i_max=0; float max=0;
        float media;
        for(t=0; t<K;t++)
        {		media=(float) s->voti[t]/NUM_THREADS;
                printf("Valutazione media del film n.%d (%s): %f\n", t+1, s->film[t], media);
                if (media>max)
                {	max=media;
                    i_max=t;
                }
        }
        printf("\n\nIL FILM VINCITORE E': %s, con voto %f !\n",  s->film[i_max], max);
        s->indiceVincitore=i_max; // salvo nella struttura il vincitore del sondaggio
        sem_post(&b->barriera);
    }
    sem_post(&b->mb);
    sem_wait(&b->barriera);
    sem_post(&b->barriera);
    return;
}
        
void visione(int th, sondaggio *s)
{   
    pthread_mutex_lock(&s->m);
    printf("thread %d sta eseguendo download del film %s .. \n", th, s->film[s->indiceVincitore]);
    pthread_mutex_unlock(&s->m);
}
        
        
        

 
void *spettatore(void *t) // codice spettatore
{   long tid, result=0;
    tid = (long)t;
    esprimi_pareri(&S, tid);
    passBarriera(&B, &S);
    visione(tid, &S);
    pthread_exit((void*) result);
}


int main (int argc, char *argv[])
{  pthread_t thread[NUM_THREADS];
   int rc, i_max;
   long t;
   float media, max;
   void *status;
   
   init(&S);
   initBarriera(&B);
  
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
  return 0;
}
