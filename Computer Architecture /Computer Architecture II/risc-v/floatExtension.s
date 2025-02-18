# RISC-V Floating Point Extension: Educational Example
# This RISC-V assembly example demonstrates the use of the floating point (F) extension. The code performs basic floating point operations such as addition, subtraction, multiplication, division, as well as conversions between integer and floating point values. These operations are crucial for scientific computing, graphics, signal processing, and any application requiring real-number computations.
# Below is the complete code with detailed comments:


    .section .data
# Define floating point constants and an integer value in the data section.
float1:    .float 1.5       # Single-precision float constant (1.5)
float2:    .float 2.5       # Single-precision float constant (2.5)
int_val:   .word 10         # Integer value (10)

    .section .text
    .global main
main:
    # -----------------------------------------------------------
    # Load floating point constants into floating point registers.
    # -----------------------------------------------------------
    flw    f0, float1       # Load the value 1.5 from memory into register f0.
    flw    f1, float2       # Load the value 2.5 from memory into register f1.

    # -----------------------------------------------------------
    # Floating Point Addition: f2 = f0 + f1
    # Expected result: 1.5 + 2.5 = 4.0
    # -----------------------------------------------------------
    fadd.s f2, f0, f1       # Perform single-precision addition and store result in f2.

    # -----------------------------------------------------------
    # Floating Point Multiplication: f3 = f0 * f1
    # Expected result: 1.5 * 2.5 = 3.75
    # -----------------------------------------------------------
    fmul.s f3, f0, f1       # Multiply f0 and f1, result in f3.

    # -----------------------------------------------------------
    # Integer to Floating Point Conversion:
    # Convert the integer in memory (int_val) to a float.
    # -----------------------------------------------------------
    lw     t0, int_val      # Load integer value 10 from memory into register t0.
    fcvt.s.w f4, t0        # Convert the integer in t0 to a single-precision float in f4 (10.0).

    # -----------------------------------------------------------
    # Floating Point Subtraction: f5 = f4 - f0
    # Expected result: 10.0 - 1.5 = 8.5
    # -----------------------------------------------------------
    fsub.s f5, f4, f0       # Subtract f0 from f4, result stored in f5.

    # -----------------------------------------------------------
    # Floating Point to Integer Conversion:
    # Convert the float in f3 (3.75) to an integer (rounding toward zero).
    # Expected result: 3
    # -----------------------------------------------------------
    fcvt.w.s t1, f3        # Convert f3 to an integer and store it in t1.

    # -----------------------------------------------------------
    # Floating Point Division: f6 = f4 / f1
    # Expected result: 10.0 / 2.5 = 4.0
    # -----------------------------------------------------------
    fdiv.s f6, f4, f1       # Divide f4 by f1, result stored in f6.

    # -----------------------------------------------------------
    # Floating Point Comparison:
    # Check if f0 is less than f1.
    # The instruction sets t2 to 1 if (f0 < f1) is true, else 0.
    # Expected result: 1 (since 1.5 < 2.5)
    # -----------------------------------------------------------
    flt.s  t2, f0, f1       # Compare f0 and f1; result in t2.

    # -----------------------------------------------------------
    # Summary of Computation:
    # f2 : 4.0   (Addition result)
    # f3 : 3.75  (Multiplication result)
    # f4 : 10.0  (Integer to float conversion result)
    # f5 : 8.5   (Subtraction result)
    # t1 : 3     (Float to integer conversion result)
    # f6 : 4.0   (Division result)
    # t2 : 1     (Comparison result: true)
    #
    # In a complete system, these values might be stored, used in further calculations,
    # or printed out using system calls. For educational purposes, they remain in registers.
    # -----------------------------------------------------------

    # -----------------------------------------------------------
    # Exit the program.
    # For RISC-V Linux, use the ECALL with the exit syscall.
    # -----------------------------------------------------------
    li     a7, 93         # Load syscall number for exit (93) into a7.
    li     a0, 0          # Set exit code 0 (successful termination) in a0.
    ecall                 # Make the system call to exit.


# Importance of the Floating Point Extension
# Precision and Range: The F extension provides hardware support for single-precision floating point arithmetic, enabling calculations with fractional numbers and a wide range of magnitudes.
# Performance: Hardware acceleration of floating point operations significantly speeds up computation compared to software emulation.
# Scientific and Engineering Applications: Many algorithms in science, engineering, and graphics rely on efficient floating point arithmetic.
# Standard Compliance: The RISC-V F extension follows IEEE 754 standards, ensuring consistency and portability of floating point computations across different platforms.
