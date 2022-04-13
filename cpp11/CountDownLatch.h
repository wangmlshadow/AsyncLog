#pragma once
// #include "Condition.h"
// #include "MutexLock.h"
#include "noncopyable.h"
#include <mutex>
#include <condition_variable>

// CountDownLatch的主要作用是确保Thread中传进去的func真的启动了以后
// 外层的start才返回
class CountDownLatch : noncopyable {
 public:
  explicit CountDownLatch(int count);
  void wait();
  void countDown();

 private:
  //mutable MutexLock mutex_;
  std::mutex mutex_;
  //Condition condition_;
  std::condition_variable condition_;
  int count_;
};