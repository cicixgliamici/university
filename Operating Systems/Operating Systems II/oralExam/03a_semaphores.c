#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

/*
 * A counting semaphore implementation using POSIX threads.
 *
 * The semaphore structure maintains a counter that represents the number of available resources.
 * It uses a mutex to ensure mutual exclusion when accessing or modifying the counter and a condition
 * variable to allow threads to wait until a resource is available.
 *
 * The main functions provided are:
 *  - semaphore_init: Initialize the semaphore with a given initial count.
 *  - semaphore_wait: Decrement the semaphore (wait operation). If the count is zero, the thread blocks until it can decrement.
 *  - semaphore_signal: Increment the semaphore (signal operation) and wake up a waiting thread if any.
 *
 * Theoretical background:
 * -----------------------
 * Semaphores are a synchronization mechanism used to control access to a common resource.
 * A counting semaphore has a non-negative integer value. The two primary operations are:
 *
 *   wait (also known as P or down operation):
 *     - If the semaphore's count is greater than zero, decrement it and allow the thread to continue.
 *     - If the count is zero, block the thread until another thread increments the semaphore.
 *
 *   signal (also known as V or up operation):
 *     - Increment the semaphore's count.
 *     - If there are any threads waiting, wake one up so it can proceed.
 *
 * This mechanism ensures that only a limited number of threads can access the resource simultaneously,
 * which is crucial for preventing race conditions and ensuring data consistency in concurrent environments.
 */

typedef struct {
    int count;                  // The count representing available resources
    pthread_mutex_t mutex;      // Mutex to protect access to count
    pthread_cond_t cond;        // Condition variable for blocking and waking up threads
} semaphore;

/*
 * Initialize the semaphore.
 * 'initial' is the initial count of the semaphore (number of resources available).
 */
void semaphore_init(semaphore *sem, int initial) {
    sem->count = initial;
    pthread_mutex_init(&(sem->mutex), NULL);
    pthread_cond_init(&(sem->cond), NULL);
}

/*
 * Wait (or P) operation on the semaphore.
 * Decrements the semaphore's count. If the count is zero, the thread will block until it can decrement.
 */
void semaphore_wait(semaphore *sem) {
    // Lock the mutex to ensure exclusive access to the count.
    pthread_mutex_lock(&(sem->mutex));

    // Loop to handle spurious wake-ups:
    while (sem->count <= 0) {
        // Block the thread until a resource becomes available.
        pthread_cond_wait(&(sem->cond), &(sem->mutex));
    }
    // A resource is now available, so decrement the count.
    sem->count--;

    // Unlock the mutex to allow other threads to access the semaphore.
    pthread_mutex_unlock(&(sem->mutex));
}

/*
 * Signal (or V) operation on the semaphore.
 * Increments the semaphore's count, indicating that a resource has been released.
 * If there are any waiting threads, one is woken up.
 */
void semaphore_signal(semaphore *sem) {
    // Lock the mutex to ensure exclusive access.
    pthread_mutex_lock(&(sem->mutex));

    // Increment the count (release a resource).
    sem->count++;

    // Signal one waiting thread that a resource is now available.
    pthread_cond_signal(&(sem->cond));

    // Unlock the mutex.
    pthread_mutex_unlock(&(sem->mutex));
}

/* Example usage of the semaphore to protect a critical section */
#define NUM_THREADS 5

semaphore sem; // Global semaphore

// A shared resource (for example, a counter)
int shared_resource = 0;

void* thread_function(void* arg) {
    int thread_id = *((int*)arg);

    // Simulate waiting for resource access using semaphore_wait
    printf("Thread %d: Waiting to access the shared resource.\n", thread_id);
    semaphore_wait(&sem);

    // Begin critical section
    printf("Thread %d: Entered critical section.\n", thread_id);
    int local = shared_resource;
    local++;  // Modify shared resource
    sleep(1); // Simulate some work being done
    shared_resource = local;
    printf("Thread %d: Exiting critical section. Shared resource = %d\n", thread_id, shared_resource);
    // End critical section

    // Signal semaphore to release the resource
    semaphore_signal(&sem);

    pthread_exit(NULL);
}

int main() {
    pthread_t threads[NUM_THREADS];
    int thread_ids[NUM_THREADS];

    // Initialize semaphore with a count of 1 for mutual exclusion (binary semaphore)
    semaphore_init(&sem, 1);

    // Create multiple threads that use the shared resource.
    for (int i = 0; i < NUM_THREADS; i++) {
        thread_ids[i] = i + 1;
        if (pthread_create(&threads[i], NULL, thread_function, (void*)&thread_ids[i]) != 0) {
            perror("Failed to create thread");
            exit(EXIT_FAILURE);
        }
    }

    // Wait for all threads to complete.
    for (int i = 0; i < NUM_THREADS; i++) {
        pthread_join(threads[i], NULL);
    }

    // Clean up (destroy mutex and condition variable)
    pthread_mutex_destroy(&sem.mutex);
    pthread_cond_destroy(&sem.cond);

    printf("Final value of shared resource: %d\n", shared_resource);
    return 0;
}

/*
 * Explanation:
 *
 * This C program implements a counting semaphore to manage access to a shared resource among multiple threads.
 * 
 * Key Components:
 * 1. Semaphore Structure:
 *    - count: Represents the number of available resources.
 *    - mutex: Used to guarantee mutually exclusive access to the count variable.
 *    - cond: A condition variable that allows threads to wait until a resource becomes available.
 *
 * 2. Main Functions:
 *    - semaphore_init: Initializes the semaphore with a specified initial count and sets up the mutex and condition variable.
 *    - semaphore_wait (P operation): Decrements the semaphore's count. If the count is zero, the thread will block until a resource is released.
 *    - semaphore_signal (V operation): Increments the semaphore's count, indicating that a resource has been released, and signals a waiting thread if there is one.
 *
 * 3. Usage:
 *    - In the main function, a binary semaphore (count = 1) is initialized to ensure exclusive access to the critical section.
 *    - Multiple threads are created, each waiting to access a shared resource.
 *    - When a thread enters the critical section, it increments the shared_resource variable, simulating resource usage.
 *    - The semaphore ensures that only one thread at a time can modify the shared resource, preventing race conditions.
 *
 * The theoretical foundation of semaphores is critical in concurrent programming, ensuring that shared resources are safely managed
 * and that simultaneous access does not corrupt data.
 */
