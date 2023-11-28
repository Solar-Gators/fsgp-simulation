import subprocess
import sys
from mystic.solvers import fmin
from mystic.monitors import VerboseMonitor
from mystic.penalty import quadratic_inequality

cli_program = "./strategy-simulation.exe"


# Call CLI program and return output
def call_cli_program(x):
    return subprocess.run(
        [cli_program] + list(map(str, x)), capture_output=True, text=True
    ).stdout


# Cache to store the output for the current x to avoid redundant CLI calls
output_cache = {}


# Function to get the output from cache or CLI call
def get_output(x):
    x_tuple = tuple(x)
    if x_tuple not in output_cache:
        output_cache[x_tuple] = call_cli_program(x)
    return output_cache[x_tuple]


# Objective function
def objective(x):
    output = get_output(x)
    try:
        time_elapsed = float(output.split("Time Elapsed (s):")[1].split("\n")[0])
        return (
            time_elapsed
            if time_elapsed != float("inf") or time_elapsed < 0
            else sys.float_info.max
        )
    except ValueError:
        return sys.float_info.max


# Constraint function for energy consumption
def constraint_energy(x):
    output = get_output(x)
    energy_consumption = float(
        output.split("Energy Consumption (W):")[1].split("\n")[0]
    )
    # Return positive if constraint is satisfied and negative if not
    return 5000 - energy_consumption


# Constraint function for initial and final velocity
def constraint_velocity(x):
    output = get_output(x)
    acceptable_difference = 0.5
    max_velocity = 40.0
    initial_velocity = float(output.split("Initial Velocity (m/s):")[1].split("\n")[0])
    final_velocity = float(output.split("Final Velocity (m/s):")[1].split("\n")[0])

    # Penalize if velocities are out of bounds
    if (
        initial_velocity < 0
        or initial_velocity > max_velocity
        or final_velocity < 0
        or final_velocity > max_velocity
    ):
        return -1

    # Penalize if the velocity difference percentage is too high
    if initial_velocity == final_velocity == 0:
        velocity_difference_percentage = 0
    else:
        velocity_difference_percentage = (
            abs(initial_velocity - final_velocity)
            / max(initial_velocity, final_velocity)
            * 100
        )
    if velocity_difference_percentage > acceptable_difference:
        return acceptable_difference - velocity_difference_percentage

    return 0  # No penalty if within acceptable bounds


# Combine penalties
def total_penalty(x):
    return penalty_energy(x) + penalty_velocity(x)


# Define the penalty functions using mystic's penalty method
@quadratic_inequality(constraint_energy)
def penalty_energy(x):
    return 0.0


@quadratic_inequality(constraint_velocity)
def penalty_velocity(x):
    return 0.0


# Combine penalties
def total_penalty(x):
    return penalty_energy(x) + penalty_velocity(x)


# Initial guess
x0 = [9, 0, -1, 2, -1, 0.55, -3.5, -1.4]

# Use VerboseMonitor to get the convergence information
mon = VerboseMonitor(10)


# Custom callback function
def custom_callback(x):
    y = objective(x)  # Compute the objective function value
    mon(x, y)  # Call the VerboseMonitor with both x and y


# Solve the optimization problem using the custom callback
res = fmin(
    objective,
    x0,
    penalty=total_penalty,  # Pass the combined penalty function
    disp=True,
    maxiter=100,
    callback=custom_callback,  # Use the custom callback function
)

print(res)

# Clear output_cache to ensure fresh output for final call
output_cache.clear()
print(call_cli_program(res))
