#include "Logging.h"
//#include "CurrentThread.h"
//#include "Thread.h"
#include "AsyncLogging.h"
#include <assert.h>
#include <iostream>
#include <mutex>
#include <time.h>  
#include <sys/time.h>



//static pthread_once_t once_control_ = PTHREAD_ONCE_INIT;
// static SingletonAsyncLogging *SingletonAsyncLogging_;
// AsyncLogging* SingletonAsyncLogging::AsyncLogger_ = nullptr;
// std::mutex SingletonAsyncLogging::mutex_;

static AsyncLogging* AsyncLogger_;
//static std::mutex mutex_;
static std::once_flag flag;
std::string Logger::logFileName_ = "./WebServer.log";

void once_init()
{
    AsyncLogger_ = new AsyncLogging(Logger::getLogFileName());
    AsyncLogger_->start(); 
}

// void once_init() {
//     if (AsyncLogger_ == nullptr) {
//         std::unique_lock<std::mutex> lock(mutex_);
//         if (AsyncLogger_ == nullptr) {
//             AsyncLogger_ = new AsyncLogging(Logger::getLogFileName());
//         }
//     }
// }

void output(const char* msg, int len)
{
    // pthread_once(&once_control_, once_init);
    // SingletonAsyncLogging_ = new SingletonAsyncLogging();
    // SingletonAsyncLogging_->getInstance(Logger::getLogFileName())->start();
    // SingletonAsyncLogging_->getInstance(Logger::getLogFileName())->append(msg, len);
    //once_init();
    std::call_once(flag, once_init);
    //AsyncLogger_->start();
    AsyncLogger_->append(msg, len);
}

Logger::Impl::Impl(const char *fileName, int line)
  : stream_(),
    line_(line),
    basename_(fileName)
{
    formatTime();
}

void Logger::Impl::formatTime()
{
    struct timeval tv;
    time_t time;
    char str_t[26] = {0};
    gettimeofday (&tv, NULL);
    time = tv.tv_sec;
    struct tm* p_time = localtime(&time);   
    strftime(str_t, 26, "%Y-%m-%d %H:%M:%S\n", p_time);
    stream_ << str_t;
}

Logger::Logger(const char *fileName, int line)
  : impl_(fileName, line)
{ }

Logger::~Logger()
{
    impl_.stream_ << " -- " << impl_.basename_ << ':' << impl_.line_ << '\n';
    const LogStream::Buffer& buf(stream().buffer());
    output(buf.data(), buf.length());
}