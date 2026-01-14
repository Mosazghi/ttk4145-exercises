// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you
// use gcc directly

#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>

int i = 0;
sem_t sem;

// Note the return type: void*
void *incrementingThreadFunction() {
  for (int j = 0; j < 1000000; ++j) {
    sem_wait(&sem);
    printf("INCREMENTING\n");
    ++i;
    sem_post(&sem);
  }

  return NULL;
}

void *decrementingThreadFunction() {
  for (int j = 0; j < 1000001; ++j) {
    sem_wait(&sem);
    printf("DECREMENTING\n");
    --i;
    sem_post(&sem);
  }
  return NULL;
}

int main() {
  pthread_t thread_1;
  pthread_t thread_2;
  sem_init(&sem, 0, 1);

  if (pthread_create(&thread_1, NULL, incrementingThreadFunction, NULL) != 0) {
    fprintf(stderr, "Thread 1 creation failed");
    return 1;
  }

  if (pthread_create(&thread_2, NULL, decrementingThreadFunction, NULL) != 0) {
    fprintf(stderr, "Thread 2 creation failed");
    return 1;
  }

  pthread_join(thread_1, NULL);
  pthread_join(thread_2, NULL);

  printf("The magic number is: %d\n", i);
  return 0;
}
