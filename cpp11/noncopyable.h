#pragma once

class noncopyable {
 protected:
  noncopyable() {}
  ~noncopyable() {}

 private:
  noncopyable(const noncopyable&)=delete;
  const noncopyable& operator=(const noncopyable&)=delete;
};