# Essential Integrals for Calculus 1

## Basic Integrals

| Function          | Integral                            | Notes/Constraints                   |
|-------------------|-------------------------------------|-------------------------------------|
| $c$              | $cx + C$                           | $c$ is a constant                  |
| $x^n$            | $\frac{x^{n+1}}{n+1} + C$          | $n \neq -1$                        |
| $\frac{1}{x}$    | $\ln|x| + C$                       | $x \neq 0$                         |
| $e^x$            | $e^x + C$                          |                                     |
| $a^x$            | $\frac{a^x}{\ln(a)} + C$           | $a > 0$, $a \neq 1$                |

## Exponential and Logarithmic

| Function          | Integral                            | Notes/Constraints                   |
|-------------------|-------------------------------------|-------------------------------------|
| $\ln(x)$         | $x\ln(x) - x + C$                  | $x > 0$                            |
| $\log_a(x)$      | $\frac{x(\ln(x) - 1)}{\ln(a)} + C$ | $x > 0$, $a \neq 1$                |

## Trigonometric Functions

| Function          | Integral                            | Notes/Constraints                   |
|-------------------|-------------------------------------|-------------------------------------|
| $\sin(x)$        | $-\cos(x) + C$                     |                                     |
| $\cos(x)$        | $\sin(x) + C$                      |                                     |
| $\tan(x)$        | $-\ln|\cos(x)| + C$               | $x \neq \frac{\pi}{2} + k\pi$      |
| $\cot(x)$        | $\ln|\sin(x)| + C$                 | $x \neq k\pi$                      |
| $\sec(x)$        | $\ln|\sec(x) + \tan(x)| + C$       | $x \neq \frac{\pi}{2} + k\pi$      |
| $\csc(x)$        | $-\ln|\csc(x) + \cot(x)| + C$      | $x \neq k\pi$                      |
| $\sec^2(x)$     | $\tan(x) + C$                      |                                     |
| $\csc^2(x)$     | $-\cot(x) + C$                     |                                     |

## Inverse Trigonometric Functions

| Function              | Integral                              | Domain                     |
|-----------------------|---------------------------------------|----------------------------|
| $\frac{1}{\sqrt{1 - x^2}}$ | $\arcsin(x) + C$            | $-1 < x < 1$             |
| $\frac{1}{1 + x^2}$        | $\arctan(x) + C$            | $x \in \mathbb{R}$       |
| $\frac{1}{x\sqrt{x^2 - 1}}$| $\text{arcsec}|x| + C$      | $|x| > 1$               |

## Hyperbolic Functions

| Function           | Integral                           | Notes/Constraints          |
|--------------------|------------------------------------|----------------------------|
| $\sinh(x)$        | $\cosh(x) + C$                    |                            |
| $\cosh(x)$        | $\sinh(x) + C$                    |                            |
| $\tanh(x)$        | $\ln(\cosh(x)) + C$               |                            |
| $\coth(x)$        | $\ln|\sinh(x)| + C$               | $x \neq 0$                |
| $\text{sech}^2(x)$| $\tanh(x) + C$                    |                            |
| $\text{csch}^2(x)$| $-\coth(x) + C$                   | $x \neq 0$                |

## Inverse Hyperbolic Functions

| Function               | Integral                              | Domain                     |
|------------------------|---------------------------------------|----------------------------|
| $\frac{1}{\sqrt{x^2 + 1}}$ | $\text{arsinh}(x) + C$       | $x \in \mathbb{R}$        |
| $\frac{1}{\sqrt{x^2 - 1}}$ | $\text{arcosh}(x) + C$       | $x > 1$                  |
| $\frac{1}{1 - x^2}$        | $\text{artanh}(x) + C$       | $-1 < x < 1$             |
| $\frac{1}{1 - x^2}$        | $\text{arcoth}(x) + C$       | $|x| > 1$               |

## Integration Rules

| Rule Name           | Formula                                                                 |
|---------------------|-------------------------------------------------------------------------|
| Power Rule          | $\int x^n dx = \frac{x^{n+1}}{n+1} + C$ for $n \neq -1$                |
| Substitution        | $\int f(g(x))g'(x) dx = \int f(u) du$ where $u = g(x)$                 |
| Integration by Parts| $\int u \, dv = uv - \int v \, du$                                     |
| Linearity           | $\int [af(x) + bg(x)] dx = a\int f(x) dx + b\int g(x) dx$              |
