from scipy.optimize import minimize
import subprocess
import sys

subprocess.run(["go", "build", "."])
cli_program = "./strategy-simulation"
output = None

# Objective function
def objective(x):
    global output
    # Call your CLI program with the current parameters
    result = subprocess.run([cli_program] + list(map(str, x)), capture_output=True, text=True)
    output = result.stdout if result.stdout is not None else ""

    # Parse the output and return the value to minimize
    # If the output is infinity, return the maximum float value
    try:
        time_elapsed = float(output.split("Time Elapsed (s):")[1].split("\n")[0])
        return time_elapsed if time_elapsed != float('inf') else sys.float_info.max
    except ValueError:
        # Handle the case where the conversion to float fails
        return sys.float_info.max

# Constraint function
def constraint(x):
    global output

    # Parse the output and return the value for the constraint
    # Constraint is the Energy Consumption (W)
    if output:
        return 1000 - float(output.split("Energy Consumption (W):")[1].split("\n")[0])
    else:
        return 1000

# Constraint function for initial and final velocity
def velocity_constraint(x):
    global output
    acceptable_difference = 0.5  # 0.5% or an absolute threshold for very small velocities

    # Parse the output and get the initial and final velocities
    if output:
        initial_velocity = float(output.split("Initial Velocity (m/s):")[1].split("\n")[0])
        final_velocity = float(output.split("Final Velocity (m/s):")[1].split("\n")[0])

        # Check for initial velocity close to zero
        if abs(initial_velocity) < 1e-8:  # A threshold for considering velocity as zero
            # Use an absolute difference for very small velocities
            velocity_difference = abs(final_velocity)
        else:
            # Calculate the difference in percentage for non-zero initial velocities
            velocity_difference = abs(initial_velocity - final_velocity) / initial_velocity * 100

        # Velocity difference should be within the acceptable threshold
        return acceptable_difference - velocity_difference
    else:
        return acceptable_difference

# Initial guess
x0 = [0.00001]*8

# Define the constraints
con1 = {'type': 'ineq', 'fun': constraint}
con2 = {'type': 'ineq', 'fun': velocity_constraint}

# Solve the optimization problem
res = minimize(objective, x0, method='SLSQP', constraints=[con1, con2], options={'disp': True, 'maxiter': 100})

# Print the optimal solution
print(res.x)

