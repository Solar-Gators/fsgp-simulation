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
    return 5000 - energy_consumption


# Constraint function for initial and final velocity
def constraint_velocity(x):
    output = get_output(x)

    # in %
    acceptable_difference = 0.5

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

    # Constraint is satisfied if the difference is less than or equal to the acceptable limit
    return acceptable_difference - abs(velocity_difference)


# Define the constraints
con1 = {"type": "ineq", "fun": constraint_energy}
con2 = {"type": "ineq", "fun": constraint_velocity}

# Initial guess
x0 = [9, 0, -1, 2, -1, 0.55, -3.5, -1.4]

# Solve the optimization problem
res = minimize(
    objective,
    x0,
    constraints=[con1, con2],
    options={"disp": True, "maxiter": 50},
)

print(res.x)
print()

# Clear output_cache to ensure fresh output for final call
output_cache.clear()
print(call_cli_program(res.x))
