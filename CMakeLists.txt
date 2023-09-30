cmake_minimum_required(VERSION 3.25)
project(strategy_simulation)

set(CMAKE_CXX_STANDARD 20)

find_package(Eigen3 CONFIG REQUIRED)
target_link_libraries(your_target_name PRIVATE Eigen3::Eigen)

add_executable(strategy_simulation src/strategysim.cpp src/lib/matplotlibcpp.h)
