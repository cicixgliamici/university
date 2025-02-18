# Basic Arithmetic Operations
# Performs addition, subtraction, multiplication, and division

.text
.globl _start

_start:
    li x5, 10          # Load 10 into register x5
    li x6, 20          # Load 20 into register x6

    add x7, x5, x6     # x7 = x5 + x6 (10 + 20)
    sub x8, x6, x5     # x8 = x6 - x5 (20 - 10)
    mul x9, x5, x6     # x9 = x5 * x6 (10 * 20)
    div x10, x6, x5    # x10 = x6 / x5 (20 / 10)

    # Infinite loop to halt the program
    j .
