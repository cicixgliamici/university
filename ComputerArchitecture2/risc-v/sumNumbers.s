# Program: Sum of first n natural numbers
# Demonstrates:
#   - While loop
#   - Function call
#   - Load/Store instructions
# Registers:
#   x10: n (input)
#   x11: sum (output)

.data
n:      .word 5         # Input value: n = 5
result: .word 0         # Result storage

.text
.globl _start

# Main program
_start:
    la x1, n                # Load the address of n into x1
    lw x10, 0(x1)           # Load the value of n into x10
    li x11, 0               # Initialize sum = 0 in x11
    jal x1, sum_numbers     # Call sum_numbers function
    la x2, result           # Load the address of result into x2
    sw x11, 0(x2)           # Store the result (x11) into memory

    # Infinite loop to terminate program
    j .

# Function: sum_numbers
# Input: x10 (n)
# Output: x11 (sum)
sum_numbers:
    li x12, 0           # Initialize counter i = 0 (x12)
    li x11, 0           # Initialize sum = 0 (x11)

while_loop:
    bge x12, x10, end_loop  # Exit loop if i >= n
    add x11, x11, x12       # sum += i
    addi x12, x12, 1        # i += 1
    j while_loop            # Repeat the loop

end_loop:
    ret                 # Return to caller
