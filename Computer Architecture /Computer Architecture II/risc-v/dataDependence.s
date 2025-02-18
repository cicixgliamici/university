# Example illustrating RAW, WAW, and WAR hazards
# Assume x1 = 10, x2 = 20, and x3 = 30 before execution

# RAW: Read After Write (True Dependency)
# Instruction 2 depends on the result of Instruction 1
add x4, x1, x2     #Instruction 1: x4 = x1 + x2
sub x5, x4, x3     #Instruction 2: x5 = x4 - x3

# WAW: Write After Write (Output Dependency)
# Both instructions write to the same register (x6)
add x6, x1, x3     #Instruction 3: x6 = x1 + x3
mul x6, x2, x3     #Instruction 4: x6 = x2 * x3

# WAR: Write After Read (Anti-Dependency)
# Instruction 6 writes to a register (x4) that is read by Instruction 5
add x7, x4, x1     #Instruction 5: x7 = x4 + x1
add x4, x2, x3     #Instruction 6: x4 = x2 + x3
