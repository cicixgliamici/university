#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define NUM_THREADS 9
#define K 2

// Data structure for storing survey information
typedef struct {
    int ratings[K];               
    char movies[K][40];           
    int opinions_count;           
    int winning_index;           
    pthread_mutex_t mutex;        
} survey;

// Data structure for implementing a barrier synchronization mechanism
typedef struct {
    int completed;                
    sem_t barrier_mutex;          
    sem_t barrier_gate;           
} barrier_sync;

survey Survey;
barrier_sync Barrier;

// Initialize the survey database
void initialize_survey(survey *s) {
    int i;
    s->opinions_count = 0;
    for (i = 0; i < K; i++) {
        printf("What is the name of movie number %d? ", i + 1);
        scanf("%s", s->movies[i]);
        s->ratings[i] = 0;
    }
    pthread_mutex_init(&s->mutex, NULL);
}

// Initialize the barrier as empty
void initialize_barrier(barrier_sync *b) {
    sem_init(&b->barrier_mutex, 0, 1);
    sem_init(&b->barrier_gate, 0, 0);
    b->completed = 0;
}

// Function for a viewer (thread) to submit their ratings
void submit_opinions(survey *s, int thread_id) {
    int i, rating;
    pthread_mutex_lock(&s->mutex);
    printf("\n\n FILLING OUT QUESTIONNAIRE for Viewer %d:\n", thread_id);
    for (i = 0; i < K; i++) {
        printf("Rating for movie %s [0 to 10]? ", s->movies[i]);
        scanf("%d", &rating);
        s->ratings[i] += rating;
    }
    s->opinions_count++;
    printf("END OF QUESTIONNAIRE for Viewer %d\n PARTIAL SURVEY RESULTS:\n", thread_id);
    for (i = 0; i < K; i++) {
        printf("Average rating for movie %s: %f\n", s->movies[i], (float)(s->ratings[i]) / s->opinions_count);
    }
    pthread_mutex_unlock(&s->mutex);
}

// Barrier implementation
void pass_barrier(barrier_sync *b, survey *s) {
    sem_wait(&b->barrier_mutex);
    b->completed++;

    // If all threads have reached the barrier, determine the winning movie
    if (b->completed == NUM_THREADS) {
        int t, max_index = 0;
        float max_rating = 0, avg_rating;

        for (t = 0; t < K; t++) {
            avg_rating = (float) s->ratings[t] / NUM_THREADS;
            printf("Average rating for movie %d (%s): %f\n", t + 1, s->movies[t], avg_rating);
            if (avg_rating > max_rating) {
                max_rating = avg_rating;
                max_index = t;
            }
        }
        printf("\n\nTHE WINNING MOVIE IS: %s, with a rating of %f!\n", s->movies[max_index], max_rating);
        s->winning_index = max_index;
        sem_post(&b->barrier_gate);
    }
    sem_post(&b->barrier_mutex);

    // Ensure all threads wait until the barrier is released
    sem_wait(&b->barrier_gate);
    sem_post(&b->barrier_gate);
}

// Simulate viewing the winning movie
void watch_movie(int thread_id, survey *s) {
    pthread_mutex_lock(&s->mutex);
    printf("Thread %d is downloading the movie %s...\n", thread_id, s->movies[s->winning_index]);
    pthread_mutex_unlock(&s->mutex);
}

// Function executed by each thread
void *viewer(void *thread_data) {
    long thread_id = (long)thread_data;
    submit_opinions(&Survey, thread_id);
    pass_barrier(&Barrier, &Survey);
    watch_movie(thread_id, &Survey);
    pthread_exit(NULL);
}

int main(int argc, char *argv[]) {
    pthread_t threads[NUM_THREADS];
    int rc;
    long t;

    // Initialize the survey and barrier
    initialize_survey(&Survey);
    initialize_barrier(&Barrier);

    // Create threads
    for (t = 0; t < NUM_THREADS; t++) {
        rc = pthread_create(&threads[t], NULL, viewer, (void *)t);
        if (rc) {
            printf("ERROR: Unable to create thread %ld, error code %d\n", t, rc);
            exit(-1);
        }
    }

    // Wait for all threads to complete
    for (t = 0; t < NUM_THREADS; t++) {
        rc = pthread_join(threads[t], NULL);
        if (rc) {
            printf("ERROR joining thread %ld, error code %d\n", t, rc);
        }
    }

    return 0;
}
