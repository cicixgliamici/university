-- Advanced Example: Mathematical Operations in Ada
-- Demonstrates functions, procedures, parameters, loops, and exception handling

with Ada.Text_IO;           -- For text input/output
with Ada.Integer_Text_IO;   -- For integer input/output
use Ada.Text_IO, Ada.Integer_Text_IO; -- Combine use clauses

procedure Math_Operations is
    -- Function declaration (visible to entire procedure)
    function Factorial(n: in Integer) return Integer;
    
    -- Procedure declaration
    procedure Calculate_Fibonacci(n: in Integer; result: out Integer);
    
    -- Function implementation
    function Factorial(n: in Integer) return Integer is
    begin
        -- Recursive factorial calculation
        if n = 0 then
            return 1;
        else
            return n * Factorial(n - 1);
        end if;
    end Factorial;

    -- Procedure implementation
    procedure Calculate_Fibonacci(n: in Integer; result: out Integer) is
        a: Integer := 0;
        b: Integer := 1;
        temp: Integer;
    begin
        -- Iterative Fibonacci calculation
        if n = 0 then
            result := a;
            return;
        end if;
        
        for i in 2..n loop
            temp := a + b;
            a := b;
            b := temp;
        end loop;
        
        result := b;
    end Calculate_Fibonacci;

    -- Main variables
    number: Integer;
    fib_result: Integer;
begin
    -- User input with validation
    loop
        Put("Enter a positive integer (0-20): ");
        Get(number);
        exit when number >= 0 and number <= 20;
        Put_Line("Invalid input! Please try again.");
    end loop;

    -- Function call and output
    Put("Factorial of "); Put(number); Put(" is: ");
    Put(Factorial(number), Width => 1);
    New_Line;

    -- Procedure call and output
    Calculate_Fibonacci(number, fib_result);
    Put("Fibonacci("); Put(number); Put(") is: ");
    Put(fib_result, Width => 1);
    New_Line;

exception
    -- Basic exception handling
    when others =>
        Put_Line("An error occurred during calculation!");
end Math_Operations;
