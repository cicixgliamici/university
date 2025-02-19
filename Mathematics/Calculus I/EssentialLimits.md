# Essential Limits for Calculus 1

## Basic Limits
| Function              | Limit Expression                   | Result          | Notes                     |
|-----------------------|------------------------------------|-----------------|---------------------------|
| Constant              | `$$\lim_{x \to a} c$$`            | `$$c$$`         | `$$c$$` is constant       |
| Identity Function     | `$$\lim_{x \to a} x$$`            | `$$a$$`         |                           |
| Linear Function       | `$$\lim_{x \to a} (mx + b)$$`     | `$$ma + b$$`    |                           |
| Polynomial            | `$$\lim_{x \to a} p(x)$$`         | `$$p(a)$$`      | For polynomial functions  |

## Exponential & Logarithmic Limits
| Limit Expression                                  | Result            | Notes                          |
|---------------------------------------------------|-------------------|--------------------------------|
| `$$\lim_{x \to 0} \frac{e^x - 1}{x}$$`            | `$$1$$`           |                                |
| `$$\lim_{x \to 0} \frac{a^x - 1}{x}$$`            | `$$\ln a$$`       | `$$a > 0$$`                    |
| `$$\lim_{x \to 0} \frac{\ln(1 + x)}{x}$$`         | `$$1$$`           |                                |
| `$$\lim_{x \to \infty} \left(1 + \frac{1}{x}\right)^x$$` | `$$e$$`     |                                |
| `$$\lim_{x \to 0} (1 + x)^{1/x}$$`                | `$$e$$`           |                                |

## Trigonometric Limits
| Limit Expression                      | Result    | Notes                          |
|---------------------------------------|-----------|--------------------------------|
| `$$\lim_{x \to 0} \frac{\sin x}{x}$$` | `$$1$$`   | Angles in radians             |
| `$$\lim_{x \to 0} \frac{\tan x}{x}$$` | `$$1$$`  |                                |
| `$$\lim_{x \to 0} \frac{1 - \cos x}{x}$$` | `$$0$$` |                        |
| `$$\lim_{x \to 0} \frac{1 - \cos x}{x^2}$$` | `$$\frac{1}{2}$$` |                |

## Limits at Infinity
| Limit Expression                          | Result          | Notes                      |
|-------------------------------------------|-----------------|----------------------------|
| `$$\lim_{x \to \infty} \frac{1}{x}$$`     | `$$0$$`         |                            |
| `$$\lim_{x \to \infty} x^{1/x}$$`         | `$$1$$`         |                            |
| `$$\lim_{x \to \infty} \frac{\ln x}{x}$$` | `$$0$$`         |                            |
| `$$\lim_{x \to 0^+} x \ln x$$`            | `$$0$$`         |                            |

## Special Limits
| Limit Expression                              | Result          | Notes                          |
|-----------------------------------------------|-----------------|--------------------------------|
| `$$\lim_{x \to 0} \frac{\arcsin x}{x}$$`      | `$$1$$`         |                                |
| `$$\lim_{x \to 0} \frac{\arctan x}{x}$$`      | `$$1$$`         |                                |
| `$$\lim_{x \to 0} \frac{e^{kx} - 1}{x}$$`     | `$$k$$`         | `$$k$$` constant              |
| `$$\lim_{x \to \infty} \left(1 + \frac{k}{x}\right)^{mx}$$` | `$$e^{mk}$$` |                                |

## Limit Laws
| Law Name          | Mathematical Expression                          |
|-------------------|--------------------------------------------------|
| Sum Rule          | `$$\lim_{x \to a} [f(x) + g(x)] = L + M$$`       |
| Difference Rule   | `$$\lim_{x \to a} [f(x) - g(x)] = L - M$$`       |
| Product Rule      | `$$\lim_{x \to a} [f(x) \cdot g(x)] = L \cdot M$$` |
| Quotient Rule
