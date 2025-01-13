#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#define NUM_THREADS 3
#define K 2

// We need a data structure organized as follows: Each thread corresponds to a person.
// Each person gives K opinions on K movies, assigning a rating between 1 and 10.
// (The fact that the number of opinions matches the number of movies is purely coincidental.)

typedef struct {
    int ratings[K];
    char movies[K][40];
    int opinions_count;
    pthread_mutex_t mutex;
} survey;

survey Survey;

// Initialize the "database" of movies
void initialize(survey *s) {
    int i;
    s->opinions_count = 0;
    for (i = 0; i < K; i++) {
        printf("What is the name of movie number %d? ", i + 1);
        scanf("%s", s->movies[i]);
        s->ratings[i] = 0;
    }
    pthread_mutex_init(&s->mutex, NULL);
}

// Function executed by each "user" (thread)
void submit_opinions(survey *s, int thread_id) {         
    int i, rating;
    // Temporarily lock the database to avoid concurrent writes
    pthread_mutex_lock(&s->mutex);
    printf("\n\n COMPLETING QUESTIONNAIRE for Viewer %d:\n", thread_id); 
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
    // Unlock to prevent deadlocks
    pthread_mutex_unlock(&s->mutex);
}

// Viewer = thread
void *viewer(void *thread_data) {
    long thread_id, result = 0;
    thread_id = (int)thread_data;
    submit_opinions(&Survey, thread_id);
    printf("Viewer %ld completed the questionnaire...\n", thread_id);
    pthread_exit((void *)result);
}

// Create and synchronize all threads
int main(int argc, char *argv[]) {
    pthread_t threads[NUM_THREADS];
    int return_code, top_movie_index;
    long thread_id;
    float average, highest_average;
    void *status;

    initialize(&Survey);

    // After initializing the database, create all threads to proceed independently
    for (thread_id = 0; thread_id < NUM_THREADS; thread_id++) {
        return_code = pthread_create(&threads[thread_id], NULL, viewer, (void *)thread_id); 
        if (return_code) {
            printf("ERROR: %d\n", return_code);
            exit(-1);  
        }
    }

    // Wait for all threads to finish to ensure all questionnaires are completed
    for (thread_id = 0; thread_id < NUM_THREADS; thread_id++) {
        return_code = pthread_join(threads[thread_id], &status);
        if (return_code) {
            printf("ERROR joining thread %ld, code %d\n", thread_id, return_code);
        }
    }

    printf("\n\n--- RESULTS ---\n");
    top_movie_index = 0;
    highest_average = 0;
    for (thread_id = 0; thread_id < K; thread_id++) {
        average = (float)Survey.ratings[thread_id] / NUM_THREADS;
        printf("Average rating for movie %ld (%s): %f\n", thread_id + 1, Survey.movies[thread_id], average);
        if (average > highest_average) {
            highest_average = average;
            top_movie_index = thread_id;
        }
    }
    printf("\n\n THE WINNING MOVIE IS: %s, with a rating of %f!\n", Survey.movies[top_movie_index], highest_average);
    return 0;
}
