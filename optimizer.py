import os
import platform
import subprocess
import sys

from mystic.monitors import VerboseMonitor
from mystic.solvers import *
from mystic.strategy import Best1Bin
from mystic.termination import VTR

CALLS_BETWEEN_IMAGE = 20
MAX_VELOCITY = 40.0
MAX_ENERGY_CONS = 1300
MAX_ACCELERATION = 3.0  # m/s^2

try:
    subprocess.run(["go", "build", "."])
except:
    print("Ensure Go is installed! Using binaries...\n")

if platform.system() == "Windows":
    cli_program = "./strategy-simulation.exe"
else:
    cli_program = "./strategy-simulation"


def call_cli_program(x, endArg):
    return subprocess.run(
        [cli_program] + list(map(str, x)) + [str(endArg)],
        capture_output=True,
        text=True,
    ).stdout


def get_expected_argument_count():
    output = subprocess.run([cli_program], capture_output=True, text=True).stdout
    try:
        return int(output.split("Expected argument count:")[1].split("\n")[0])
    except (IndexError, ValueError):
        print("Could not determine the expected argument count from the CLI program.")
        sys.exit(1)


output_cache = {}
i = 0


def get_output(x):
    global i
    autoEndArg = "" if i % CALLS_BETWEEN_IMAGE == 0 else "none"
    i += 1
    x_tuple = tuple(x)
    if x_tuple not in output_cache:
        output_cache[x_tuple] = call_cli_program(x, autoEndArg)
    return output_cache[x_tuple]


def objective(x):
    output = get_output(x)
    try:
        # Parse the output for the required values
        time_elapsed = float(output.split("Time Elapsed (s):")[1].split("\n")[0])
        energy_consumption = float(
            output.split("Energy Consumption (W):")[1].split("\n")[0]
        )
        initial_velocity = abs(
            float(output.split("Initial Velocity (m/s):")[1].split("\n")[0])
        )
        final_velocity = float(output.split("Final Velocity (m/s):")[1].split("\n")[0])
        max_velocity = float(output.split("Max Velocity (m/s):")[1].split("\n")[0])
        min_velocity = float(output.split("Min Velocity (m/s):")[1].split("\n")[0])
        max_acceleration = float(
            output.split("Max Acceleration (m/s^2):")[1].split("\n")[0]
        )
        min_acceleration = float(
            output.split("Min Acceleration (m/s^2):")[1].split("\n")[0]
        )

        # Check energy consumption constraint
        if energy_consumption > MAX_ENERGY_CONS or energy_consumption < 0:
            time_elapsed += (abs(energy_consumption - MAX_ENERGY_CONS) + 1) * 100

        # Check velocity constraints
        if max_velocity > MAX_VELOCITY:
            time_elapsed += (max_velocity - MAX_VELOCITY + 1) * 100

        if min_velocity < 0:
            time_elapsed += (abs(min_velocity) + 1) * 100

        if max_acceleration > MAX_ACCELERATION:
            time_elapsed += (max_acceleration - MAX_ACCELERATION + 1) * 100

        if abs(min_acceleration) > MAX_ACCELERATION:
            time_elapsed += (abs(min_acceleration) - MAX_ACCELERATION + 1) * 100

        # Check the percentage difference constraint
        velocity_difference = abs(final_velocity - initial_velocity)
        time_elapsed += (abs(velocity_difference) + 1) * 100

        # If all constraints are satisfied, return the time elapsed
        return (
            time_elapsed
            if time_elapsed != float("inf") and time_elapsed >= 0
            else sys.float_info.max
        )
    except (ValueError, IndexError, OverflowError):
        # If parsing fails, return max float value as penalty
        return sys.float_info.max


# Initialization
expected_args = get_expected_argument_count()
npts = 20  # Number of points in the lattice (adjust based on problem size)
bounds = [(0, MAX_VELOCITY)] * expected_args  # Assuming bounds are known
mon = VerboseMonitor(10)

# Configure and solve using LatticeSolver
cube_root_npts = int(round(npts ** (1 / expected_args)))  # For 3D: npts ** (1/3)
nbins = (cube_root_npts,) * expected_args  # Adjust this based on your problem

# Initialization for target value
target_value = 0.01  # Set this to your desired target for the objective function

# Configure and solve using LatticeSolver
solver = SparsitySolver(expected_args, npts=npts)
solver.SetEvaluationMonitor(mon)
# Use VTR with the target value directly
solver.Solve(objective, termination=VTR(target_value), strategy=Best1Bin, disp=True)

res = solver.Solution()
print("Optimized Result:", res)
print("Objective Value:", objective(res))
output_cache.clear()
print(call_cli_program(res, ""))
