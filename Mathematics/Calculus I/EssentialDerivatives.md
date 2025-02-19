
# Essential Derivatives for Calculus 1

## Basic Derivatives

| Function          | Derivative                          | Notes/Constraints                   |
|-------------------|-------------------------------------|-------------------------------------|
| $c$            | $0$                              | $c$ is a constant                 |
| $x^n$          | $nx^{n-1}$                       | $n \in \mathbb{R}$                |
| $e^x$          | $e^x$                            |                                     |
| $a^x$          | $a^x \ln(a)$                     | $a > 0$                          |
| $\ln(x)$       | $\frac{1}{x}$                    | $x > 0$                          |
| $\log_a(x)$    | $\frac{1}{x \ln(a)}$             | $x > 0, \ a \neq 1$              |
| $\sin(x)$      | $\cos(x)$                        |                                     |
| $\cos(x)$      | $-\sin(x)$                       |                                     |
| $\tan(x)$      | $\sec^2(x)$                      | $x \neq \frac{\pi}{2} + k\pi$      |
| $\cot(x)$      | $-\csc^2(x)$                     | $x \neq k\pi$                     |
| $\sec(x)$      | $\sec(x)\tan(x)$                  | $x \neq \frac{\pi}{2} + k\pi$      |
| $\csc(x)$      | $-\csc(x)\cot(x)$                 | $x \neq k\pi$                     |

## Inverse Trigonometric Functions

| Function              | Derivative                              | Domain                     |
|-----------------------|-----------------------------------------|----------------------------|
| $\arcsin(x)$       | $\frac{1}{\sqrt{1 - x^2}}$            | $-1 < x < 1$             |
| $\arccos(x)$       | $-\frac{1}{\sqrt{1 - x^2}}$           | $-1 < x < 1$             |
| $\arctan(x)$       | $\frac{1}{1 + x^2}$                   | $x \in \mathbb{R}$       |
| $\text{arccot}(x)$ | $-\frac{1}{1 + x^2}$                   | $x \in \mathbb{R}$       |
| $\text{arcsec}(x)$ | $\frac{1}{|x|\sqrt{x^2 - 1}}$         | $|x| > 1$               |
| $\text{arccsc}(x)$ | $-\frac{1}{|x|\sqrt{x^2 - 1}}$        | $|x| > 1$               |

## Hyperbolic Functions

| Function           | Derivative                           |
|--------------------|--------------------------------------|
| $\sinh(x)$      | $\cosh(x)$                        |
| $\cosh(x)$      | $\sinh(x)$                        |
| $\tanh(x)$      | $\text{sech}^2(x)$                 |
| $\coth(x)$      | $-\text{csch}^2(x)$                |
| $\text{sech}(x)$| $-\text{sech}(x)\tanh(x)$           |
| $\text{csch}(x)$| $-\text{csch}(x)\coth(x)$           |

## Inverse Hyperbolic Functions

| Function               | Derivative                              | Domain                     |
|------------------------|-----------------------------------------|----------------------------|
| $\text{arsinh}(x)$   | $\frac{1}{\sqrt{x^2 + 1}}$             | $x \in \mathbb{R}$        |
| $\text{arcosh}(x)$   | $\frac{1}{\sqrt{x^2 - 1}}$             | $x > 1$                 |
| $\text{artanh}(x)$   | $\frac{1}{1 - x^2}$                    | $-1 < x < 1$             |
| $\text{arcoth}(x)$   | $\frac{1}{1 - x^2}$                    | $|x| > 1$               |

## Differentiation Rules

| Rule Name       | Formula                                                  |
|-----------------|----------------------------------------------------------|
| Product Rule    | $ (fg)' = f'g + fg' $                                  |
| Quotient Rule   | $ \left(\frac{f}{g}\right)' = \frac{f'g - fg'}{g^2} $   |
| Chain Rule      | $ \frac{d}{dx}[f(g(x))] = f'(g(x)) \cdot g'(x) $         |
