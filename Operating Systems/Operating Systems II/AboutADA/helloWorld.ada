-- Simple "Hello, World!" program in Ada
-- Demonstrates basic program structure, package usage, and I/O operations

-- Ada programs are typically organized using packages
-- We need to explicitly include required packages
with Ada.Text_IO;           -- Standard package for text input/output operations
use Ada.Text_IO;             -- Allows direct access to package contents without prefix

-- Every Ada program has at least one procedure
-- Procedure name must match filename (Hello_World.adb for this example)
procedure Hello_World is     -- 'is' begins the procedure/function declaration
-- Declaration section (variables, constants, etc.) would go here
-- We don't need any for this simple example

begin  -- Marks the start of executable statements
    -- Use Put_Line procedure from Ada.Text_IO package
    -- Automatically adds a newline after output
    -- Without 'use' clause, we'd need to write Ada.Text_IO.Put_Line(...)
    Put_Line("Hello, World!");  -- String literals are in double quotes
    
    -- Note: Ada is case-insensitive, but conventional formatting uses:
    -- - Lowercase for reserved words
    -- - Mixed case for identifiers
    -- - Underscores for word separation
    
end Hello_World;  -- End of procedure with repeated name for clarity


-- Package Inclusion: with clause imports packages
-- Namespace Management: use clause avoids repetitive package prefixes
-- Procedure Structure: Mandatory main procedure format
-- Case Insensitivity: Language convention despite case flexibility
-- String Handling: Double-quoted literals and automatic newline
-- File Naming: Relationship between procedure name and filename
-- This demonstrates Ada's emphasis on explicit declarations, strong typing, and modular design while maintaining readability.
