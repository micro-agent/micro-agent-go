# Complete Python Getting Started Course

## Table of Contents

1. [Introduction to Python](#1-introduction-to-python)
2. [Getting Started](#2-getting-started)
3. [Basic Syntax and Concepts](#3-basic-syntax-and-concepts)
4. [Variables and Data Types](#4-variables-and-data-types)
5. [Operators](#5-operators)
6. [Control Flow](#6-control-flow)
7. [Functions](#7-functions)
8. [Data Structures](#8-data-structures)
9. [Strings](#9-strings)
10. [File Input/Output](#10-file-inputoutput)
11. [Error Handling](#11-error-handling)
12. [Modules and Packages](#12-modules-and-packages)
13. [Object-Oriented Programming Basics](#13-object-oriented-programming-basics)
14. [Standard Library Overview](#14-standard-library-overview)
15. [Best Practices](#15-best-practices)
16. [Exercises and Projects](#16-exercises-and-projects)

---

## 1. Introduction to Python

### What is Python?

Python is a high-level, interpreted programming language created by Guido van Rossum and first released in 1991. It emphasizes code readability and simplicity, making it an excellent choice for beginners.

### Key Features

- **Easy to learn and read**
- **Interpreted language** (no compilation needed)
- **Cross-platform** (Windows, macOS, Linux)
- **Extensive standard library**
- **Large community and ecosystem**
- **Dynamically typed**
- **Object-oriented and functional programming support**

### Why Choose Python?

- Beginner-friendly syntax
- Rapid development and prototyping
- Versatile applications (web, data science, AI, automation)
- Strong community support
- Extensive third-party libraries
- Great for scripting and automation

### Use Cases

- Web development (Django, Flask)
- Data science and analytics
- Machine learning and AI
- Automation and scripting
- Scientific computing
- Game development
- Desktop applications

---

## 2. Getting Started

### Installation

#### On Windows
1. Visit https://python.org/downloads/
2. Download the latest Python version
3. Run the installer and check "Add Python to PATH"
4. Verify installation: open Command Prompt and type `python --version`

#### On macOS
```bash
# Using Homebrew (recommended)
brew install python

# Or download from python.org
```

#### On Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install python3 python3-pip

# CentOS/RHEL
sudo yum install python3 python3-pip
```

### Verify Installation
```bash
python --version
python3 --version
pip --version
```

### Setting up Development Environment

#### Python IDLE (built-in)
- Comes with Python installation
- Good for beginners and simple scripts

#### Visual Studio Code
1. Install VS Code
2. Install Python extension
3. Configure Python interpreter

#### PyCharm
- Professional IDE for Python
- Community edition is free

#### Jupyter Notebook
```bash
pip install jupyter
jupyter notebook
```

### Running Python Code

#### Interactive Mode (REPL)
```bash
python
>>> print("Hello, World!")
Hello, World!
>>> exit()
```

#### Script Mode
Create a file `hello.py`:
```python
print("Hello, World!")
```

Run it:
```bash
python hello.py
```

---

## 3. Basic Syntax and Concepts

### Your First Python Program

```python
print("Hello, World!")
```

### Python Philosophy (The Zen of Python)
```python
import this
```

### Code Structure

```python
# This is a comment
print("Hello, World!")  # This is also a comment

"""
This is a
multi-line comment
or docstring
"""

# Python uses indentation to define code blocks
if True:
    print("This is indented")
    print("This too")
print("This is not indented")
```

### Indentation Rules

- Python uses indentation (spaces or tabs) to define code blocks
- Standard is 4 spaces per indentation level
- Be consistent throughout your code

```python
# Correct indentation
if True:
    print("Correct")
    if True:
        print("Nested correct")

# Incorrect indentation (will cause IndentationError)
# if True:
# print("Wrong")
```

### Case Sensitivity

```python
# Python is case-sensitive
name = "John"
Name = "Jane"
print(name)  # Outputs: John
print(Name)  # Outputs: Jane
```

### Statements and Expressions

```python
# Statement: performs an action
print("Hello")

# Expression: evaluates to a value
result = 5 + 3
```

---

## 4. Variables and Data Types

### Variables

```python
# Variable assignment
name = "Alice"
age = 25
height = 5.6
is_student = True

# Multiple assignment
x, y, z = 1, 2, 3
a = b = c = 10

# Variable naming rules
# Valid names
first_name = "John"
age_2 = 25
_private = "hidden"

# Invalid names (will cause errors)
# 2age = 25        # Can't start with number
# first-name = ""  # Can't use hyphens
# class = "test"   # Can't use keywords
```

### Data Types

#### Numeric Types

```python
# Integer
age = 25
year = 2024

# Float
height = 5.9
price = 19.99

# Complex (rarely used)
complex_num = 3 + 4j

# Type checking
print(type(age))    # <class 'int'>
print(type(height)) # <class 'float'>
```

#### String Type

```python
# Single quotes
name = 'Alice'

# Double quotes
message = "Hello, World!"

# Triple quotes (multi-line)
long_text = """This is a
multi-line
string"""

# Escape characters
escaped = "She said \"Hello\""
newline = "Line 1\nLine 2"
tab = "Column1\tColumn2"
```

#### Boolean Type

```python
is_valid = True
is_complete = False

# Boolean from expressions
is_adult = age >= 18
has_name = len(name) > 0
```

#### None Type

```python
# Represents absence of value
result = None
optional_value = None

# Checking for None
if result is None:
    print("No result")
```

### Type Conversion

```python
# Implicit conversion
result = 5 + 2.0  # Result is 7.0 (float)

# Explicit conversion
age_str = "25"
age_int = int(age_str)      # String to integer
age_float = float(age_str)  # String to float

price = 19.99
price_str = str(price)      # Float to string
price_int = int(price)      # Float to integer (truncates)

# Boolean conversion
bool_from_int = bool(1)     # True
bool_from_str = bool("")    # False (empty string)
bool_from_str2 = bool("hi") # True (non-empty string)
```

### Getting User Input

```python
# Basic input (always returns string)
name = input("What's your name? ")
print(f"Hello, {name}!")

# Converting input to other types
age_str = input("What's your age? ")
age = int(age_str)

# More concise
age = int(input("What's your age? "))
```

---

## 5. Operators

### Arithmetic Operators

```python
a = 10
b = 3

# Basic arithmetic
addition = a + b        # 13
subtraction = a - b     # 7
multiplication = a * b  # 30
division = a / b        # 3.333...
floor_division = a // b # 3 (integer division)
modulus = a % b         # 1 (remainder)
exponentiation = a ** b # 1000 (10^3)

# Operator precedence (same as math)
result = 2 + 3 * 4      # 14, not 20
result2 = (2 + 3) * 4   # 20
```

### Assignment Operators

```python
x = 10

# Augmented assignment
x += 5   # x = x + 5, result: 15
x -= 3   # x = x - 3, result: 12
x *= 2   # x = x * 2, result: 24
x /= 4   # x = x / 4, result: 6.0
x //= 2  # x = x // 2, result: 3.0
x %= 2   # x = x % 2, result: 1.0
x **= 3  # x = x ** 3, result: 1.0
```

### Comparison Operators

```python
a = 10
b = 5

# Comparison operators (return boolean)
is_equal = a == b        # False
not_equal = a != b       # True
greater = a > b          # True
less = a < b             # False
greater_equal = a >= b   # True
less_equal = a <= b      # False

# Chaining comparisons
age = 25
is_adult = 18 <= age <= 65  # True
```

### Logical Operators

```python
# Boolean values
is_sunny = True
is_warm = False

# Logical operators
nice_weather = is_sunny and is_warm    # False
okay_weather = is_sunny or is_warm     # True
bad_weather = not is_sunny             # False

# With expressions
age = 25
has_license = True
can_drive = age >= 16 and has_license  # True
```

### Identity and Membership Operators

```python
# Identity operators
a = [1, 2, 3]
b = [1, 2, 3]
c = a

print(a is b)    # False (different objects)
print(a is c)    # True (same object)
print(a is not b) # True

# Membership operators
numbers = [1, 2, 3, 4, 5]
print(3 in numbers)     # True
print(6 not in numbers) # True

text = "Hello, World!"
print("Hello" in text)  # True
print("hello" in text)  # False (case sensitive)
```

---

## 6. Control Flow

### Conditional Statements

#### if Statement

```python
age = 18

if age >= 18:
    print("You are an adult")
```

#### if-else Statement

```python
age = 16

if age >= 18:
    print("You are an adult")
else:
    print("You are a minor")
```

#### if-elif-else Statement

```python
score = 85

if score >= 90:
    grade = "A"
elif score >= 80:
    grade = "B"
elif score >= 70:
    grade = "C"
elif score >= 60:
    grade = "D"
else:
    grade = "F"

print(f"Your grade is: {grade}")
```

#### Nested if Statements

```python
weather = "sunny"
temperature = 75

if weather == "sunny":
    if temperature > 70:
        print("Perfect day for a picnic!")
    else:
        print("Sunny but a bit cold")
else:
    print("Not a sunny day")
```

#### Conditional Expressions (Ternary Operator)

```python
age = 20
status = "adult" if age >= 18 else "minor"
print(status)  # Output: adult

# More complex example
x = 10
y = 20
max_value = x if x > y else y
```

### Loops

#### for Loop

```python
# Basic for loop with range
for i in range(5):
    print(i)  # Prints 0, 1, 2, 3, 4

# Range with start and stop
for i in range(1, 6):
    print(i)  # Prints 1, 2, 3, 4, 5

# Range with step
for i in range(0, 10, 2):
    print(i)  # Prints 0, 2, 4, 6, 8

# Iterating over strings
name = "Python"
for letter in name:
    print(letter)

# Iterating over lists
fruits = ["apple", "banana", "orange"]
for fruit in fruits:
    print(fruit)

# Using enumerate to get index and value
for index, fruit in enumerate(fruits):
    print(f"{index}: {fruit}")
```

#### while Loop

```python
# Basic while loop
count = 0
while count < 5:
    print(count)
    count += 1

# While loop with user input
password = ""
while password != "secret":
    password = input("Enter password: ")
print("Access granted!")

# Infinite loop (be careful!)
# while True:
#     print("This runs forever")
#     break  # Use break to exit
```

#### Loop Control Statements

```python
# break statement
for i in range(10):
    if i == 5:
        break  # Exit the loop
    print(i)  # Prints 0, 1, 2, 3, 4

# continue statement
for i in range(10):
    if i % 2 == 0:
        continue  # Skip even numbers
    print(i)  # Prints 1, 3, 5, 7, 9

# else clause with loops
for i in range(5):
    print(i)
else:
    print("Loop completed normally")

# else clause doesn't execute if loop is broken
for i in range(5):
    if i == 3:
        break
    print(i)
else:
    print("This won't print")
```

#### Nested Loops

```python
# Multiplication table
for i in range(1, 4):
    for j in range(1, 4):
        print(f"{i} x {j} = {i * j}")
    print()  # Empty line between tables
```

---

## 7. Functions

### Defining Functions

```python
# Basic function
def greet():
    print("Hello, World!")

# Call the function
greet()
```

### Functions with Parameters

```python
# Function with parameters
def greet_person(name):
    print(f"Hello, {name}!")

greet_person("Alice")

# Function with multiple parameters
def add_numbers(a, b):
    result = a + b
    print(f"{a} + {b} = {result}")

add_numbers(5, 3)
```

### Return Values

```python
# Function that returns a value
def add_numbers(a, b):
    return a + b

result = add_numbers(5, 3)
print(result)  # Output: 8

# Function with multiple return values
def get_name_age():
    name = "John"
    age = 25
    return name, age

person_name, person_age = get_name_age()
print(f"{person_name} is {person_age} years old")
```

### Default Parameters

```python
# Function with default parameter
def greet(name, greeting="Hello"):
    print(f"{greeting}, {name}!")

greet("Alice")              # Uses default greeting
greet("Bob", "Hi")          # Uses custom greeting
greet("Charlie", greeting="Hey")  # Named parameter
```

### Keyword Arguments

```python
def create_profile(name, age, city="Unknown", country="Unknown"):
    print(f"Name: {name}")
    print(f"Age: {age}")
    print(f"City: {city}")
    print(f"Country: {country}")

# Positional arguments
create_profile("Alice", 25)

# Mixed positional and keyword arguments
create_profile("Bob", 30, country="USA")

# All keyword arguments
create_profile(name="Charlie", age=35, city="New York", country="USA")
```

### Variable-Length Arguments

```python
# *args for variable positional arguments
def sum_all(*numbers):
    total = 0
    for num in numbers:
        total += num
    return total

result = sum_all(1, 2, 3, 4, 5)
print(result)  # Output: 15

# **kwargs for variable keyword arguments
def print_info(**info):
    for key, value in info.items():
        print(f"{key}: {value}")

print_info(name="Alice", age=25, city="New York")

# Combining all types
def complex_function(required_arg, *args, **kwargs):
    print(f"Required: {required_arg}")
    print(f"Args: {args}")
    print(f"Kwargs: {kwargs}")

complex_function("test", 1, 2, 3, name="Alice", age=25)
```

### Local and Global Variables

```python
# Global variable
global_var = "I'm global"

def my_function():
    # Local variable
    local_var = "I'm local"
    print(global_var)  # Can access global
    print(local_var)

my_function()
# print(local_var)  # Error: local_var not accessible

# Modifying global variables
counter = 0

def increment():
    global counter
    counter += 1

increment()
print(counter)  # Output: 1
```

### Lambda Functions

```python
# Lambda function (anonymous function)
square = lambda x: x ** 2
print(square(5))  # Output: 25

# Lambda with multiple parameters
add = lambda x, y: x + y
print(add(3, 4))  # Output: 7

# Lambda functions are often used with built-in functions
numbers = [1, 2, 3, 4, 5]
squares = list(map(lambda x: x ** 2, numbers))
print(squares)  # Output: [1, 4, 9, 16, 25]

even_numbers = list(filter(lambda x: x % 2 == 0, numbers))
print(even_numbers)  # Output: [2, 4]
```

---

## 8. Data Structures

### Lists

```python
# Creating lists
fruits = ["apple", "banana", "orange"]
numbers = [1, 2, 3, 4, 5]
mixed_list = ["hello", 42, True, 3.14]
empty_list = []

# Accessing elements (zero-indexed)
print(fruits[0])    # Output: apple
print(fruits[-1])   # Output: orange (last item)
print(fruits[-2])   # Output: banana (second to last)

# Slicing
print(fruits[0:2])  # Output: ['apple', 'banana']
print(fruits[1:])   # Output: ['banana', 'orange']
print(fruits[:2])   # Output: ['apple', 'banana']
print(fruits[:])    # Output: ['apple', 'banana', 'orange'] (copy)

# Modifying lists
fruits[1] = "grape"  # Change element
fruits.append("kiwi")  # Add to end
fruits.insert(0, "strawberry")  # Insert at index
removed = fruits.pop()  # Remove and return last item
fruits.remove("apple")  # Remove first occurrence

# List methods
print(len(fruits))  # Length
print("banana" in fruits)  # Check membership
fruits.sort()  # Sort in place
fruits.reverse()  # Reverse in place
fruits.clear()  # Remove all items

# List comprehensions
squares = [x ** 2 for x in range(10)]
even_squares = [x ** 2 for x in range(10) if x % 2 == 0]
```

### Tuples

```python
# Creating tuples
coordinates = (3, 4)
colors = ("red", "green", "blue")
single_item = (42,)  # Note the comma
empty_tuple = ()

# Accessing elements (same as lists)
print(coordinates[0])  # Output: 3
print(colors[-1])      # Output: blue

# Tuples are immutable
# coordinates[0] = 5  # Error: can't modify

# Tuple unpacking
x, y = coordinates
print(f"x: {x}, y: {y}")

# Multiple assignment using tuples
a, b, c = colors

# Returning multiple values from functions
def get_coordinates():
    return 10, 20

x, y = get_coordinates()
```

### Dictionaries

```python
# Creating dictionaries
person = {
    "name": "Alice",
    "age": 25,
    "city": "New York"
}

# Alternative creation methods
person2 = dict(name="Bob", age=30, city="Boston")
empty_dict = {}

# Accessing values
print(person["name"])      # Output: Alice
print(person.get("age"))   # Output: 25
print(person.get("country", "Unknown"))  # Default value

# Modifying dictionaries
person["age"] = 26         # Update existing key
person["country"] = "USA"  # Add new key
del person["city"]         # Delete key

# Dictionary methods
print(person.keys())       # All keys
print(person.values())     # All values
print(person.items())      # Key-value pairs

# Iterating over dictionaries
for key in person:
    print(f"{key}: {person[key]}")

for key, value in person.items():
    print(f"{key}: {value}")

# Dictionary comprehensions
squares_dict = {x: x ** 2 for x in range(5)}
# Output: {0: 0, 1: 1, 2: 4, 3: 9, 4: 16}
```

### Sets

```python
# Creating sets
fruits = {"apple", "banana", "orange"}
numbers = {1, 2, 3, 4, 5}
empty_set = set()  # Note: {} creates an empty dict

# Adding and removing
fruits.add("kiwi")
fruits.remove("banana")  # Raises error if not found
fruits.discard("grape")  # Doesn't raise error if not found

# Set operations
set1 = {1, 2, 3, 4, 5}
set2 = {4, 5, 6, 7, 8}

union = set1 | set2           # {1, 2, 3, 4, 5, 6, 7, 8}
intersection = set1 & set2    # {4, 5}
difference = set1 - set2      # {1, 2, 3}
symmetric_diff = set1 ^ set2  # {1, 2, 3, 6, 7, 8}

# Set comprehensions
even_squares = {x ** 2 for x in range(10) if x % 2 == 0}
```

### Choosing the Right Data Structure

```python
# Lists: ordered, mutable, allow duplicates
shopping_list = ["milk", "bread", "eggs", "milk"]

# Tuples: ordered, immutable, allow duplicates
rgb_color = (255, 128, 0)

# Dictionaries: unordered, mutable, keys must be unique
student_grades = {"Alice": 85, "Bob": 92, "Charlie": 78}

# Sets: unordered, mutable, no duplicates
unique_visitors = {"alice", "bob", "charlie"}
```

---

## 9. Strings

### String Basics

```python
# Creating strings
single_quotes = 'Hello'
double_quotes = "World"
triple_quotes = """Multi-line
string"""

# Raw strings (ignore escape characters)
path = r"C:\Users\Name\Documents"

# Unicode strings
unicode_text = "Hello 🌍"
```

### String Indexing and Slicing

```python
text = "Python Programming"

# Indexing
print(text[0])   # P
print(text[-1])  # g
print(text[7])   # P

# Slicing
print(text[0:6])    # Python
print(text[7:])     # Programming
print(text[:6])     # Python
print(text[::2])    # Pto rgamn (every 2nd character)
print(text[::-1])   # gnimmargorP nohtyP (reverse)
```

### String Methods

```python
text = "  Hello, World!  "

# Case methods
print(text.upper())        # "  HELLO, WORLD!  "
print(text.lower())        # "  hello, world!  "
print(text.title())        # "  Hello, World!  "
print(text.capitalize())   # "  hello, world!  "
print(text.swapcase())     # "  hELLO, wORLD!  "

# Whitespace methods
print(text.strip())        # "Hello, World!"
print(text.lstrip())       # "Hello, World!  "
print(text.rstrip())       # "  Hello, World!"

# Search and replace
print(text.find("World"))    # 9 (index of first occurrence)
print(text.count("l"))       # 3
print(text.replace("World", "Python"))  # "  Hello, Python!  "

# Split and join
words = "apple,banana,orange".split(",")  # ['apple', 'banana', 'orange']
joined = "-".join(words)                  # "apple-banana-orange"

# Boolean methods
print("123".isdigit())      # True
print("abc".isalpha())      # True
print("abc123".isalnum())   # True
print("Hello World".istitle())  # True
```

### String Formatting

#### Old-style formatting (%)

```python
name = "Alice"
age = 25
print("My name is %s and I'm %d years old" % (name, age))
```

#### str.format() method

```python
name = "Bob"
age = 30
print("My name is {} and I'm {} years old".format(name, age))
print("My name is {0} and I'm {1} years old".format(name, age))
print("My name is {name} and I'm {age} years old".format(name=name, age=age))
```

#### f-strings (formatted string literals) - Python 3.6+

```python
name = "Charlie"
age = 35
print(f"My name is {name} and I'm {age} years old")

# Expressions in f-strings
x = 10
y = 20
print(f"The sum of {x} and {y} is {x + y}")

# Formatting numbers
pi = 3.14159
print(f"Pi is approximately {pi:.2f}")  # Pi is approximately 3.14
```

### String Escape Characters

```python
# Common escape characters
print("Hello\nWorld")     # Newline
print("Hello\tWorld")     # Tab
print("She said \"Hi\"")  # Quote within string
print("C:\\Users\\Name")  # Backslash
print("Line 1\\\nLine 2") # Backslash + newline
```

### Working with Characters

```python
# ASCII values
print(ord('A'))     # 65
print(chr(65))      # A

# Checking character types
char = 'A'
print(char.isalpha())   # True
print(char.isdigit())   # False
print(char.isupper())   # True
print(char.islower())   # False
```

---

## 10. File Input/Output

### Opening and Closing Files

```python
# Basic file opening (manual close)
file = open("example.txt", "r")
content = file.read()
file.close()

# Using with statement (automatic close)
with open("example.txt", "r") as file:
    content = file.read()
    # File automatically closed when exiting the block
```

### File Modes

```python
# Reading modes
with open("file.txt", "r") as file:     # Read text
    pass
with open("file.txt", "rb") as file:    # Read binary
    pass

# Writing modes
with open("file.txt", "w") as file:     # Write (overwrites)
    pass
with open("file.txt", "a") as file:     # Append
    pass
with open("file.txt", "x") as file:     # Exclusive creation (fails if exists)
    pass

# Read and write
with open("file.txt", "r+") as file:    # Read and write
    pass
```

### Reading Files

```python
# Read entire file
with open("example.txt", "r") as file:
    content = file.read()
    print(content)

# Read line by line
with open("example.txt", "r") as file:
    for line in file:
        print(line.strip())  # strip() removes newline

# Read specific number of characters
with open("example.txt", "r") as file:
    first_10_chars = file.read(10)

# Read all lines into a list
with open("example.txt", "r") as file:
    lines = file.readlines()

# Read one line at a time
with open("example.txt", "r") as file:
    first_line = file.readline()
    second_line = file.readline()
```

### Writing Files

```python
# Write string to file
with open("output.txt", "w") as file:
    file.write("Hello, World!\n")
    file.write("This is the second line.\n")

# Write multiple lines
lines = ["Line 1\n", "Line 2\n", "Line 3\n"]
with open("output.txt", "w") as file:
    file.writelines(lines)

# Append to file
with open("output.txt", "a") as file:
    file.write("This line is appended.\n")
```

### Working with CSV Files

```python
import csv

# Writing CSV
data = [
    ["Name", "Age", "City"],
    ["Alice", 25, "New York"],
    ["Bob", 30, "Boston"],
    ["Charlie", 35, "Chicago"]
]

with open("people.csv", "w", newline="") as file:
    writer = csv.writer(file)
    writer.writerows(data)

# Reading CSV
with open("people.csv", "r") as file:
    reader = csv.reader(file)
    for row in reader:
        print(row)

# Using DictReader for column names
with open("people.csv", "r") as file:
    reader = csv.DictReader(file)
    for row in reader:
        print(f"{row['Name']} is {row['Age']} years old")
```

### File and Directory Operations

```python
import os

# Check if file exists
if os.path.exists("example.txt"):
    print("File exists")

# Get file information
import os.path
print(os.path.getsize("example.txt"))  # File size in bytes
print(os.path.getctime("example.txt"))  # Creation time

# Directory operations
os.mkdir("new_directory")              # Create directory
os.makedirs("path/to/directory")       # Create nested directories
os.listdir(".")                        # List directory contents
os.getcwd()                            # Get current directory
os.chdir("new_directory")              # Change directory

# Remove files and directories
os.remove("file.txt")                  # Remove file
os.rmdir("directory")                  # Remove empty directory
import shutil
shutil.rmtree("directory")             # Remove directory and contents
```

### Error Handling with Files

```python
try:
    with open("nonexistent.txt", "r") as file:
        content = file.read()
except FileNotFoundError:
    print("File not found!")
except PermissionError:
    print("Permission denied!")
except Exception as e:
    print(f"An error occurred: {e}")
```

---

## 11. Error Handling

### Understanding Exceptions

```python
# Common exceptions
print(10 / 0)           # ZeroDivisionError
print(int("hello"))     # ValueError
print(my_list[10])      # IndexError (if list has < 11 items)
print(my_dict["key"])   # KeyError (if key doesn't exist)
```

### try-except Blocks

```python
# Basic exception handling
try:
    number = int(input("Enter a number: "))
    result = 10 / number
    print(f"Result: {result}")
except ValueError:
    print("That's not a valid number!")
except ZeroDivisionError:
    print("Cannot divide by zero!")
```

### Multiple Exception Types

```python
# Handle multiple exceptions
try:
    file = open("data.txt", "r")
    data = file.read()
    number = int(data)
    result = 100 / number
except (FileNotFoundError, IOError):
    print("File error occurred")
except ValueError:
    print("Invalid number in file")
except ZeroDivisionError:
    print("Cannot divide by zero")
except Exception as e:
    print(f"Unexpected error: {e}")
```

### else and finally Clauses

```python
try:
    number = int(input("Enter a number: "))
    result = 10 / number
except ValueError:
    print("Invalid number!")
except ZeroDivisionError:
    print("Cannot divide by zero!")
else:
    # Runs if no exception occurred
    print(f"Result: {result}")
finally:
    # Always runs
    print("Operation completed")
```

### Raising Exceptions

```python
def validate_age(age):
    if age < 0:
        raise ValueError("Age cannot be negative")
    if age > 150:
        raise ValueError("Age seems unrealistic")
    return True

try:
    age = int(input("Enter your age: "))
    validate_age(age)
    print(f"Age {age} is valid")
except ValueError as e:
    print(f"Error: {e}")
```

### Custom Exceptions

```python
class CustomError(Exception):
    """Custom exception class"""
    pass

class ValidationError(Exception):
    """Raised when validation fails"""
    def __init__(self, message, code=None):
        super().__init__(message)
        self.code = code

def process_data(data):
    if not data:
        raise ValidationError("Data cannot be empty", code="EMPTY_DATA")
    if len(data) < 5:
        raise ValidationError("Data too short", code="SHORT_DATA")

try:
    process_data("hi")
except ValidationError as e:
    print(f"Validation error: {e}")
    print(f"Error code: {e.code}")
```

### Best Practices

```python
# Be specific with exceptions
try:
    value = int(user_input)
except ValueError:  # Specific exception
    print("Please enter a valid number")

# Don't catch all exceptions unless necessary
try:
    risky_operation()
except Exception:  # Too broad, avoid unless necessary
    pass

# Use finally for cleanup
file = None
try:
    file = open("data.txt", "r")
    # Process file
except IOError:
    print("File error")
finally:
    if file:
        file.close()

# Better: use context managers
try:
    with open("data.txt", "r") as file:
        # Process file
        pass
except IOError:
    print("File error")
```

---

## 12. Modules and Packages

### What are Modules?

A module is a file containing Python code. It can define functions, classes, and variables, and can also include runnable code.

### Importing Modules

```python
# Import entire module
import math
print(math.pi)
print(math.sqrt(16))

# Import specific functions
from math import pi, sqrt
print(pi)
print(sqrt(16))

# Import with alias
import math as m
print(m.pi)

from math import sqrt as square_root
print(square_root(16))

# Import all (not recommended)
from math import *
print(pi)  # Available without math. prefix
```

### Creating Your Own Modules

Create a file named `my_module.py`:
```python
# my_module.py
def greet(name):
    return f"Hello, {name}!"

def add(a, b):
    return a + b

PI = 3.14159

class Calculator:
    def multiply(self, a, b):
        return a * b
```

Using your module:
```python
# main.py
import my_module

print(my_module.greet("Alice"))
print(my_module.add(5, 3))
print(my_module.PI)

calc = my_module.Calculator()
print(calc.multiply(4, 5))
```

### Module Search Path

```python
import sys

# See where Python looks for modules
for path in sys.path:
    print(path)

# Add new path
sys.path.append("/path/to/my/modules")
```

### Standard Library Modules

#### os module
```python
import os

print(os.getcwd())          # Current directory
print(os.listdir("."))      # List directory contents
os.mkdir("new_folder")      # Create directory
print(os.environ["PATH"])   # Environment variables
```

#### datetime module
```python
from datetime import datetime, date, time

now = datetime.now()
print(now)

today = date.today()
print(today)

# Formatting dates
formatted = now.strftime("%Y-%m-%d %H:%M:%S")
print(formatted)

# Parsing dates
date_string = "2024-01-15"
parsed_date = datetime.strptime(date_string, "%Y-%m-%d")
print(parsed_date)
```

#### random module
```python
import random

# Random number between 0 and 1
print(random.random())

# Random integer in range
print(random.randint(1, 10))

# Random choice from list
colors = ["red", "green", "blue"]
print(random.choice(colors))

# Shuffle list
numbers = [1, 2, 3, 4, 5]
random.shuffle(numbers)
print(numbers)
```

#### json module
```python
import json

# Python dict to JSON string
data = {"name": "Alice", "age": 25, "city": "New York"}
json_string = json.dumps(data)
print(json_string)

# JSON string to Python dict
parsed_data = json.loads(json_string)
print(parsed_data)

# Write to file
with open("data.json", "w") as file:
    json.dump(data, file)

# Read from file
with open("data.json", "r") as file:
    loaded_data = json.load(file)
    print(loaded_data)
```

### Packages

A package is a directory containing multiple modules. It must contain an `__init__.py` file.

Create package structure:
```
my_package/
    __init__.py
    math_operations.py
    string_operations.py
```

`my_package/__init__.py`:
```python
# Can be empty or contain initialization code
from .math_operations import add, subtract
from .string_operations import reverse_string

__version__ = "1.0.0"
```

`my_package/math_operations.py`:
```python
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b
```

`my_package/string_operations.py`:
```python
def reverse_string(s):
    return s[::-1]

def count_vowels(s):
    return sum(1 for char in s.lower() if char in 'aeiou')
```

Using the package:
```python
import my_package

print(my_package.add(5, 3))
print(my_package.reverse_string("hello"))

# Or import specific modules
from my_package import math_operations
print(math_operations.subtract(10, 3))
```

### Installing Third-Party Packages

```bash
# Install package using pip
pip install requests

# Install specific version
pip install requests==2.25.1

# Install from requirements file
pip install -r requirements.txt

# List installed packages
pip list

# Show package information
pip show requests
```

Using installed packages:
```python
import requests

response = requests.get("https://api.github.com/users/octocat")
print(response.json())
```

---

## 13. Object-Oriented Programming Basics

### Classes and Objects

```python
# Define a class
class Dog:
    # Class variable (shared by all instances)
    species = "Canis familiaris"
    
    # Constructor method
    def __init__(self, name, age):
        # Instance variables
        self.name = name
        self.age = age
    
    # Instance method
    def bark(self):
        return f"{self.name} says Woof!"
    
    def get_info(self):
        return f"{self.name} is {self.age} years old"

# Create objects (instances)
dog1 = Dog("Buddy", 3)
dog2 = Dog("Luna", 5)

# Access attributes and methods
print(dog1.name)        # Buddy
print(dog1.bark())      # Buddy says Woof!
print(dog2.get_info())  # Luna is 5 years old
print(Dog.species)      # Canis familiaris
```

### Instance vs Class Variables

```python
class Counter:
    # Class variable
    total_instances = 0
    
    def __init__(self, initial_value=0):
        # Instance variable
        self.value = initial_value
        # Modify class variable
        Counter.total_instances += 1
    
    def increment(self):
        self.value += 1

counter1 = Counter(10)
counter2 = Counter(20)

print(counter1.value)           # 10
print(counter2.value)           # 20
print(Counter.total_instances)  # 2
```

### Methods Types

```python
class MathUtils:
    pi = 3.14159
    
    def __init__(self, value):
        self.value = value
    
    # Instance method
    def double(self):
        return self.value * 2
    
    # Class method
    @classmethod
    def get_pi(cls):
        return cls.pi
    
    # Static method
    @staticmethod
    def add(a, b):
        return a + b

# Usage
math_obj = MathUtils(5)
print(math_obj.double())        # 10 (instance method)
print(MathUtils.get_pi())       # 3.14159 (class method)
print(MathUtils.add(3, 4))      # 7 (static method)
```

### Inheritance

```python
# Parent class
class Animal:
    def __init__(self, name, species):
        self.name = name
        self.species = species
    
    def make_sound(self):
        return "Some generic animal sound"
    
    def get_info(self):
        return f"{self.name} is a {self.species}"

# Child class
class Dog(Animal):
    def __init__(self, name, breed):
        # Call parent constructor
        super().__init__(name, "Dog")
        self.breed = breed
    
    # Override parent method
    def make_sound(self):
        return "Woof!"
    
    # Add new method
    def fetch(self):
        return f"{self.name} is fetching the ball!"

class Cat(Animal):
    def __init__(self, name, indoor=True):
        super().__init__(name, "Cat")
        self.indoor = indoor
    
    def make_sound(self):
        return "Meow!"

# Usage
dog = Dog("Buddy", "Golden Retriever")
cat = Cat("Whiskers")

print(dog.get_info())       # Buddy is a Dog
print(dog.make_sound())     # Woof!
print(dog.fetch())          # Buddy is fetching the ball!

print(cat.get_info())       # Whiskers is a Cat
print(cat.make_sound())     # Meow!
```

### Encapsulation (Private Attributes)

```python
class BankAccount:
    def __init__(self, account_number, initial_balance=0):
        self.account_number = account_number
        self._balance = initial_balance  # Protected (convention)
        self.__pin = "1234"              # Private (name mangling)
    
    def deposit(self, amount):
        if amount > 0:
            self._balance += amount
            return True
        return False
    
    def withdraw(self, amount):
        if 0 < amount <= self._balance:
            self._balance -= amount
            return True
        return False
    
    def get_balance(self):
        return self._balance
    
    def _internal_method(self):  # Protected method
        return "Internal processing"
    
    def __private_method(self):  # Private method
        return "Private processing"

account = BankAccount("12345", 1000)
print(account.get_balance())    # 1000
account.deposit(500)
print(account.get_balance())    # 1500

# These work but are discouraged
print(account._balance)         # 1500 (protected)
print(account._internal_method())  # Internal processing

# This is name-mangled
# print(account.__private_method())  # Error
print(account._BankAccount__private_method())  # Works but ugly
```

### Special Methods (Magic Methods)

```python
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __str__(self):
        # For print() and str()
        return f"Point({self.x}, {self.y})"
    
    def __repr__(self):
        # For repr() and interactive display
        return f"Point(x={self.x}, y={self.y})"
    
    def __add__(self, other):
        # For + operator
        return Point(self.x + other.x, self.y + other.y)
    
    def __eq__(self, other):
        # For == operator
        return self.x == other.x and self.y == other.y
    
    def __len__(self):
        # For len()
        return 2  # Always 2 coordinates
    
    def __getitem__(self, index):
        # For indexing p[0], p[1]
        if index == 0:
            return self.x
        elif index == 1:
            return self.y
        else:
            raise IndexError("Point index out of range")

# Usage
p1 = Point(3, 4)
p2 = Point(1, 2)

print(p1)           # Point(3, 4)
print(repr(p1))     # Point(x=3, y=4)
print(p1 + p2)      # Point(4, 6)
print(p1 == p2)     # False
print(len(p1))      # 2
print(p1[0])        # 3
print(p1[1])        # 4
```

### Property Decorators

```python
class Circle:
    def __init__(self, radius):
        self._radius = radius
    
    @property
    def radius(self):
        return self._radius
    
    @radius.setter
    def radius(self, value):
        if value < 0:
            raise ValueError("Radius cannot be negative")
        self._radius = value
    
    @property
    def area(self):
        return 3.14159 * self._radius ** 2
    
    @property
    def diameter(self):
        return 2 * self._radius

circle = Circle(5)
print(circle.radius)    # 5
print(circle.area)      # 78.53975
print(circle.diameter)  # 10

circle.radius = 7       # Uses setter
print(circle.area)      # 153.93791

# circle.radius = -1    # Raises ValueError
```

---

## 14. Standard Library Overview

### Collections Module

```python
from collections import defaultdict, Counter, namedtuple, deque

# defaultdict - never raises KeyError
dd = defaultdict(list)
dd['key1'].append('value1')  # No KeyError
print(dd)  # defaultdict(<class 'list'>, {'key1': ['value1']})

# Counter - count occurrences
text = "hello world"
counter = Counter(text)
print(counter)  # Counter({'l': 3, 'o': 2, 'h': 1, 'e': 1, ' ': 1, 'w': 1, 'r': 1, 'd': 1})
print(counter.most_common(3))  # [('l', 3), ('o', 2), ('h', 1)]

# namedtuple - tuple with named fields
Point = namedtuple('Point', ['x', 'y'])
p = Point(3, 4)
print(p.x, p.y)  # 3 4

# deque - double-ended queue
dq = deque([1, 2, 3])
dq.appendleft(0)    # Add to left
dq.append(4)        # Add to right
print(dq)           # deque([0, 1, 2, 3, 4])
```

### itertools Module

```python
import itertools

# count - infinite counting
# for i in itertools.count(10, 2):  # Start at 10, step by 2
#     if i > 20:
#         break
#     print(i)  # 10, 12, 14, 16, 18, 20

# cycle - infinite cycling
colors = ['red', 'green', 'blue']
color_cycle = itertools.cycle(colors)
for _ in range(6):
    print(next(color_cycle))  # red, green, blue, red, green, blue

# combinations and permutations
letters = ['A', 'B', 'C']
print(list(itertools.combinations(letters, 2)))  # [('A', 'B'), ('A', 'C'), ('B', 'C')]
print(list(itertools.permutations(letters, 2)))  # [('A', 'B'), ('A', 'C'), ('B', 'A'), ('B', 'C'), ('C', 'A'), ('C', 'B')]

# chain - flatten iterables
list1 = [1, 2, 3]
list2 = [4, 5, 6]
chained = list(itertools.chain(list1, list2))
print(chained)  # [1, 2, 3, 4, 5, 6]
```

### functools Module

```python
from functools import partial, reduce

# partial - partially apply functions
def multiply(x, y):
    return x * y

double = partial(multiply, 2)  # x is always 2
print(double(5))  # 10

# reduce - apply function cumulatively
numbers = [1, 2, 3, 4, 5]
sum_all = reduce(lambda x, y: x + y, numbers)
print(sum_all)  # 15
```

### pathlib Module (Modern file handling)

```python
from pathlib import Path

# Create path objects
current_dir = Path('.')
home_dir = Path.home()
file_path = Path('data') / 'file.txt'  # OS-independent path joining

# Path operations
print(file_path.exists())       # Check if exists
print(file_path.is_file())      # Check if file
print(file_path.is_dir())       # Check if directory
print(file_path.suffix)         # File extension
print(file_path.stem)           # Filename without extension
print(file_path.parent)         # Parent directory

# Create directory
new_dir = Path('new_directory')
new_dir.mkdir(exist_ok=True)

# List files
for file in Path('.').glob('*.py'):
    print(file)
```

### urllib Module (Web requests)

```python
from urllib.request import urlopen
from urllib.parse import urlencode, urlparse

# Simple GET request
response = urlopen('https://httpbin.org/get')
data = response.read().decode('utf-8')
print(data)

# Parse URLs
url = 'https://example.com/path?param=value'
parsed = urlparse(url)
print(parsed.scheme)  # https
print(parsed.netloc)  # example.com
print(parsed.path)    # /path
```

### argparse Module (Command-line arguments)

```python
import argparse

parser = argparse.ArgumentParser(description='A simple calculator')
parser.add_argument('operation', choices=['add', 'subtract', 'multiply', 'divide'])
parser.add_argument('x', type=float, help='First number')
parser.add_argument('y', type=float, help='Second number')
parser.add_argument('--verbose', '-v', action='store_true', help='Verbose output')

args = parser.parse_args()

if args.operation == 'add':
    result = args.x + args.y
elif args.operation == 'subtract':
    result = args.x - args.y
elif args.operation == 'multiply':
    result = args.x * args.y
elif args.operation == 'divide':
    result = args.x / args.y

if args.verbose:
    print(f"Performing {args.operation} on {args.x} and {args.y}")
print(f"Result: {result}")

# Usage: python script.py add 5 3 --verbose
```

---

## 15. Best Practices

### Code Style (PEP 8)

```python
# Good naming conventions
user_name = "john_doe"          # Snake case for variables and functions
USER_CONSTANT = "ADMIN"         # Uppercase for constants
class UserAccount:              # PascalCase for classes
    pass

# Function and variable naming
def calculate_total_price(items, tax_rate):
    """Calculate total price including tax."""
    subtotal = sum(item.price for item in items)
    tax_amount = subtotal * tax_rate
    return subtotal + tax_amount

# Imports organization
import os
import sys
from collections import defaultdict
import requests

# Line length (79 characters max)
long_string = ("This is a very long string that should be split "
               "across multiple lines for better readability")

# Whitespace
def function_name(param1, param2):  # Spaces around operators
    if param1 == param2:            # Spaces around operators
        return param1 + param2      # Spaces around operators
    
# List comprehensions
squares = [x**2 for x in range(10)]
even_squares = [x**2 for x in range(10) if x % 2 == 0]
```

### Documentation

```python
def calculate_distance(point1, point2):
    """
    Calculate the Euclidean distance between two points.
    
    Args:
        point1 (tuple): First point coordinates (x, y)
        point2 (tuple): Second point coordinates (x, y)
    
    Returns:
        float: The distance between the points
    
    Example:
        >>> calculate_distance((0, 0), (3, 4))
        5.0
    """
    x1, y1 = point1
    x2, y2 = point2
    return ((x2 - x1)**2 + (y2 - y1)**2)**0.5

class BankAccount:
    """
    A simple bank account class.
    
    Attributes:
        account_number (str): The account number
        balance (float): Current account balance
    """
    
    def __init__(self, account_number, initial_balance=0):
        """
        Initialize a new bank account.
        
        Args:
            account_number (str): Unique account identifier
            initial_balance (float, optional): Starting balance. Defaults to 0.
        """
        self.account_number = account_number
        self.balance = initial_balance
```

### Error Handling Best Practices

```python
# Be specific with exceptions
def divide_numbers(a, b):
    """Divide two numbers with proper error handling."""
    try:
        result = a / b
    except ZeroDivisionError:
        print("Error: Cannot divide by zero")
        return None
    except TypeError:
        print("Error: Both arguments must be numbers")
        return None
    else:
        return result

# Use finally for cleanup
def read_file_safely(filename):
    """Read file with proper resource management."""
    file = None
    try:
        file = open(filename, 'r')
        return file.read()
    except FileNotFoundError:
        print(f"Error: File {filename} not found")
        return None
    finally:
        if file:
            file.close()

# Better: use context managers
def read_file_better(filename):
    """Read file using context manager."""
    try:
        with open(filename, 'r') as file:
            return file.read()
    except FileNotFoundError:
        print(f"Error: File {filename} not found")
        return None
```

### Writing Readable Code

```python
# Use meaningful variable names
# Bad
d = {"n": "John", "a": 25}
t = d["a"] * 365 * 24 * 60 * 60

# Good
person = {"name": "John", "age": 25}
seconds_lived = person["age"] * 365 * 24 * 60 * 60

# Extract magic numbers into constants
# Bad
if age >= 18:
    can_vote = True

# Good
VOTING_AGE = 18
if age >= VOTING_AGE:
    can_vote = True

# Use functions to break down complex logic
# Bad
total = 0
for item in cart:
    if item.category == "electronics":
        total += item.price * 0.9  # 10% discount
    else:
        total += item.price

# Good
def calculate_item_price(item):
    """Calculate price with applicable discounts."""
    ELECTRONICS_DISCOUNT = 0.1
    if item.category == "electronics":
        return item.price * (1 - ELECTRONICS_DISCOUNT)
    return item.price

total = sum(calculate_item_price(item) for item in cart)
```

### Performance Tips

```python
# Use list comprehensions instead of loops when appropriate
# Slower
squares = []
for x in range(1000):
    squares.append(x**2)

# Faster
squares = [x**2 for x in range(1000)]

# Use built-in functions
# Slower
def find_max(numbers):
    max_num = numbers[0]
    for num in numbers:
        if num > max_num:
            max_num = num
    return max_num

# Faster
max_num = max(numbers)

# Use sets for membership testing
# Slower with lists
large_list = list(range(10000))
if 5000 in large_list:  # O(n) operation
    print("Found")

# Faster with sets
large_set = set(range(10000))
if 5000 in large_set:   # O(1) operation
    print("Found")
```

### Testing Your Code

```python
# Simple testing with assert
def add_numbers(a, b):
    """Add two numbers."""
    return a + b

# Test the function
assert add_numbers(2, 3) == 5
assert add_numbers(-1, 1) == 0
assert add_numbers(0, 0) == 0

# Using doctest
def factorial(n):
    """
    Calculate factorial of n.
    
    >>> factorial(5)
    120
    >>> factorial(0)
    1
    >>> factorial(1)
    1
    """
    if n <= 1:
        return 1
    return n * factorial(n - 1)

if __name__ == "__main__":
    import doctest
    doctest.testmod()
```

---

## 16. Exercises and Projects

### Beginner Exercises

#### Exercise 1: Basic Calculator
```python
def calculator():
    """Simple calculator with basic operations."""
    print("Simple Calculator")
    print("Operations: +, -, *, /")
    
    try:
        num1 = float(input("Enter first number: "))
        operation = input("Enter operation (+, -, *, /): ")
        num2 = float(input("Enter second number: "))
        
        if operation == '+':
            result = num1 + num2
        elif operation == '-':
            result = num1 - num2
        elif operation == '*':
            result = num1 * num2
        elif operation == '/':
            if num2 != 0:
                result = num1 / num2
            else:
                print("Error: Division by zero!")
                return
        else:
            print("Error: Invalid operation!")
            return
        
        print(f"{num1} {operation} {num2} = {result}")
    
    except ValueError:
        print("Error: Please enter valid numbers!")

# calculator()
```

#### Exercise 2: Number Guessing Game
```python
import random

def guessing_game():
    """Number guessing game."""
    number = random.randint(1, 100)
    attempts = 0
    max_attempts = 7
    
    print("Welcome to the Number Guessing Game!")
    print("I'm thinking of a number between 1 and 100.")
    print(f"You have {max_attempts} attempts to guess it.")
    
    while attempts < max_attempts:
        try:
            guess = int(input("Enter your guess: "))
            attempts += 1
            
            if guess == number:
                print(f"Congratulations! You guessed it in {attempts} attempts!")
                return
            elif guess < number:
                print("Too low!")
            else:
                print("Too high!")
            
            remaining = max_attempts - attempts
            if remaining > 0:
                print(f"You have {remaining} attempts left.")
        
        except ValueError:
            print("Please enter a valid number!")
    
    print(f"Game over! The number was {number}.")

# guessing_game()
```

#### Exercise 3: Word Counter
```python
def word_counter():
    """Count words, characters, and lines in text."""
    print("Text Analysis Tool")
    print("Enter your text (press Enter twice to finish):")
    
    lines = []
    while True:
        line = input()
        if line == "":
            break
        lines.append(line)
    
    text = "\n".join(lines)
    
    # Count statistics
    char_count = len(text)
    char_count_no_spaces = len(text.replace(" ", ""))
    word_count = len(text.split())
    line_count = len(lines)
    
    print("\n--- Analysis Results ---")
    print(f"Characters (with spaces): {char_count}")
    print(f"Characters (without spaces): {char_count_no_spaces}")
    print(f"Words: {word_count}")
    print(f"Lines: {line_count}")

# word_counter()
```

### Intermediate Projects

#### Project 1: To-Do List Manager
```python
import json
from datetime import datetime

class TodoManager:
    """A simple to-do list manager."""
    
    def __init__(self, filename="todos.json"):
        self.filename = filename
        self.todos = self.load_todos()
    
    def load_todos(self):
        """Load todos from file."""
        try:
            with open(self.filename, 'r') as file:
                return json.load(file)
        except FileNotFoundError:
            return []
    
    def save_todos(self):
        """Save todos to file."""
        with open(self.filename, 'w') as file:
            json.dump(self.todos, file, indent=2)
    
    def add_todo(self, task, priority="medium"):
        """Add a new todo."""
        todo = {
            "id": len(self.todos) + 1,
            "task": task,
            "priority": priority,
            "completed": False,
            "created": datetime.now().isoformat()
        }
        self.todos.append(todo)
        self.save_todos()
        print(f"Added: {task}")
    
    def list_todos(self):
        """List all todos."""
        if not self.todos:
            print("No todos found.")
            return
        
        print("\n--- Your To-Do List ---")
        for todo in self.todos:
            status = "✓" if todo["completed"] else "○"
            priority = todo["priority"].upper()
            print(f"{status} [{priority}] {todo['task']} (ID: {todo['id']})")
    
    def complete_todo(self, todo_id):
        """Mark a todo as completed."""
        for todo in self.todos:
            if todo["id"] == todo_id:
                todo["completed"] = True
                self.save_todos()
                print(f"Completed: {todo['task']}")
                return
        print("Todo not found.")
    
    def delete_todo(self, todo_id):
        """Delete a todo."""
        for i, todo in enumerate(self.todos):
            if todo["id"] == todo_id:
                deleted_task = self.todos.pop(i)["task"]
                self.save_todos()
                print(f"Deleted: {deleted_task}")
                return
        print("Todo not found.")
    
    def run(self):
        """Main application loop."""
        while True:
            print("\n--- To-Do Manager ---")
            print("1. Add todo")
            print("2. List todos")
            print("3. Complete todo")
            print("4. Delete todo")
            print("5. Exit")
            
            choice = input("Choose an option: ")
            
            if choice == "1":
                task = input("Enter task: ")
                priority = input("Priority (low/medium/high): ") or "medium"
                self.add_todo(task, priority)
            elif choice == "2":
                self.list_todos()
            elif choice == "3":
                try:
                    todo_id = int(input("Enter todo ID to complete: "))
                    self.complete_todo(todo_id)
                except ValueError:
                    print("Invalid ID.")
            elif choice == "4":
                try:
                    todo_id = int(input("Enter todo ID to delete: "))
                    self.delete_todo(todo_id)
                except ValueError:
                    print("Invalid ID.")
            elif choice == "5":
                print("Goodbye!")
                break
            else:
                print("Invalid choice.")

# To run the application:
# todo_manager = TodoManager()
# todo_manager.run()
```

#### Project 2: Contact Book
```python
import json
import re

class ContactBook:
    """A simple contact book application."""
    
    def __init__(self, filename="contacts.json"):
        self.filename = filename
        self.contacts = self.load_contacts()
    
    def load_contacts(self):
        """Load contacts from file."""
        try:
            with open(self.filename, 'r') as file:
                return json.load(file)
        except FileNotFoundError:
            return {}
    
    def save_contacts(self):
        """Save contacts to file."""
        with open(self.filename, 'w') as file:
            json.dump(self.contacts, file, indent=2)
    
    def validate_email(self, email):
        """Validate email format."""
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return re.match(pattern, email) is not None
    
    def validate_phone(self, phone):
        """Validate phone number format."""
        # Simple validation - only digits, spaces, hyphens, parentheses
        pattern = r'^[\d\s\-\(\)]+$'
        return re.match(pattern, phone) is not None and len(phone.replace(' ', '').replace('-', '').replace('(', '').replace(')', '')) >= 10
    
    def add_contact(self):
        """Add a new contact."""
        print("\n--- Add New Contact ---")
        name = input("Name: ").strip()
        
        if not name:
            print("Name cannot be empty.")
            return
        
        if name.lower() in self.contacts:
            print("Contact already exists.")
            return
        
        phone = input("Phone: ").strip()
        if phone and not self.validate_phone(phone):
            print("Invalid phone number format.")
            return
        
        email = input("Email: ").strip()
        if email and not self.validate_email(email):
            print("Invalid email format.")
            return
        
        address = input("Address: ").strip()
        
        self.contacts[name.lower()] = {
            "name": name,
            "phone": phone,
            "email": email,
            "address": address
        }
        
        self.save_contacts()
        print(f"Contact '{name}' added successfully!")
    
    def search_contacts(self):
        """Search for contacts."""
        if not self.contacts:
            print("No contacts found.")
            return
        
        query = input("Enter name to search: ").strip().lower()
        
        matches = []
        for key, contact in self.contacts.items():
            if query in key:
                matches.append(contact)
        
        if matches:
            print(f"\n--- Search Results ({len(matches)} found) ---")
            for contact in matches:
                self.display_contact(contact)
        else:
            print("No contacts found.")
    
    def display_contact(self, contact):
        """Display a single contact."""
        print(f"Name: {contact['name']}")
        if contact['phone']:
            print(f"Phone: {contact['phone']}")
        if contact['email']:
            print(f"Email: {contact['email']}")
        if contact['address']:
            print(f"Address: {contact['address']}")
        print("-" * 30)
    
    def list_all_contacts(self):
        """List all contacts."""
        if not self.contacts:
            print("No contacts found.")
            return
        
        print(f"\n--- All Contacts ({len(self.contacts)}) ---")
        for contact in sorted(self.contacts.values(), key=lambda x: x['name']):
            self.display_contact(contact)
    
    def delete_contact(self):
        """Delete a contact."""
        if not self.contacts:
            print("No contacts found.")
            return
        
        name = input("Enter name to delete: ").strip().lower()
        
        if name in self.contacts:
            deleted_name = self.contacts[name]['name']
            del self.contacts[name]
            self.save_contacts()
            print(f"Contact '{deleted_name}' deleted successfully!")
        else:
            print("Contact not found.")
    
    def run(self):
        """Main application loop."""
        while True:
            print("\n--- Contact Book ---")
            print("1. Add contact")
            print("2. Search contacts")
            print("3. List all contacts")
            print("4. Delete contact")
            print("5. Exit")
            
            choice = input("Choose an option: ")
            
            if choice == "1":
                self.add_contact()
            elif choice == "2":
                self.search_contacts()
            elif choice == "3":
                self.list_all_contacts()
            elif choice == "4":
                self.delete_contact()
            elif choice == "5":
                print("Goodbye!")
                break
            else:
                print("Invalid choice.")

# To run the application:
# contact_book = ContactBook()
# contact_book.run()
```

### Practice Challenges

1. **Palindrome Checker**: Write a function that checks if a string is a palindrome (reads the same forwards and backwards).

2. **FizzBuzz**: Print numbers 1-100, but replace multiples of 3 with "Fizz", multiples of 5 with "Buzz", and multiples of both with "FizzBuzz".

3. **File Organizer**: Create a script that organizes files in a directory by their extensions.

4. **Password Generator**: Build a secure password generator with customizable length and character sets.

5. **Simple Web Scraper**: Use the `requests` library to fetch data from a website and extract specific information.

6. **CSV Data Analyzer**: Read a CSV file and provide statistics (mean, median, mode) for numerical columns.

7. **Log File Parser**: Parse log files and extract useful information (error counts, most common errors, etc.).

8. **Unit Converter**: Create a program that converts between different units (temperature, length, weight, etc.).

---

## Conclusion

Congratulations! You've completed the Python Getting Started course. You now have a solid foundation in Python programming, including:

- **Basic syntax and concepts**
- **Data types and structures**
- **Control flow and functions**
- **Object-oriented programming basics**
- **File handling and error management**
- **Standard library usage**
- **Best practices and code style**

### Next Steps

1. **Practice regularly** with coding exercises
2. **Build projects** to apply your knowledge
3. **Explore advanced Python topics** (decorators, generators, context managers)
4. **Learn popular libraries** (requests, pandas, numpy, matplotlib)
5. **Choose a specialization** (web development, data science, automation)

### Additional Resources

- **Official Python Documentation**: https://docs.python.org/
- **Python.org Tutorial**: https://docs.python.org/tutorial/
- **PEP 8 Style Guide**: https://pep8.org/
- **Python Package Index (PyPI)**: https://pypi.org/
- **Real Python**: https://realpython.com/

Remember: The best way to learn programming is by writing code. Start with small projects and gradually work your way up to more complex applications. Good luck with your Python journey!