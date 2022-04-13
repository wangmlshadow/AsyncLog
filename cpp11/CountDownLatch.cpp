#include "CountDownLatch.h"

CountDownLatch::CountDownLatch(int count)
    : count_(count) {}

void CountDownLatch::wait() {
  //MutexLockGuard lock(mutex_);
  std::unique_lock<std::mutex> lock(mutex_);
  //while (count_ > 0) condition_.wait();
  condition_.wait(lock, [=] {return count_ <= 0;});
}

void CountDownLatch::countDown() {
  //MutexLockGuard lock(mutex_);
  std::unique_lock<std::mutex> lock(mutex_);
  --count_;
  if (count_ == 0) condition_.notify_all();
}