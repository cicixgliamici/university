# Program: Load and Store in RISC-V
# Demonstrates loading and storing memory values.

.data
value1: .word 15       # Memory location with value 15
value2: .word 25       # Memory location with value 25
result: .word 0        # Memory location to store result

.text
.globl _start

_start:
    la x1, value1      # Load address of value1 into x1
    lw x2, 0(x1)       # Load word at address in x1 into x2 (x2 = 15)

    la x3, value2      # Load address of value2 into x3
    lw x4, 0(x3)       # Load word at address in x3 into x4 (x4 = 25)

    add x5, x2, x4     # x5 = x2 + x4 (x5 = 15 + 25 = 40)

    la x6, result      # Load address of result into x6
    sw x5, 0(x6)       # Store x5 value into memory at address x6

    # Infinite loop to stop the program
    j .
