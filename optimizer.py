import subprocess
import platform
import sys
import os
from mystic.solvers import *
from mystic.monitors import *
from mystic.constraints import *
from mystic.symbolic import *

NUM_ARGUMENTS = 8
CALLS_BETWEEN_IMAGE = 20

try:
    subprocess.run(["go", "build", "."])
except:
    print("Ensure Go is installed! Using binaries...\n")

if platform.system() == "Windows":
    cli_program = "./strategy-simulation.exe"
else:
    cli_program = "./strategy-simulation"


# Call CLI program and return output
def call_cli_program(x, endArg):
    return subprocess.run(
        [cli_program] + list(map(str, x)) + [str(endArg)],
        capture_output=True,
        text=True,
    ).stdout


# Cache to store the output for the current x to avoid redundant CLI calls
output_cache = {}


i = 0


# Function to get the output from cache or CLI call
def get_output(x):
    global i
    if i % CALLS_BETWEEN_IMAGE == 0:
        autoEndArg = ""
    else:
        autoEndArg = "none"
    i += 1
    x_tuple = tuple(x)
    if x_tuple not in output_cache:
        output_cache[x_tuple] = call_cli_program(x, autoEndArg)
    return output_cache[x_tuple]


# Objective function with constraints
def objective(x):
    # in %
    acceptable_difference = 0.05
    max_velocity = 40.0

    output = get_output(x)
    try:
        # Parse the output for the required values
        time_elapsed = float(output.split("Time Elapsed (s):")[1].split("\n")[0])
        energy_consumption = float(
            output.split("Energy Consumption (W):")[1].split("\n")[0]
        )
        initial_velocity = float(
            output.split("Initial Velocity (m/s):")[1].split("\n")[0]
        )
        final_velocity = float(output.split("Final Velocity (m/s):")[1].split("\n")[0])

        # Check energy consumption constraint
        if energy_consumption > 1300 or energy_consumption < 0:
            return sys.float_info.max

        # Check velocity constraints
        if not (0 < initial_velocity < max_velocity) or not (
            0 < final_velocity < max_velocity
        ):
            return sys.float_info.max

        # Check the percentage difference constraint
        velocity_difference = abs(final_velocity - initial_velocity)
        if (
            velocity_difference / max(initial_velocity, final_velocity)
            > acceptable_difference / 100
        ):
            return sys.float_info.max

        # If all constraints are satisfied, return the time elapsed
        return (
            time_elapsed
            if time_elapsed != float("inf") and time_elapsed >= 0
            else sys.float_info.max
        )
    except (ValueError, IndexError):
        # If parsing fails, return max float value as penalty
        return sys.float_info.max


mon = VerboseMonitor(10)


def custom_callback(x):
    y = objective(x)
    mon(x, y)


# Initial guess
x0 = [0] * NUM_ARGUMENTS
x0[0] = 1

# Solve the optimization problem using the constraints
res = fmin_powell(
    objective,
    x0,
    disp=True,
    maxiter=2000,
    callback=custom_callback,
)

print(res)
output_cache.clear()
print(call_cli_program(res, ""))
