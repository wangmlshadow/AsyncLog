#pragma once
#include <functional>
#include <string>
#include <vector>
#include <mutex>
#include <thread>
#include <condition_variable>
#include "CountDownLatch.h"
#include "LogStream.h"
//#include "MutexLock.h"
//#include "Thread.h"
#include "noncopyable.h"


class AsyncLogging : noncopyable {
 public:
  AsyncLogging(const std::string basename, int flushInterval = 2);
  ~AsyncLogging() {
    if (running_) stop();
  }
  void append(const char* logline, int len);

  void start() {
    running_ = true;
    // thread_.start();
    // latch_.wait();
    thread_ = std::unique_ptr<std::thread>(new std::thread(&AsyncLogging::threadFunc, this));
    latch_.wait();
  }

  void stop() {
    running_ = false;
    //cond_.notify();
    //thread_.join();
    cond_.notify_all();
    thread_->join();
  }

 private:
  void threadFunc();
  typedef FixedBuffer<kLargeBuffer> Buffer;
  typedef std::vector<std::shared_ptr<Buffer>> BufferVector;
  typedef std::shared_ptr<Buffer> BufferPtr;
  const int flushInterval_;
  bool running_;
  std::string basename_;
  //Thread thread_;
  std::unique_ptr<std::thread> thread_;
  //MutexLock mutex_;
  std::mutex mutex_;
  //Condition cond_;
  std::condition_variable cond_;
  BufferPtr currentBuffer_;
  BufferPtr nextBuffer_;
  BufferVector buffers_;
  CountDownLatch latch_;
};

// class SingletonAsyncLogging : noncopyable {
// public:
//     SingletonAsyncLogging() {}
//     SingletonAsyncLogging(const SingletonAsyncLogging&)=delete;
//     SingletonAsyncLogging& operator=(const SingletonAsyncLogging&)=delete;
// public:
//     static AsyncLogging* getInstance(const std::string basename) {
//         if (AsyncLogger_ == nullptr) {
//             std::unique_lock<std::mutex> lock(mutex_);
//             if (AsyncLogger_ == nullptr) {
//                 AsyncLogger_ = new AsyncLogging(basename);
//             }
//         }
//         return AsyncLogger_;
//     }
// private:
//     static AsyncLogging* AsyncLogger_;
//     static std::mutex mutex_;
// };

// AsyncLogging* SingletonAsyncLogging::AsyncLogger_ = nullptr;
// std::mutex SingletonAsyncLogging::mutex_;