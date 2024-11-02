//  !! GIVEN EXAMPLE FROM UNIVERSITY - ONLY SLIGHT CHANGES !!

#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define N 20
#define K 4

//Given an empty array of N=20 elements, we will divide it in K=4 portions
//Then we will randomize it and we will create M=N/K=5 threads
//We will search each maximum in the M set and then find the maximum above all sets
int V[N];


void *Calcolo(void *t) {                               //"Worker", the segment where we search the maximum 
    int first, result=0;
    first = (int)t;
    for (int i=first; i<first+K; i++)
	if (V[i]>result)
		result=V[i];
    printf("Local Maximum: %d\n", result); 
    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t thread[N/K];
    int rc;
    int M, t, first, status, max=0;
    M=N/K;	                                      //Number of threads/worker
    srand(time(0));                                   //Randomize using the seed "time"
    printf("Vector initialization V:\n");
    for(int i = 0; i < N; i++){
	   V[i]=1+rand()%200;                         //Random numbers between 1-200
	   printf("%d\t", V[i]);
   	}	
    printf("\n");
    for(t=0; t<M; t++) {
        printf("Main: Thread creation n.%d\n", t);
	first=t*K;                                                       //Giving each thread the index of the first element of the sets they will use
        rc = pthread_create(&thread[t], NULL, Calcolo, (void *)first);
        if (rc) {
            printf("Error: %d\n", rc);
            exit(-1); 
	}
    }
    for(t=0; t<M; t++) {
        rc = pthread_join(thread[t], (void *)&status);
	if (rc)
		printf("ERROR join thread %ld code %d\n", t, rc);
	else {
		printf("Finished thread %ld with answer %d\n",t,status);
		if (status>max)
			max=status;
	}
   }
   printf("True answer: %d\n", max);
}
