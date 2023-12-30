<h1 align="center"><b>FSGP Simulation</b></h1>

<h4 align="center"><b>Go Race Simulation + Python Solver</b></h4>

<p align="center">

<a href="https://github.com/Solar-Gators/Pit-GUI/issues" target="_blank">
<img src="https://img.shields.io/github/issues/Solar-Gators/fsgp-simulation?style=flat-square" alt="issues"/>
</a>

<a href="https://github.com/Solar-Gators/Pit-GUI/pulls" target="_blank">
<img src="https://img.shields.io/github/issues-pr/Solar-Gators/fsgp-simulation?style=flat-square" alt="pull-requests"/>
</a>

</p>

## 
This project is comprised of 2 parts:
- **Race simulation CLI program** - written in Go to takes in a given race strategy (velocity at each point in the track) and uses physics constants of Sunrider to calculate energy consumed and time elapsed for a lap of the specified strategy.
- **Optimization solver** - written in Python using mystic to run thousands of iterations of the CLI simulation, utilizing gradient descent to find the optimal race strategy for a given track layout and target energy consumption.


https://github.com/Solar-Gators/fsgp-simulation/assets/26682594/74a2178b-8bd8-4b02-8d92-e0df9a3e3819
