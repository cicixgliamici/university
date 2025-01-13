# Understanding Pthreads and Synchronization in C

## Introduction to Pthreads

Pthreads (POSIX threads) is a standard API for creating and managing threads in C. Threads enable parallel execution within the same process, sharing memory space and resources efficiently. This capability is especially useful for multi-core processors, where parallelism can significantly enhance performance.

### Key Components of Pthreads

1. **Thread Creation**: `pthread_create` is used to spawn a new thread.
2. **Thread Joining**: `pthread_join` waits for a thread to complete its execution.
3. **Thread Exit**: `pthread_exit` allows a thread to terminate.
4. **Mutexes**: `pthread_mutex_t` provides mutual exclusion to synchronize access to shared resources.

---

## Synchronization Techniques

### 1. **Mutexes**
Mutexes ensure mutual exclusion, meaning only one thread can access a critical section at a time. This is achieved using the following:

- `pthread_mutex_init`: Initializes a mutex.
- `pthread_mutex_lock`: Locks the mutex, preventing other threads from entering the critical section.
- `pthread_mutex_unlock`: Unlocks the mutex, allowing other threads to proceed.
- `pthread_mutex_destroy`: Destroys a mutex when it's no longer needed.

#### Example
```c
pthread_mutex_t mutex;
pthread_mutex_init(&mutex, NULL);

pthread_mutex_lock(&mutex);
// Critical section
pthread_mutex_unlock(&mutex);

pthread_mutex_destroy(&mutex);
```

### 2. **Semaphores**
Semaphores are used for signaling and are useful in scenarios where threads need to coordinate.

- `sem_init`: Initializes a semaphore.
- `sem_wait`: Decrements the semaphore value (blocks if the value is 0).
- `sem_post`: Increments the semaphore value, potentially unblocking a waiting thread.
- `sem_destroy`: Destroys the semaphore.

#### Example
```c
sem_t semaphore;
sem_init(&semaphore, 0, 1);

sem_wait(&semaphore);
// Critical section
sem_post(&semaphore);

sem_destroy(&semaphore);
```

### 3. **Barriers**
Barriers ensure that all threads reach a certain point before any of them can proceed.

In the provided code, barriers are implemented using semaphores to hold threads until a condition is met (e.g., all threads have completed a specific task).

---

## Example Code Analysis

### Example 1: Finding Maximum Value in an Array
This code divides an array into segments, assigns each segment to a thread, and computes the local maximum in parallel. Finally, the global maximum is determined by comparing the local maxima.

Key Features:
- Threads independently compute local maxima using `pthread_create`.
- The main thread waits for all threads to complete using `pthread_join`.
- Mutexes are not needed here since each thread operates on independent segments.

#### Code Highlights
```c
pthread_t thread[N / K];
for (t = 0; t < M; t++) {
    pthread_create(&thread[t], NULL, Calcolo, (void *)first);
}

for (t = 0; t < M; t++) {
    pthread_join(thread[t], (void *)&status);
}
```

### Example 2: Survey with Mutexes
This example uses threads to simulate a survey where each thread corresponds to a user submitting ratings for movies. A mutex ensures synchronized access to the shared data structure.

Key Features:
- Mutex locks protect the critical section when updating shared data.
- Each thread submits ratings and computes partial results independently.

#### Code Highlights
```c
pthread_mutex_lock(&s->mutex);
// Update shared ratings array
pthread_mutex_unlock(&s->mutex);
```

### Example 3: Survey with Barrier Synchronization
This example extends the previous survey with a barrier mechanism to ensure all threads complete their work before the winning movie is determined.

Key Features:
- A semaphore-based barrier synchronizes threads, ensuring that all reach a common point.
- The barrier unlocks when all threads have completed their tasks.

#### Code Highlights
```c
sem_wait(&b->barrier_mutex);
if (b->completed == NUM_THREADS) {
    sem_post(&b->barrier_gate);
}
sem_post(&b->barrier_mutex);
sem_wait(&b->barrier_gate);
sem_post(&b->barrier_gate);
```

---

## Summary

Pthreads and synchronization techniques such as mutexes, semaphores, and barriers are essential tools for parallel programming in C. They enable safe and efficient execution of concurrent tasks. By understanding and applying these concepts, developers can maximize the performance and reliability of their multithreaded applications.
