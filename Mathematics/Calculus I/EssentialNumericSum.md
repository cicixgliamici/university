# Essential Numerical Series

## Common Series and Their Sums

| Series Name          | General Form                                  | Sum (if it exists)                       | Constraints                 |
|----------------------|---------------------------------------------|-----------------------------------------|----------------------------|
| Geometric Series    | $\sum_{n=0}^{\infty} ar^n$                   | $\frac{a}{1 - r}$                      | $\|r\| < 1$                   |
| Harmonic Series     | $\sum_{n=1}^{\infty} \frac{1}{n}$            | Divergent                              |                              |
| p-Series           | $\sum_{n=1}^{\infty} \frac{1}{n^p}$          | Convergent if $p > 1$, Divergent if $p \leq 1$ | $p > 0$                      |
| Alternating Harmonic Series | $\sum_{n=1}^{\infty} \frac{(-1)^{n+1}}{n}$ | $\ln(2)$                              |                              |
| Telescoping Series  | $\sum_{n=1}^{\infty} (b_n - b_{n+1})$        | Depends on $\lim_{n \to \infty} b_n$  |                              |

## Convergence Tests

| Test Name            | Condition for Convergence                      |
|----------------------|------------------------------------------------|
| Ratio Test          | $\lim_{n \to \infty} \left\| \frac{a_{n+1}}{a_n} \right\| < 1$ |
| Root Test           | $\limsup_{n \to \infty} \sqrt[n]{\|a_n\|} < 1$    |
| Integral Test       | If $f(n)$ is positive, decreasing, and $\int_1^{\infty} f(x) dx$ converges, then so does $\sum a_n$. |
| Comparison Test     | If $0 \leq a_n \leq b_n$ and $\sum b_n$ converges, then $\sum a_n$ converges. |
| Limit Comparison Test | If $\lim_{n \to \infty} \frac{a_n}{b_n} = c > 0$, and $\sum b_n$ converges, then $\sum a_n$ converges. |
| Alternating Series Test | If $a_n$ is decreasing and $\lim_{n \to \infty} a_n = 0$, then $\sum (-1)^n a_n$ converges. |

## Power Series

| Series Name         | General Form                                       | Radius of Convergence |
|--------------------|-------------------------------------------------|----------------------|
| General Power Series | $\sum_{n=0}^{\infty} c_n (x - a)^n$            | $R = \limsup_{n \to \infty} |c_n|^{1/n}$ |
| Taylor Series      | $\sum_{n=0}^{\infty} \frac{f^{(n)}(a)}{n!} (x-a)^n$ | Depends on function   |
| Maclaurin Series   | $\sum_{n=0}^{\infty} \frac{f^{(n)}(0)}{n!} x^n$  | Special case of Taylor Series |

## Special Taylor Series Expansions

| Function             | Expansion (Maclaurin Series)                   | Interval of Convergence |
|----------------------|----------------------------------------------|------------------------|
| $e^x$               | $\sum_{n=0}^{\infty} \frac{x^n}{n!}$         | $(-\infty, \infty)$  |
| $\sin(x)$           | $\sum_{n=0}^{\infty} \frac{(-1)^n x^{2n+1}}{(2n+1)!}$ | $(-\infty, \infty)$  |
| $\cos(x)$           | $\sum_{n=0}^{\infty} \frac{(-1)^n x^{2n}}{(2n)!}$ | $(-\infty, \infty)$  |
| $\frac{1}{1-x}$     | $\sum_{n=0}^{\infty} x^n$                    | $\|x\| < 1$              |
| $\ln(1+x)$         | $\sum_{n=1}^{\infty} (-1)^{n+1} \frac{x^n}{n}$ | $-1 < x \leq 1$       |


