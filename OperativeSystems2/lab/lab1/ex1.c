#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define N 20
#define K 4

/*
 * This program demonstrates the use of pthreads in C for parallel computation.
 * Pthreads (POSIX threads) are a threading API used to create and manage threads
 * in C programs. Threads allow tasks to execute concurrently within the same process,
 * sharing memory space. Here, we divide an array into segments, calculate the local 
 * maximum of each segment in parallel using threads, and then determine the overall maximum.
 */

// Array of N=20 elements
int V[N];

// Thread worker function: finds the maximum element in the segment assigned to the thread
void *Calcolo(void *t) {                               
    int first, result = 0; 			//Initialize the result to 0
    first = (int)t; 				//Starting index of the segment assigned to this thread
    for (int i = first; i < first + K; i++)     //Iterate over the assigned segment
        if (V[i] > result)
            result = V[i]; 			//Update the local maximum
    printf("Local Maximum: %d\n", result); 
    pthread_exit((void*) result); 		//Return the local maximum
}	

int main (int argc, char *argv[]) {
    pthread_t thread[N / K]; 			//Array to store thread identifiers
    int rc; 					//Return code for thread functions
    int M, t, first, status, max = 0; 		//Variables for number of threads, thread status, and global maximum
    M = N / K; 					//Number of threads/workers
    srand(time(0)); 				//Initialize the random seed with the current time
    
    // Initialize the array with random numbers between 1 and 200
    printf("Vector initialization V:\n");
    for (int i = 0; i < N; i++) {
        V[i] = 1 + rand() % 200; 		//Random numbers between 1-200
        printf("%d\t", V[i]);
    }
    printf("\n");
    
    // Create threads to compute local maxima
    for (t = 0; t < M; t++) {
        printf("Main: Thread creation n.%d\n", t);
        first = t * K; 				//Index of the first element in the segment assigned to the thread
        rc = pthread_create(&thread[t], NULL, Calcolo, (void *)first);
        if (rc) { 				//Check for errors in thread creation
            printf("Error: %d\n", rc);
            exit(-1); 
        }
    }
    
    // Wait for all threads to complete and find the global maximum
    for (t = 0; t < M; t++) {
        rc = pthread_join(thread[t], (void *)&status);
        if (rc) 				//Check for errors in thread joining
            printf("ERROR joining thread %ld code %d\n", t, rc);
        else {
            printf("Finished thread %ld with answer %d\n", t, status);
            if (status > max) 			//Update the global maximum if necessary
                max = status;
        }
    }
    
    // Print the final global maximum
    printf("True answer: %d\n", max);
}
