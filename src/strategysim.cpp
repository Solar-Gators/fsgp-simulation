#include <iostream>

int main(int argc, char *argv[]) {
    if(argc > 1) {
        std::cout << "First argument: " << argv[1] << std::endl;
    }
    if(argc > 2) {
        std::cout << "Second argument: " << argv[2] << std::endl;
    }
    return 0;
}
