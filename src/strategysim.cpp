#include <iostream>
#include <vector>
#include <Eigen/Dense>
#include <matplotlibcpp.h>

constexpr int segmentCount = 40;
constexpr double segmentLengths[segmentCount] = { /*your lengths here*/ };

std::vector<double> create_piecewise(std::vector<double>& args) {
    std::vector<double> result;
    double c = args[0];
    for (int i = 1; i < args.size(); i += 3) {
        Eigen::Vector3d parabola(args[i], args[i+1], c);
        for (int j = 0; j < segmentLengths[i / 3]; ++j) {
            result.push_back(parabola[0] * j * j + parabola[1] * j + parabola[2]);
        }
        c = result.back();
        if (i + 2 < args.size()) {
            Eigen::Vector2d line(args[i+2], c);
            for (int j = 0; j < segmentLengths[i / 3 + 1]; ++j) {
                result.push_back(line[0] * j + line[1]);
            }
            c = result.back();
        }
    }
    return result;
}

int main(int argc, char* argv[]) {
    std::vector<double> args(argc - 1);
    for (int i = 1; i < argc; ++i) {
        args[i - 1] = std::stod(argv[i]);
    }
    std::vector<double> piecewise = create_piecewise(args);
    matplotlibcpp::plot(piecewise);
    matplotlibcpp::show();
    return 0;
}