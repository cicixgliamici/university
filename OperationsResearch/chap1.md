# Summary of Core Formulas in Operations Research

## Introduction

Operations Research (OR) applies mathematical models to support decision-making in complex systems. The field relies on various techniques—including linear programming, integer programming, and queueing theory—to optimize performance and resource allocation. This summary outlines several key formulas and modeling approaches common to OR.

## Linear Programming

A classic example in production planning is modeled as a linear programming (LP) problem. Consider the following formulation:

**Objective Function:**

$$
\max z = 30x_1 + 50x_2
$$

**Constraints:**

$$
\begin{aligned}
x_1 &\leq 4 \\
2x_2 &\leq 10 \\
3x_1 + 2x_2 &\leq 18 \\
x_1, \; x_2 &\geq 0
\end{aligned}
$$

This LP model optimizes profit while ensuring that production does not exceed available capacities.

## Queueing Theory

For a single-server queue (commonly referred to as an M/M/1 queue) with arrival rate \(\lambda\) and service rate \(\mu\) (where \(\lambda < \mu\)), the average waiting time \(W\) in the system is given by:

$$
W = \frac{\lambda}{\mu (\mu - \lambda)}
$$

This formula helps in assessing the performance of service processes under steady-state conditions.

## Assignment Problem Model

The assignment problem involves assigning \(n\) tasks to \(n\) candidates with the goal of minimizing total cost or time. It is modeled with decision variables \(x_{ij}\) defined as follows:

$$
x_{ij} =
\begin{cases}
1, & \text{if candidate } i \text{ is assigned to task } j \\
0, & \text{otherwise}
\end{cases}
$$

**Objective Function:**

$$
\min \sum_{i=1}^{n} \sum_{j=1}^{n} t_{ij} x_{ij}
$$

**Constraints:**

$$
\begin{aligned}
\sum_{i=1}^{n} x_{ij} &= 1 \quad \text{for each task } j \\
\sum_{j=1}^{n} x_{ij} &= 1 \quad \text{for each candidate } i \\
x_{ij} &\in \{0, 1\} \quad \text{for all } i, j
\end{aligned}
$$

This ensures that each task is assigned to exactly one candidate and vice versa, optimizing the overall assignment based on the cost matrix \(T = (t_{ij})\).

## Modeling Framework

The general process of developing an OR model typically involves these steps:

- **Problem Definition:** Clearly identify the decision-making issue.
- **System Analysis:** Gather data such as arrival rates (\(\lambda\)) and service rates (\(\mu\)) for queueing systems, or resource constraints for production systems.
- **Model Formulation:** Develop a mathematical representation using appropriate techniques (e.g., LP, ILP, or simulation).
- **Solution and Testing:** Apply algorithms (e.g., Simplex, branch-and-bound) to solve the model and validate predictions against real data.
- **Implementation:** Execute the solution in practice and monitor performance for necessary adjustments.
