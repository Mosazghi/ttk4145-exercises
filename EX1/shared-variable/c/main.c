// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you
// use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t mtx;

// Note the return type: void*
void *incrementingThreadFunction() {
  for (int j = 0; j < 1000000; ++j) {
    pthread_mutex_lock(&mtx);
    printf("INC\n");
    ++i;
    pthread_mutex_unlock(&mtx);
  }

  return NULL;
}

void *decrementingThreadFunction() {
  for (int j = 0; j < 1000001; ++j) {
    pthread_mutex_lock(&mtx);
    printf("DEC\n");
    --i;
    pthread_mutex_unlock(&mtx);
  }
  return NULL;
}

int main() {
  pthread_t thread_1;
  pthread_t thread_2;
  pthread_mutex_init(&mtx, NULL);
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
