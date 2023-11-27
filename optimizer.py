from concurrent.futures import ThreadPoolExecutor
from scipy.optimize import minimize
import subprocess

subprocess.run(["go", "build", "."])
cli_program = "./strategy-simulation"
output = None

# Objective function
def objective(x):
    global output
    # Call your CLI program with the current parameters
    with ThreadPoolExecutor() as executor:
        futures = [executor.submit(subprocess.run, [cli_program] + list(map(str, x)), capture_output=True, text=True) for _ in range(4)]
    results = [future.result() for future in futures]
    # aggregate the results
    output = [result.stdout if result.stdout is not None else "" for result in results]

    # Parse the output and return the value to minimize
    # We will minimize the Time Elapsed (s)
    return sum(float(out.split("Time Elapsed (s):")[1].split("\n")[0]) for out in output)

# Constraint function
def constraint(x):
    global output

    # Parse the output and return the value for the constraint
    # Constraint is the Energy Consumption (W)
    if output:
        return 1000 - sum(float(out.split("Energy Consumption (W):")[1].split("\n")[0]) for out in output)
    else:
        return 1000

# Constraint function for initial and final velocity
def velocity_constraint(x):
    global output

    # Parse the output and get the initial and final velocities
    if output:
        initial_velocities = [float(out.split("Initial Velocity (m/s):")[1].split("\n")[0]) for out in output]
        final_velocities = [float(out.split("Final Velocity (m/s):")[1].split("\n")[0]) for out in output]

        # Calculate the difference in percentage
        velocity_differences = [abs(initial_velocity - final_velocity) / initial_velocity * 100 for initial_velocity, final_velocity in zip(initial_velocities, final_velocities)]

        # Velocity difference should be within 0.5%
        return 0.5 - max(velocity_differences)
    else:
        return 0.5

# Initial guess
x0 = [1]*8

# Define the constraints
con1 = {'type': 'ineq', 'fun': constraint}
con2 = {'type': 'ineq', 'fun': velocity_constraint}

# Solve the optimization problem
res = minimize(objective, x0, method='SLSQP', constraints=[con1, con2], options={'disp': True, 'maxiter': 100})

# Print the optimal solution
print(res.x)
