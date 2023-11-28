from scipy.optimize import minimize
import subprocess
import sys

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
    # Convert x to a tuple to use it as a dictionary key (lists are not hashable)
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
    # Return non-negative if constraint is satisfied, negative otherwise
    if energy_consumption < 0:
        return -1

    if energy_consumption > 5000:
        return -1
    return 0


# Constraint function for initial and final velocity
def constraint_velocity(x):
    output = get_output(x)

    # in %
    acceptable_difference = 2.0

    max_velocity = 40.0  # Maximum allowed velocity
    initial_velocity = float(output.split("Initial Velocity (m/s):")[1].split("\n")[0])
    final_velocity = float(output.split("Final Velocity (m/s):")[1].split("\n")[0])

    # Check if either velocity is outside the acceptable range [0, max_velocity]
    if initial_velocity < 0 or initial_velocity > max_velocity:
        return -1
    if final_velocity < 0 or final_velocity > max_velocity:
        return -1

    # Calculate the percentage difference between the initial and final velocities
    if initial_velocity == final_velocity == 0:
        velocity_difference = 0
    else:
        # Otherwise, calculate the difference as a percentage
        velocity_difference = (
            abs(initial_velocity - final_velocity)
            / max(initial_velocity, final_velocity)
            * 100
        )

    # Check if the difference is within the acceptable limit
    if velocity_difference > acceptable_difference:
        return -1

    # If all constraints are satisfied, return 0 or a positive value
    return 0


# Initial guess
x0 = [9, 0, -1, 2, -1, 0.55, -3.5, -1.4]

# Define the constraints
con1 = {"type": "eq", "fun": constraint_energy}
con2 = {"type": "eq", "fun": constraint_velocity}

# Solve the optimization problem
res = minimize(
    objective,
    x0,
    method="SLSQP",
    constraints=[con1, con2],
    options={"disp": True, "maxiter": 20},
)

print(res.x)
print()

# Clear output_cache to ensure fresh output for final call
output_cache.clear()
print(call_cli_program(res.x))
