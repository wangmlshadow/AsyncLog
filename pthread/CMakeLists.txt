cmake_minimum_required(VERSION 2.8)

SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++0x")
#SET(CMAKE_CXX_STANDARD 11)
#SET(CMAKE_C_COMPILER /usr/bin/gcc)
#SET(CMAKE_CXX_COMPILER /usr/bin/gcc)

SET(CMAKE_BUILD_TYPE "Debug")
SET(CMAKE_CXX_FLAGS_DEBUG "$ENV{CXXFLAGS} -O0 -Wall -g -ggdb")
SET(CMAKE_CXX_FLAGS_RELEASE "$ENV{CXXFLAGS} -O3 -Wall")
set(THREADS_PREFER_PTHREAD_FLAG ON)
project(Demo)


find_package(Threads)

aux_source_directory(. DIR_SRCS)

add_executable(Demo ${DIR_SRCS})

target_link_libraries(Demo ${CMAKE_THREAD_LIBS_INIT})