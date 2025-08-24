# Advanced Python Programming Course - Part 1

## Table of Contents

1. [Advanced Functions and Decorators](#1-advanced-functions-and-decorators)
2. [Generators and Iterators](#2-generators-and-iterators)
3. [Context Managers](#3-context-managers)
4. [Advanced Object-Oriented Programming](#4-advanced-object-oriented-programming)
5. [Metaclasses](#5-metaclasses)
6. [Descriptors](#6-descriptors)
7. [Advanced Data Structures](#7-advanced-data-structures)
8. [Memory Management and Performance](#8-memory-management-and-performance)
9. [Concurrency - Threading](#9-concurrency-threading)
10. [Concurrency - Multiprocessing](#10-concurrency-multiprocessing)

---

## 1. Advanced Functions and Decorators

### Function Introspection

```python
import inspect
from functools import wraps

def example_function(a: int, b: str = "default", *args, **kwargs) -> str:
    """Example function for introspection."""
    return f"{a}, {b}, {args}, {kwargs}"

# Inspect function properties
print(inspect.signature(example_function))
print(inspect.getfullargspec(example_function))
print(inspect.getsource(example_function))

# Get annotations
print(example_function.__annotations__)
print(example_function.__doc__)
print(example_function.__name__)
```

### Advanced Decorators

#### Basic Decorator with Arguments

```python
def repeat(times):
    """Decorator that repeats function execution."""
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            results = []
            for _ in range(times):
                result = func(*args, **kwargs)
                results.append(result)
            return results
        return wrapper
    return decorator

@repeat(3)
def greet(name):
    return f"Hello, {name}!"

print(greet("Alice"))  # ['Hello, Alice!', 'Hello, Alice!', 'Hello, Alice!']
```

#### Class-Based Decorators

```python
class CountCalls:
    """Decorator class to count function calls."""
    
    def __init__(self, func):
        self.func = func
        self.count = 0
        wraps(func)(self)
    
    def __call__(self, *args, **kwargs):
        self.count += 1
        print(f"Call {self.count} of {self.func.__name__}")
        return self.func(*args, **kwargs)
    
    def __get__(self, obj, objtype):
        """Support instance methods."""
        if obj is None:
            return self
        return functools.partial(self.__call__, obj)

@CountCalls
def say_hello():
    return "Hello!"

say_hello()  # Call 1 of say_hello
say_hello()  # Call 2 of say_hello
```

#### Decorator with State

```python
import time
from functools import wraps
from collections import defaultdict

class RateLimiter:
    """Rate limiting decorator."""
    
    def __init__(self, max_calls=10, time_window=60):
        self.max_calls = max_calls
        self.time_window = time_window
        self.calls = defaultdict(list)
    
    def __call__(self, func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            now = time.time()
            key = f"{func.__name__}:{id(args)}:{id(kwargs)}"
            
            # Clean old calls
            self.calls[key] = [
                call_time for call_time in self.calls[key]
                if now - call_time < self.time_window
            ]
            
            if len(self.calls[key]) >= self.max_calls:
                raise Exception(f"Rate limit exceeded for {func.__name__}")
            
            self.calls[key].append(now)
            return func(*args, **kwargs)
        return wrapper

@RateLimiter(max_calls=3, time_window=10)
def api_call():
    return "API response"

# First 3 calls work, 4th raises exception
```

#### Property Decorators Advanced

```python
class LazyProperty:
    """Lazy evaluation property descriptor."""
    
    def __init__(self, func):
        self.func = func
        self.name = func.__name__
    
    def __get__(self, obj, owner):
        if obj is None:
            return self
        
        value = self.func(obj)
        setattr(obj, self.name, value)
        return value

class Circle:
    def __init__(self, radius):
        self.radius = radius
    
    @LazyProperty
    def area(self):
        print("Calculating area...")  # Only prints once
        return 3.14159 * self.radius ** 2

circle = Circle(5)
print(circle.area)  # Calculates and caches
print(circle.area)  # Uses cached value
```

### Function Factories

```python
def create_validator(min_val=None, max_val=None, type_check=None):
    """Factory for creating validation functions."""
    
    def validator(value):
        if type_check and not isinstance(value, type_check):
            raise TypeError(f"Expected {type_check.__name__}, got {type(value).__name__}")
        
        if min_val is not None and value < min_val:
            raise ValueError(f"Value {value} is less than minimum {min_val}")
        
        if max_val is not None and value > max_val:
            raise ValueError(f"Value {value} is greater than maximum {max_val}")
        
        return True
    
    return validator

# Create specific validators
age_validator = create_validator(min_val=0, max_val=150, type_check=int)
price_validator = create_validator(min_val=0, type_check=(int, float))

age_validator(25)      # OK
# age_validator(-5)    # ValueError
# age_validator("25")  # TypeError
```

### Closures and Nonlocal

```python
def create_accumulator(initial=0):
    """Create a closure that accumulates values."""
    total = initial
    
    def accumulate(value):
        nonlocal total
        total += value
        return total
    
    def get_total():
        return total
    
    def reset():
        nonlocal total
        total = initial
    
    # Return multiple functions
    accumulate.get_total = get_total
    accumulate.reset = reset
    return accumulate

acc = create_accumulator(10)
print(acc(5))           # 15
print(acc(3))           # 18
print(acc.get_total())  # 18
acc.reset()
print(acc.get_total())  # 10
```

---

## 2. Generators and Iterators

### Advanced Generator Patterns

#### Generator Expressions vs List Comprehensions

```python
import sys

# Memory comparison
numbers = range(1000000)

# List comprehension - uses memory
squares_list = [x**2 for x in numbers]
print(f"List size: {sys.getsizeof(squares_list)} bytes")

# Generator expression - lazy evaluation
squares_gen = (x**2 for x in numbers)
print(f"Generator size: {sys.getsizeof(squares_gen)} bytes")

# Generator functions
def fibonacci_generator():
    """Infinite Fibonacci sequence generator."""
    a, b = 0, 1
    while True:
        yield a
        a, b = b, a + b

fib = fibonacci_generator()
print([next(fib) for _ in range(10)])
```

#### Generator Pipelines

```python
def read_data():
    """Simulate reading data."""
    for i in range(100):
        yield f"data_item_{i}"

def process_data(data_stream):
    """Process data stream."""
    for item in data_stream:
        yield item.upper()

def filter_data(data_stream, pattern):
    """Filter data stream."""
    for item in data_stream:
        if pattern in item:
            yield item

def batch_data(data_stream, batch_size):
    """Batch data stream."""
    batch = []
    for item in data_stream:
        batch.append(item)
        if len(batch) == batch_size:
            yield batch
            batch = []
    if batch:
        yield batch

# Create processing pipeline
pipeline = batch_data(
    filter_data(
        process_data(read_data()),
        "ITEM_1"
    ),
    batch_size=5
)

for batch in pipeline:
    print(f"Batch: {batch}")
    if len(list(pipeline)) > 2:  # Stop after a few batches
        break
```

#### Generator Send and Throw

```python
def coroutine_example():
    """Generator that can receive values."""
    print("Coroutine started")
    try:
        while True:
            value = yield
            if value is None:
                print("Received None, continuing...")
            else:
                print(f"Received: {value}")
    except GeneratorExit:
        print("Coroutine ending")
    except Exception as e:
        print(f"Exception in coroutine: {e}")

# Using the coroutine
coro = coroutine_example()
next(coro)  # Prime the generator

coro.send("Hello")
coro.send(42)
coro.send(None)

# Throw exception
try:
    coro.throw(ValueError, "Something went wrong")
except StopIteration:
    pass

coro.close()
```

### Custom Iterator Classes

```python
class RangeIterator:
    """Custom range iterator with step functionality."""
    
    def __init__(self, start, stop, step=1):
        self.current = start
        self.stop = stop
        self.step = step
    
    def __iter__(self):
        return self
    
    def __next__(self):
        if (self.step > 0 and self.current >= self.stop) or \
           (self.step < 0 and self.current <= self.stop):
            raise StopIteration
        
        value = self.current
        self.current += self.step
        return value

class InfiniteCounter:
    """Iterator that counts infinitely."""
    
    def __init__(self, start=0, step=1):
        self.value = start
        self.step = step
    
    def __iter__(self):
        return self
    
    def __next__(self):
        current = self.value
        self.value += self.step
        return current

# Usage
for i in RangeIterator(0, 10, 2):
    print(i)  # 0, 2, 4, 6, 8

counter = InfiniteCounter(10, 5)
print([next(counter) for _ in range(5)])  # [10, 15, 20, 25, 30]
```

### Itertools Advanced Patterns

```python
import itertools
from collections import deque

def sliding_window(iterable, n):
    """Create sliding window of size n."""
    it = iter(iterable)
    window = deque(itertools.islice(it, n), maxlen=n)
    if len(window) == n:
        yield tuple(window)
    for x in it:
        window.append(x)
        yield tuple(window)

def grouper(iterable, n, fillvalue=None):
    """Collect data into fixed-length chunks."""
    args = [iter(iterable)] * n
    return itertools.zip_longest(*args, fillvalue=fillvalue)

def flatten(nested_iterable):
    """Flatten nested iterables."""
    for item in nested_iterable:
        if hasattr(item, '__iter__') and not isinstance(item, (str, bytes)):
            yield from flatten(item)
        else:
            yield item

# Examples
data = [1, 2, 3, 4, 5, 6, 7, 8, 9]

# Sliding window
print(list(sliding_window(data, 3)))
# [(1, 2, 3), (2, 3, 4), (3, 4, 5), ...]

# Grouper
print(list(grouper(data, 3)))
# [(1, 2, 3), (4, 5, 6), (7, 8, 9)]

# Flatten
nested = [[1, 2], [3, [4, 5]], [6, 7, [8, [9]]]]
print(list(flatten(nested)))
# [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

---

## 3. Context Managers

### Custom Context Managers

#### Using Classes

```python
import time
import logging

class Timer:
    """Context manager for timing code execution."""
    
    def __init__(self, name="Operation"):
        self.name = name
        self.start_time = None
    
    def __enter__(self):
        self.start_time = time.time()
        print(f"Starting {self.name}...")
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        end_time = time.time()
        duration = end_time - self.start_time
        
        if exc_type is not None:
            print(f"{self.name} failed after {duration:.2f} seconds")
            print(f"Exception: {exc_type.__name__}: {exc_val}")
            return False  # Don't suppress exception
        else:
            print(f"{self.name} completed in {duration:.2f} seconds")
        
        return False

# Usage
with Timer("Database query"):
    time.sleep(0.5)  # Simulate work

class DatabaseConnection:
    """Mock database connection context manager."""
    
    def __init__(self, host, database):
        self.host = host
        self.database = database
        self.connection = None
    
    def __enter__(self):
        print(f"Connecting to {self.database} on {self.host}")
        self.connection = f"Connection to {self.database}"
        return self.connection
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        print(f"Closing connection to {self.database}")
        self.connection = None
        
        if exc_type is not None:
            print(f"Rolling back transaction due to {exc_type.__name__}")
        else:
            print("Committing transaction")

with DatabaseConnection("localhost", "mydb") as conn:
    print(f"Using {conn}")
    # Simulate database work
```

#### Using contextlib

```python
from contextlib import contextmanager, ExitStack
import tempfile
import os

@contextmanager
def temporary_directory():
    """Context manager for temporary directory."""
    temp_dir = tempfile.mkdtemp()
    try:
        yield temp_dir
    finally:
        import shutil
        shutil.rmtree(temp_dir)

@contextmanager
def suppress_stdout():
    """Context manager to suppress stdout."""
    import sys
    from io import StringIO
    
    old_stdout = sys.stdout
    sys.stdout = StringIO()
    try:
        yield
    finally:
        sys.stdout = old_stdout

@contextmanager
def change_directory(path):
    """Context manager to temporarily change directory."""
    old_cwd = os.getcwd()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(old_cwd)

# Usage examples
with temporary_directory() as temp_dir:
    print(f"Working in {temp_dir}")
    with open(os.path.join(temp_dir, "test.txt"), "w") as f:
        f.write("Hello, World!")

with suppress_stdout():
    print("This won't be printed")

print("This will be printed")
```

#### Advanced Context Manager Patterns

```python
from contextlib import ExitStack, contextmanager
import threading

class ResourcePool:
    """Context manager for resource pooling."""
    
    def __init__(self, create_resource, max_size=10):
        self.create_resource = create_resource
        self.max_size = max_size
        self.pool = []
        self.in_use = set()
        self.lock = threading.Lock()
    
    def acquire(self):
        with self.lock:
            if self.pool:
                resource = self.pool.pop()
            elif len(self.in_use) < self.max_size:
                resource = self.create_resource()
            else:
                raise Exception("No resources available")
            
            self.in_use.add(resource)
            return resource
    
    def release(self, resource):
        with self.lock:
            if resource in self.in_use:
                self.in_use.remove(resource)
                self.pool.append(resource)
    
    @contextmanager
    def get_resource(self):
        resource = self.acquire()
        try:
            yield resource
        finally:
            self.release(resource)

# Example usage
def create_connection():
    return f"Connection-{id(object())}"

pool = ResourcePool(create_connection, max_size=3)

with pool.get_resource() as conn1:
    print(f"Using {conn1}")
    with pool.get_resource() as conn2:
        print(f"Using {conn2}")

# Nested context managers with ExitStack
def process_files(filenames):
    with ExitStack() as stack:
        files = [
            stack.enter_context(open(fname, 'r'))
            for fname in filenames
        ]
        
        # Process all files
        for f in files:
            print(f"Processing {f.name}")
```

---

## 4. Advanced Object-Oriented Programming

### Multiple Inheritance and Method Resolution Order (MRO)

```python
class A:
    def method(self):
        print("A.method")
        super().method()

class B:
    def method(self):
        print("B.method")

class C(A):
    def method(self):
        print("C.method")
        super().method()

class D(A):
    def method(self):
        print("D.method")
        super().method()

class E(C, D, B):
    def method(self):
        print("E.method")
        super().method()

# Check MRO
print(E.__mro__)
# (<class '__main__.E'>, <class '__main__.C'>, <class '__main__.D'>, 
#  <class '__main__.A'>, <class '__main__.B'>, <class 'object'>)

e = E()
e.method()
# Output:
# E.method
# C.method
# D.method
# A.method
# B.method
```

### Abstract Base Classes

```python
from abc import ABC, abstractmethod, abstractproperty
from typing import Any, List

class Shape(ABC):
    """Abstract base class for shapes."""
    
    @abstractmethod
    def area(self) -> float:
        """Calculate the area of the shape."""
        pass
    
    @abstractmethod
    def perimeter(self) -> float:
        """Calculate the perimeter of the shape."""
        pass
    
    @property
    @abstractmethod
    def name(self) -> str:
        """Name of the shape."""
        pass
    
    def describe(self) -> str:
        """Describe the shape."""
        return f"{self.name}: area={self.area():.2f}, perimeter={self.perimeter():.2f}"

class Rectangle(Shape):
    def __init__(self, width: float, height: float):
        self.width = width
        self.height = height
    
    def area(self) -> float:
        return self.width * self.height
    
    def perimeter(self) -> float:
        return 2 * (self.width + self.height)
    
    @property
    def name(self) -> str:
        return "Rectangle"

class Circle(Shape):
    def __init__(self, radius: float):
        self.radius = radius
    
    def area(self) -> float:
        return 3.14159 * self.radius ** 2
    
    def perimeter(self) -> float:
        return 2 * 3.14159 * self.radius
    
    @property
    def name(self) -> str:
        return "Circle"

# Usage
shapes: List[Shape] = [
    Rectangle(5, 3),
    Circle(4)
]

for shape in shapes:
    print(shape.describe())

# Can't instantiate abstract class
# shape = Shape()  # TypeError
```

### Mixins

```python
class JsonMixin:
    """Mixin for JSON serialization."""
    
    def to_json(self):
        import json
        return json.dumps(self.__dict__)
    
    @classmethod
    def from_json(cls, json_str):
        import json
        data = json.loads(json_str)
        return cls(**data)

class ComparableMixin:
    """Mixin for comparison operations."""
    
    def __eq__(self, other):
        if not isinstance(other, self.__class__):
            return False
        return self.__dict__ == other.__dict__
    
    def __lt__(self, other):
        if not isinstance(other, self.__class__):
            return NotImplemented
        return str(self) < str(other)
    
    def __le__(self, other):
        return self < other or self == other
    
    def __gt__(self, other):
        return not self <= other
    
    def __ge__(self, other):
        return not self < other
    
    def __ne__(self, other):
        return not self == other

class Person(JsonMixin, ComparableMixin):
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def __str__(self):
        return f"{self.name}({self.age})"
    
    def __repr__(self):
        return f"Person(name='{self.name}', age={self.age})"

# Usage
person1 = Person("Alice", 30)
person2 = Person("Bob", 25)

# JSON serialization
json_data = person1.to_json()
print(json_data)

person_from_json = Person.from_json(json_data)
print(person_from_json)

# Comparison
print(person1 > person2)  # True (Alice > Bob alphabetically)
print(person1 == person_from_json)  # True
```

### Property Advanced Usage

```python
class Temperature:
    """Temperature class with validation and conversion."""
    
    def __init__(self, celsius=0):
        self._celsius = celsius
    
    @property
    def celsius(self):
        return self._celsius
    
    @celsius.setter
    def celsius(self, value):
        if value < -273.15:
            raise ValueError("Temperature cannot be below absolute zero")
        self._celsius = value
    
    @property
    def fahrenheit(self):
        return (self._celsius * 9/5) + 32
    
    @fahrenheit.setter
    def fahrenheit(self, value):
        self.celsius = (value - 32) * 5/9
    
    @property
    def kelvin(self):
        return self._celsius + 273.15
    
    @kelvin.setter
    def kelvin(self, value):
        self.celsius = value - 273.15
    
    def __str__(self):
        return f"{self.celsius}°C"

class ComputedProperty:
    """Descriptor for computed properties with caching."""
    
    def __init__(self, func):
        self.func = func
        self.name = func.__name__
        self.cache_name = f"_{self.name}_cache"
    
    def __get__(self, obj, owner):
        if obj is None:
            return self
        
        if not hasattr(obj, self.cache_name):
            value = self.func(obj)
            setattr(obj, self.cache_name, value)
        
        return getattr(obj, self.cache_name)
    
    def __set__(self, obj, value):
        setattr(obj, self.cache_name, value)
    
    def __delete__(self, obj):
        if hasattr(obj, self.cache_name):
            delattr(obj, self.cache_name)

class DataProcessor:
    def __init__(self, data):
        self.data = data
    
    @ComputedProperty
    def processed_data(self):
        print("Processing data...")  # Only called once
        return [x * 2 for x in self.data]
    
    @ComputedProperty
    def summary_stats(self):
        print("Computing stats...")  # Only called once
        data = self.processed_data
        return {
            'mean': sum(data) / len(data),
            'min': min(data),
            'max': max(data)
        }

# Usage
temp = Temperature(25)
print(f"Celsius: {temp.celsius}")
print(f"Fahrenheit: {temp.fahrenheit}")
print(f"Kelvin: {temp.kelvin}")

temp.fahrenheit = 100  # Sets celsius internally
print(f"New Celsius: {temp.celsius}")

processor = DataProcessor([1, 2, 3, 4, 5])
print(processor.processed_data)  # Processes data
print(processor.processed_data)  # Uses cache
print(processor.summary_stats)   # Computes stats (uses cached processed_data)
```

---

## 5. Metaclasses

### Understanding Metaclasses

```python
# Classes are objects too
class MyClass:
    pass

print(type(MyClass))  # <class 'type'>
print(type(type))     # <class 'type'>

# Creating classes dynamically
def init_method(self, value):
    self.value = value

def str_method(self):
    return f"DynamicClass(value={self.value})"

# Create class using type()
DynamicClass = type(
    'DynamicClass',           # Class name
    (object,),                # Base classes
    {                         # Class dictionary
        '__init__': init_method,
        '__str__': str_method,
        'class_var': 'I am dynamic'
    }
)

obj = DynamicClass(42)
print(obj)  # DynamicClass(value=42)
print(obj.class_var)  # I am dynamic
```

### Custom Metaclasses

```python
class SingletonMeta(type):
    """Metaclass that creates singleton instances."""
    
    _instances = {}
    
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*args, **kwargs)
        return cls._instances[cls]

class Database(metaclass=SingletonMeta):
    def __init__(self, connection_string="default"):
        self.connection_string = connection_string
        print(f"Creating database connection: {connection_string}")

# Both variables point to the same instance
db1 = Database("mysql://localhost")
db2 = Database("postgres://localhost")  # This won't create new instance

print(db1 is db2)  # True
print(db1.connection_string)  # mysql://localhost

class ValidatedMeta(type):
    """Metaclass that validates class definitions."""
    
    def __new__(mcs, name, bases, namespace):
        # Validate that all methods have docstrings
        for key, value in namespace.items():
            if callable(value) and not key.startswith('_'):
                if not value.__doc__:
                    raise ValueError(f"Method {key} in {name} must have a docstring")
        
        # Add automatic string representation
        if '__str__' not in namespace:
            def auto_str(self):
                attrs = ', '.join(f"{k}={v}" for k, v in self.__dict__.items())
                return f"{name}({attrs})"
            namespace['__str__'] = auto_str
        
        return super().__new__(mcs, name, bases, namespace)
    
    def __init__(cls, name, bases, namespace):
        super().__init__(name, bases, namespace)
        print(f"Created class {name} with validated methods")

class Person(metaclass=ValidatedMeta):
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def greet(self):
        """Greet someone."""
        return f"Hello, I'm {self.name}"
    
    def get_age(self):
        """Get person's age."""
        return self.age

person = Person("Alice", 30)
print(person)  # Person(name=Alice, age=30) - auto-generated __str__
```

### Metaclass for Automatic Registration

```python
class PluginMeta(type):
    """Metaclass for automatic plugin registration."""
    
    plugins = {}
    
    def __new__(mcs, name, bases, namespace):
        cls = super().__new__(mcs, name, bases, namespace)
        
        # Don't register the base class
        if name != 'Plugin':
            plugin_name = namespace.get('plugin_name', name.lower())
            mcs.plugins[plugin_name] = cls
        
        return cls

class Plugin(metaclass=PluginMeta):
    """Base plugin class."""
    plugin_name = None
    
    def execute(self):
        raise NotImplementedError

class EmailPlugin(Plugin):
    plugin_name = "email"
    
    def execute(self):
        return "Sending email..."

class DatabasePlugin(Plugin):
    plugin_name = "database"
    
    def execute(self):
        return "Connecting to database..."

class LoggingPlugin(Plugin):
    # Uses class name (lowercase) if plugin_name not specified
    
    def execute(self):
        return "Writing to log..."

# Access registered plugins
print("Registered plugins:", Plugin.plugins.keys())

# Execute plugins by name
for name, plugin_class in Plugin.plugins.items():
    plugin = plugin_class()
    print(f"{name}: {plugin.execute()}")
```

### Metaclass for ORM-like Functionality

```python
class Field:
    """Base field class."""
    
    def __init__(self, field_type, default=None, required=False):
        self.field_type = field_type
        self.default = default
        self.required = required
        self.name = None  # Will be set by metaclass
    
    def validate(self, value):
        if value is None and self.required:
            raise ValueError(f"Field {self.name} is required")
        
        if value is not None and not isinstance(value, self.field_type):
            raise TypeError(f"Field {self.name} must be of type {self.field_type.__name__}")
        
        return True

class ModelMeta(type):
    """Metaclass for ORM-like models."""
    
    def __new__(mcs, name, bases, namespace):
        # Collect fields
        fields = {}
        for key, value in list(namespace.items()):
            if isinstance(value, Field):
                value.name = key
                fields[key] = value
                del namespace[key]  # Remove from class namespace
        
        # Add fields as class attribute
        namespace['_fields'] = fields
        
        return super().__new__(mcs, name, bases, namespace)

class Model(metaclass=ModelMeta):
    """Base model class."""
    
    def __init__(self, **kwargs):
        # Set default values
        for name, field in self._fields.items():
            setattr(self, name, field.default)
        
        # Set provided values
        for name, value in kwargs.items():
            if name in self._fields:
                self._fields[name].validate(value)
                setattr(self, name, value)
            else:
                raise AttributeError(f"Unknown field: {name}")
    
    def validate(self):
        """Validate all fields."""
        for name, field in self._fields.items():
            value = getattr(self, name)
            field.validate(value)
    
    def to_dict(self):
        """Convert to dictionary."""
        return {name: getattr(self, name) for name in self._fields}

# Usage
class User(Model):
    name = Field(str, required=True)
    age = Field(int, default=0)
    email = Field(str, required=True)

# Create and validate user
user = User(name="Alice", age=30, email="alice@example.com")
user.validate()
print(user.to_dict())

# This will raise an error
# invalid_user = User(name="Bob")  # Missing required email
```

---

## 6. Descriptors

### Understanding Descriptors

```python
class LoggedAttribute:
    """Descriptor that logs attribute access."""
    
    def __init__(self, name):
        self.name = name
        self.private_name = f"_{name}"
    
    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        
        value = getattr(obj, self.private_name, None)
        print(f"Getting {self.name}: {value}")
        return value
    
    def __set__(self, obj, value):
        print(f"Setting {self.name}: {value}")
        setattr(obj, self.private_name, value)
    
    def __delete__(self, obj):
        print(f"Deleting {self.name}")
        delattr(obj, self.private_name)

class Person:
    name = LoggedAttribute("name")
    age = LoggedAttribute("age")
    
    def __init__(self, name, age):
        self.name = name  # Triggers __set__
        self.age = age    # Triggers __set__

person = Person("Alice", 30)
print(person.name)  # Triggers __get__
person.age = 31     # Triggers __set__
```

### Data Validation Descriptors

```python
class ValidatedAttribute:
    """Base descriptor for validated attributes."""
    
    def __init__(self, name):
        self.name = name
        self.private_name = f"_{name}"
    
    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        return getattr(obj, self.private_name)
    
    def __set__(self, obj, value):
        self.validate(value)
        setattr(obj, self.private_name, value)
    
    def validate(self, value):
        """Override in subclasses."""
        pass

class NonEmptyString(ValidatedAttribute):
    """Descriptor for non-empty strings."""
    
    def validate(self, value):
        if not isinstance(value, str):
            raise TypeError(f"{self.name} must be a string")
        if len(value.strip()) == 0:
            raise ValueError(f"{self.name} cannot be empty")

class PositiveNumber(ValidatedAttribute):
    """Descriptor for positive numbers."""
    
    def validate(self, value):
        if not isinstance(value, (int, float)):
            raise TypeError(f"{self.name} must be a number")
        if value <= 0:
            raise ValueError(f"{self.name} must be positive")

class EmailAddress(ValidatedAttribute):
    """Descriptor for email addresses."""
    
    def validate(self, value):
        import re
        if not isinstance(value, str):
            raise TypeError(f"{self.name} must be a string")
        
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        if not re.match(pattern, value):
            raise ValueError(f"{self.name} must be a valid email address")

class User:
    name = NonEmptyString("name")
    email = EmailAddress("email")
    age = PositiveNumber("age")
    
    def __init__(self, name, email, age):
        self.name = name
        self.email = email
        self.age = age
    
    def __repr__(self):
        return f"User(name='{self.name}', email='{self.email}', age={self.age})"

# Usage
user = User("Alice", "alice@example.com", 30)
print(user)

# These will raise validation errors:
# user.name = ""  # ValueError: name cannot be empty
# user.email = "invalid-email"  # ValueError: email must be a valid email
# user.age = -5  # ValueError: age must be positive
```

### Caching Descriptors

```python
import time
import weakref

class CachedProperty:
    """Descriptor that caches computed properties."""
    
    def __init__(self, func):
        self.func = func
        self.name = func.__name__
        self.__doc__ = func.__doc__
    
    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        
        # Use object's __dict__ to store cached value
        cache_attr = f"_cached_{self.name}"
        
        if not hasattr(obj, cache_attr):
            print(f"Computing {self.name}...")
            value = self.func(obj)
            setattr(obj, cache_attr, value)
        
        return getattr(obj, cache_attr)
    
    def __delete__(self, obj):
        # Clear cache
        cache_attr = f"_cached_{self.name}"
        if hasattr(obj, cache_attr):
            delattr(obj, cache_attr)

class TTLCachedProperty:
    """Descriptor with time-to-live caching."""
    
    def __init__(self, func, ttl=60):
        self.func = func
        self.name = func.__name__
        self.ttl = ttl
        self.__doc__ = func.__doc__
        # Use weak references to avoid memory leaks
        self.cache = weakref.WeakKeyDictionary()
    
    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        
        now = time.time()
        
        if obj in self.cache:
            value, timestamp = self.cache[obj]
            if now - timestamp < self.ttl:
                return value
        
        print(f"Computing {self.name} (TTL expired)...")
        value = self.func(obj)
        self.cache[obj] = (value, now)
        return value

class DataAnalyzer:
    def __init__(self, data):
        self.data = data
    
    @CachedProperty
    def mean(self):
        """Calculate mean (cached)."""
        return sum(self.data) / len(self.data)
    
    @CachedProperty
    def sorted_data(self):
        """Get sorted data (cached)."""
        return sorted(self.data)
    
    @TTLCachedProperty(ttl=5)
    def current_timestamp(self):
        """Get current timestamp (cached for 5 seconds)."""
        return time.time()

# Usage
analyzer = DataAnalyzer([3, 1, 4, 1, 5, 9, 2, 6])
print(analyzer.mean)         # Computes and caches
print(analyzer.mean)         # Uses cache
print(analyzer.sorted_data)  # Computes and caches

print(analyzer.current_timestamp)  # Computes
print(analyzer.current_timestamp)  # Uses cache
time.sleep(6)
print(analyzer.current_timestamp)  # Recomputes (TTL expired)
```

---

## 7. Advanced Data Structures

### Custom Collections

```python
from collections.abc import MutableSequence, MutableMapping
import bisect

class SortedList(MutableSequence):
    """A list that maintains sorted order automatically."""
    
    def __init__(self, iterable=None):
        self._items = []
        if iterable:
            for item in iterable:
                self.add(item)
    
    def __len__(self):
        return len(self._items)
    
    def __getitem__(self, index):
        return self._items[index]
    
    def __setitem__(self, index, value):
        # Remove old value and insert new one to maintain order
        old_value = self._items[index]
        del self[index]
        self.add(value)
    
    def __delitem__(self, index):
        del self._items[index]
    
    def insert(self, index, value):
        # Ignore index, insert in sorted position
        self.add(value)
    
    def add(self, value):
        """Add value in sorted position."""
        bisect.insort(self._items, value)
    
    def __repr__(self):
        return f"SortedList({self._items})"

class LRUCache(MutableMapping):
    """Least Recently Used cache implementation."""
    
    def __init__(self, maxsize=128):
        self.maxsize = maxsize
        self.data = {}
        self.access_order = []
    
    def __getitem__(self, key):
        if key not in self.data:
            raise KeyError(key)
        
        # Move to end (most recently used)
        self.access_order.remove(key)
        self.access_order.append(key)
        return self.data[key]
    
    def __setitem__(self, key, value):
        if key in self.data:
            # Update existing item
            self.access_order.remove(key)
        elif len(self.data) >= self.maxsize:
            # Remove least recently used item
            lru_key = self.access_order.pop(0)
            del self.data[lru_key]
        
        self.data[key] = value
        self.access_order.append(key)
    
    def __delitem__(self, key):
        if key not in self.data:
            raise KeyError(key)
        
        del self.data[key]
        self.access_order.remove(key)
    
    def __iter__(self):
        return iter(self.data)
    
    def __len__(self):
        return len(self.data)
    
    def __repr__(self):
        return f"LRUCache({dict(self.data)})"

# Usage
sorted_list = SortedList([3, 1, 4, 1, 5])
print(sorted_list)  # SortedList([1, 1, 3, 4, 5])

sorted_list.add(2)
print(sorted_list)  # SortedList([1, 1, 2, 3, 4, 5])

cache = LRUCache(maxsize=3)
cache['a'] = 1
cache['b'] = 2
cache['c'] = 3
cache['d'] = 4  # 'a' gets evicted

print(cache)  # LRUCache({'b': 2, 'c': 3, 'd': 4})
```

### Advanced Dictionary Patterns

```python
from collections import defaultdict, ChainMap
import json

class AttributeDict(dict):
    """Dictionary with attribute-style access."""
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.__dict__ = self
    
    def __getattr__(self, name):
        try:
            return self[name]
        except KeyError:
            raise AttributeError(f"'{type(self).__name__}' object has no attribute '{name}'")
    
    def __setattr__(self, name, value):
        self[name] = value
    
    def __delattr__(self, name):
        try:
            del self[name]
        except KeyError:
            raise AttributeError(f"'{type(self).__name__}' object has no attribute '{name}'")

class NestedDict(dict):
    """Dictionary with nested key access using dot notation."""
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        for key, value in self.items():
            if isinstance(value, dict):
                self[key] = NestedDict(value)
    
    def __setitem__(self, key, value):
        if isinstance(value, dict) and not isinstance(value, NestedDict):
            value = NestedDict(value)
        super().__setitem__(key, value)
    
    def get_nested(self, path, default=None):
        """Get value using dot-separated path."""
        keys = path.split('.')
        value = self
        
        try:
            for key in keys:
                value = value[key]
            return value
        except (KeyError, TypeError):
            return default
    
    def set_nested(self, path, value):
        """Set value using dot-separated path."""
        keys = path.split('.')
        d = self
        
        for key in keys[:-1]:
            if key not in d or not isinstance(d[key], dict):
                d[key] = NestedDict()
            d = d[key]
        
        d[keys[-1]] = value

class ObservableDict(dict):
    """Dictionary that notifies observers of changes."""
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._observers = []
    
    def add_observer(self, callback):
        """Add observer callback."""
        self._observers.append(callback)
    
    def remove_observer(self, callback):
        """Remove observer callback."""
        if callback in self._observers:
            self._observers.remove(callback)
    
    def _notify(self, action, key, value=None, old_value=None):
        """Notify all observers."""
        for callback in self._observers:
            callback(action, key, value, old_value)
    
    def __setitem__(self, key, value):
        old_value = self.get(key)
        super().__setitem__(key, value)
        action = 'update' if old_value is not None else 'create'
        self._notify(action, key, value, old_value)
    
    def __delitem__(self, key):
        old_value = self[key]
        super().__delitem__(key)
        self._notify('delete', key, None, old_value)

# Usage examples
attr_dict = AttributeDict({'name': 'Alice', 'age': 30})
print(attr_dict.name)  # Alice
attr_dict.city = 'New York'
print(attr_dict.city)  # New York

nested = NestedDict({
    'user': {
        'profile': {
            'name': 'Alice',
            'settings': {
                'theme': 'dark'
            }
        }
    }
})

print(nested.get_nested('user.profile.name'))  # Alice
nested.set_nested('user.profile.settings.language', 'en')
print(nested.get_nested('user.profile.settings.language'))  # en

# Observable dictionary
def on_change(action, key, value, old_value):
    print(f"Action: {action}, Key: {key}, Value: {value}, Old: {old_value}")

obs_dict = ObservableDict({'initial': 'value'})
obs_dict.add_observer(on_change)

obs_dict['new_key'] = 'new_value'  # Triggers observer
obs_dict['initial'] = 'updated'    # Triggers observer
del obs_dict['new_key']            # Triggers observer
```

---

## 8. Memory Management and Performance

### Memory Profiling and Optimization

```python
import sys
import tracemalloc
from functools import wraps
import time

def memory_profile(func):
    """Decorator to profile memory usage."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        tracemalloc.start()
        
        start_time = time.time()
        result = func(*args, **kwargs)
        end_time = time.time()
        
        current, peak = tracemalloc.get_traced_memory()
        tracemalloc.stop()
        
        print(f"{func.__name__}:")
        print(f"  Execution time: {end_time - start_time:.4f} seconds")
        print(f"  Current memory: {current / 1024 / 1024:.2f} MB")
        print(f"  Peak memory: {peak / 1024 / 1024:.2f} MB")
        
        return result
    return wrapper

@memory_profile
def memory_hungry_function():
    """Function that uses a lot of memory."""
    # Create large list
    large_list = [i**2 for i in range(1000000)]
    
    # Create dictionary
    large_dict = {i: str(i) for i in range(100000)}
    
    return len(large_list) + len(large_dict)

# Usage
result = memory_hungry_function()
```

### Object Size Analysis

```python
import sys
from collections import namedtuple

def get_object_size(obj, seen=None):
    """Calculate deep size of an object."""
    size = sys.getsizeof(obj)
    if seen is None:
        seen = set()
    
    obj_id = id(obj)
    if obj_id in seen:
        return 0
    
    # Mark as seen
    seen.add(obj_id)
    
    if isinstance(obj, dict):
        size += sum([get_object_size(v, seen) for v in obj.values()])
        size += sum([get_object_size(k, seen) for k in obj.keys()])
    elif hasattr(obj, '__dict__'):
        size += get_object_size(obj.__dict__, seen)
    elif hasattr(obj, '__iter__') and not isinstance(obj, (str, bytes, bytearray)):
        size += sum([get_object_size(i, seen) for i in obj])
    
    return size

# Compare different data structures
data = list(range(1000))

regular_list = data
tuple_data = tuple(data)
set_data = set(data)
dict_data = {i: i for i in data}

# Size comparison
structures = [
    ('List', regular_list),
    ('Tuple', tuple_data),
    ('Set', set_data),
    ('Dict', dict_data)
]

for name, structure in structures:
    size = get_object_size(structure)
    print(f"{name}: {size / 1024:.2f} KB")
```

### Slots for Memory Efficiency

```python
import sys

class RegularClass:
    """Regular class with __dict__."""
    
    def __init__(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

class SlottedClass:
    """Class with __slots__ for memory efficiency."""
    
    __slots__ = ['x', 'y', 'z']
    
    def __init__(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

# Memory comparison
regular_objects = [RegularClass(i, i+1, i+2) for i in range(1000)]
slotted_objects = [SlottedClass(i, i+1, i+2) for i in range(1000)]

print(f"Regular objects: {sys.getsizeof(regular_objects[0])} bytes per object")
print(f"Slotted objects: {sys.getsizeof(slotted_objects[0])} bytes per object")

# Note: __slots__ objects don't have __dict__
print(f"Regular has __dict__: {hasattr(regular_objects[0], '__dict__')}")
print(f"Slotted has __dict__: {hasattr(slotted_objects[0], '__dict__')}")
```

### Weak References

```python
import weakref
import gc

class Parent:
    def __init__(self, name):
        self.name = name
        self.children = []
    
    def add_child(self, child):
        self.children.append(child)
        child.parent = self  # This creates a circular reference

class Child:
    def __init__(self, name):
        self.name = name
        self.parent = None

class ImprovedParent:
    def __init__(self, name):
        self.name = name
        self.children = []
    
    def add_child(self, child):
        self.children.append(child)
        # Use weak reference to avoid circular reference
        child.parent = weakref.ref(self)

class ImprovedChild:
    def __init__(self, name):
        self.name = name
        self.parent = None
    
    def get_parent(self):
        if self.parent is not None:
            return self.parent()  # Call the weak reference
        return None

# Demonstration of memory leak
def create_circular_reference():
    parent = Parent("Parent")
    child = Child("Child")
    parent.add_child(child)
    return parent, child

def create_weak_reference():
    parent = ImprovedParent("Parent")
    child = ImprovedChild("Child")
    parent.add_child(child)
    return parent, child

# Test circular reference cleanup
print("Before creating circular references:")
print(f"Objects before: {len(gc.get_objects())}")

p1, c1 = create_circular_reference()
print(f"Objects after circular: {len(gc.get_objects())}")

p2, c2 = create_weak_reference()
print(f"Objects after weak ref: {len(gc.get_objects())}")

# Clear references
del p1, c1
gc.collect()
print(f"Objects after deleting circular: {len(gc.get_objects())}")

del p2, c2
gc.collect()
print(f"Objects after deleting weak ref: {len(gc.get_objects())}")

# Weak reference callbacks
def cleanup_callback(ref):
    print("Object was garbage collected!")

class Resource:
    def __init__(self, name):
        self.name = name

resource = Resource("Important Resource")
weak_ref = weakref.ref(resource, cleanup_callback)

print(f"Resource exists: {weak_ref() is not None}")
del resource
gc.collect()
print(f"Resource exists: {weak_ref() is not None}")
```

---

## 9. Concurrency - Threading

### Thread-Safe Data Structures

```python
import threading
import time
import queue
from collections import deque
from concurrent.futures import ThreadPoolExecutor

class ThreadSafeCounter:
    """Thread-safe counter using locks."""
    
    def __init__(self, initial=0):
        self._value = initial
        self._lock = threading.Lock()
    
    def increment(self, amount=1):
        with self._lock:
            self._value += amount
    
    def decrement(self, amount=1):
        with self._lock:
            self._value -= amount
    
    @property
    def value(self):
        with self._lock:
            return self._value

class ThreadSafeList:
    """Thread-safe list wrapper."""
    
    def __init__(self):
        self._list = []
        self._lock = threading.RLock()  # Reentrant lock
    
    def append(self, item):
        with self._lock:
            self._list.append(item)
    
    def extend(self, items):
        with self._lock:
            self._list.extend(items)
    
    def pop(self, index=-1):
        with self._lock:
            if self._list:
                return self._list.pop(index)
            raise IndexError("pop from empty list")
    
    def __len__(self):
        with self._lock:
            return len(self._list)
    
    def __getitem__(self, index):
        with self._lock:
            return self._list[index]
    
    def copy(self):
        with self._lock:
            return self._list.copy()

# Usage example
counter = ThreadSafeCounter()
safe_list = ThreadSafeList()

def worker(counter, safe_list, thread_id):
    for i in range(100):
        counter.increment()
        safe_list.append(f"Thread-{thread_id}-Item-{i}")
        time.sleep(0.001)  # Simulate work

# Create and start threads
threads = []
for i in range(5):
    thread = threading.Thread(target=worker, args=(counter, safe_list, i))
    threads.append(thread)
    thread.start()

# Wait for all threads to complete
for thread in threads:
    thread.join()

print(f"Final counter value: {counter.value}")
print(f"List length: {len(safe_list)}")
```

### Producer-Consumer Pattern

```python
import threading
import queue
import time
import random

class Producer:
    """Producer that generates data."""
    
    def __init__(self, queue, name):
        self.queue = queue
        self.name = name
    
    def produce(self, num_items=10):
        for i in range(num_items):
            item = f"{self.name}-item-{i}"
            
            # Simulate production time
            time.sleep(random.uniform(0.1, 0.5))
            
            self.queue.put(item)
            print(f"Produced: {item}")
        
        # Signal completion
        self.queue.put(None)

class Consumer:
    """Consumer that processes data."""
    
    def __init__(self, queue, name):
        self.queue = queue
        self.name = name
    
    def consume(self):
        while True:
            try:
                item = self.queue.get(timeout=1)
                
                if item is None:
                    print(f"{self.name}: Received shutdown signal")
                    self.queue.task_done()
                    break
                
                # Simulate processing time
                time.sleep(random.uniform(0.2, 0.8))
                
                print(f"{self.name} processed: {item}")
                self.queue.task_done()
                
            except queue.Empty:
                print(f"{self.name}: Queue is empty, shutting down")
                break

# Producer-Consumer example
def producer_consumer_example():
    # Use different queue types for different scenarios
    
    # FIFO Queue
    fifo_queue = queue.Queue(maxsize=5)
    
    # LIFO Queue (Stack)
    # lifo_queue = queue.LifoQueue(maxsize=5)
    
    # Priority Queue
    # priority_queue = queue.PriorityQueue(maxsize=5)
    
    # Create producers and consumers
    producer1 = Producer(fifo_queue, "Producer-1")
    producer2 = Producer(fifo_queue, "Producer-2")
    
    consumer1 = Consumer(fifo_queue, "Consumer-1")
    consumer2 = Consumer(fifo_queue, "Consumer-2")
    
    # Start threads
    threads = [
        threading.Thread(target=producer1.produce, args=(5,)),
        threading.Thread(target=producer2.produce, args=(5,)),
        threading.Thread(target=consumer1.consume),
        threading.Thread(target=consumer2.consume),
    ]
    
    for thread in threads:
        thread.start()
    
    # Wait for all tasks to be processed
    fifo_queue.join()
    
    print("All tasks completed")

# Run example
# producer_consumer_example()
```

### Thread Pool and Futures

```python
from concurrent.futures import ThreadPoolExecutor, as_completed
import requests
import time

def fetch_url(url):
    """Fetch URL and return response info."""
    try:
        start_time = time.time()
        response = requests.get(url, timeout=10)
        end_time = time.time()
        
        return {
            'url': url,
            'status_code': response.status_code,
            'response_time': end_time - start_time,
            'content_length': len(response.content)
        }
    except Exception as e:
        return {
            'url': url,
            'error': str(e)
        }

def concurrent_url_fetcher(urls, max_workers=5):
    """Fetch multiple URLs concurrently."""
    results = []
    
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        # Submit all tasks
        future_to_url = {
            executor.submit(fetch_url, url): url
            for url in urls
        }
        
        # Process completed tasks
        for future in as_completed(future_to_url):
            url = future_to_url[future]
            try:
                result = future.result()
                results.append(result)
                print(f"Completed: {url}")
            except Exception as e:
                print(f"Error with {url}: {e}")
                results.append({'url': url, 'error': str(e)})
    
    return results

# Example usage
urls = [
    'https://httpbin.org/delay/1',
    'https://httpbin.org/delay/2',
    'https://httpbin.org/status/200',
    'https://httpbin.org/status/404',
    'https://httpbin.org/json'
]

# Uncomment to test
# results = concurrent_url_fetcher(urls)
# for result in results:
#     print(result)
```

### Advanced Threading Patterns

```python
import threading
import time
from contextlib import contextmanager

class ReadWriteLock:
    """Reader-Writer lock implementation."""
    
    def __init__(self):
        self._read_ready = threading.Condition(threading.RLock())
        self._readers = 0
    
    def acquire_read(self):
        self._read_ready.acquire()
        try:
            self._readers += 1
        finally:
            self._read_ready.release()
    
    def release_read(self):
        self._read_ready.acquire()
        try:
            self._readers -= 1
            if self._readers == 0:
                self._read_ready.notifyAll()
        finally:
            self._read_ready.release()
    
    def acquire_write(self):
        self._read_ready.acquire()
        while self._readers > 0:
            self._read_ready.wait()
    
    def release_write(self):
        self._read_ready.release()
    
    @contextmanager
    def read_lock(self):
        self.acquire_read()
        try:
            yield
        finally:
            self.release_read()
    
    @contextmanager
    def write_lock(self):
        self.acquire_write()
        try:
            yield
        finally:
            self.release_write()

class Barrier:
    """Simple barrier implementation."""
    
    def __init__(self, num_threads):
        self.num_threads = num_threads
        self.count = 0
        self.condition = threading.Condition()
    
    def wait(self):
        with self.condition:
            self.count += 1
            if self.count == self.num_threads:
                # Last thread notifies all others
                self.condition.notify_all()
            else:
                # Wait for all threads to reach barrier
                self.condition.wait()

# Example usage of ReadWriteLock
shared_data = {'value': 0}
rw_lock = ReadWriteLock()

def reader(name, iterations):
    for i in range(iterations):
        with rw_lock.read_lock():
            value = shared_data['value']
            print(f"Reader {name}: Read value {value}")
            time.sleep(0.1)

def writer(name, iterations):
    for i in range(iterations):
        with rw_lock.write_lock():
            shared_data['value'] += 1
            print(f"Writer {name}: Wrote value {shared_data['value']}")
            time.sleep(0.2)

# Example with barrier
def worker_with_barrier(barrier, worker_id):
    print(f"Worker {worker_id}: Starting work")
    time.sleep(random.uniform(1, 3))  # Simulate work
    
    print(f"Worker {worker_id}: Waiting at barrier")
    barrier.wait()
    
    print(f"Worker {worker_id}: Past barrier, continuing")

# Uncomment to test
# barrier = Barrier(3)
# threads = [
#     threading.Thread(target=worker_with_barrier, args=(barrier, i))
#     for i in range(3)
# ]
# 
# for thread in threads:
#     thread.start()
# 
# for thread in threads:
#     thread.join()
```

---

## 10. Concurrency - Multiprocessing

### Process-Based Parallelism

```python
import multiprocessing as mp
import time
import os
from concurrent.futures import ProcessPoolExecutor

def cpu_intensive_task(n):
    """CPU-intensive task for testing multiprocessing."""
    def fibonacci(num):
        if num <= 1:
            return num
        return fibonacci(num - 1) + fibonacci(num - 2)
    
    result = fibonacci(n)
    return {
        'input': n,
        'result': result,
        'process_id': os.getpid()
    }

def compare_sequential_vs_parallel():
    """Compare sequential vs parallel execution."""
    numbers = [30, 31, 32, 33, 34]
    
    # Sequential execution
    start_time = time.time()
    sequential_results = []
    for n in numbers:
        result = cpu_intensive_task(n)
        sequential_results.append(result)
    sequential_time = time.time() - start_time
    
    # Parallel execution
    start_time = time.time()
    with ProcessPoolExecutor() as executor:
        parallel_results = list(executor.map(cpu_intensive_task, numbers))
    parallel_time = time.time() - start_time
    
    print(f"Sequential time: {sequential_time:.2f} seconds")
    print(f"Parallel time: {parallel_time:.2f} seconds")
    print(f"Speedup: {sequential_time / parallel_time:.2f}x")
    
    return sequential_results, parallel_results

# Uncomment to test
# seq_results, par_results = compare_sequential_vs_parallel()
```

### Inter-Process Communication

```python
import multiprocessing as mp
import queue
import time
import random

def producer_process(shared_queue, shared_list, shared_dict, barrier, process_id):
    """Producer process that generates data."""
    # Wait for all processes to be ready
    barrier.wait()
    
    for i in range(5):
        # Add to queue
        item = f"Process-{process_id}-Item-{i}"
        shared_queue.put(item)
        
        # Add to shared list (with lock)
        shared_list.append(item)
        
        # Update shared dictionary (with lock)
        with shared_dict.get_lock():
            shared_dict[f"key_{process_id}_{i}"] = i * process_id
        
        print(f"Producer {process_id}: Created {item}")
        time.sleep(random.uniform(0.1, 0.3))
    
    # Signal completion
    shared_queue.put(None)

def consumer_process(shared_queue, process_id):
    """Consumer process that processes data."""
    processed_count = 0
    
    while True:
        try:
            item = shared_queue.get(timeout=2)
            if item is None:
                print(f"Consumer {process_id}: Received shutdown signal")
                break
            
            # Simulate processing
            time.sleep(random.uniform(0.1, 0.5))
            processed_count += 1
            print(f"Consumer {process_id}: Processed {item}")
            
        except queue.Empty:
            print(f"Consumer {process_id}: Timeout, shutting down")
            break
    
    return processed_count

def multiprocessing_example():
    """Demonstrate multiprocessing with shared objects."""
    # Create shared objects
    manager = mp.Manager()
    shared_queue = manager.Queue()
    shared_list = manager.list()
    shared_dict = manager.dict()
    
    # Create barrier for synchronization
    num_producers = 2
    num_consumers = 2
    barrier = mp.Barrier(num_producers)
    
    # Create processes
    processes = []
    
    # Producer processes
    for i in range(num_producers):
        p = mp.Process(
            target=producer_process,
            args=(shared_queue, shared_list, shared_dict, barrier, i)
        )
        processes.append(p)
    
    # Consumer processes
    for i in range(num_consumers):
        p = mp.Process(
            target=consumer_process,
            args=(shared_queue, i)
        )
        processes.append(p)
    
    # Start all processes
    for p in processes:
        p.start()
    
    # Wait for all processes to complete
    for p in processes:
        p.join()
    
    print(f"Final shared list length: {len(shared_list)}")
    print(f"Final shared dict: {dict(shared_dict)}")

# Uncomment to test
# multiprocessing_example()
```

### Process Pools and MapReduce

```python
import multiprocessing as mp
from functools import reduce
import operator
import time

def map_function(chunk):
    """Map function: process a chunk of data."""
    process_id = mp.current_process().pid
    result = {
        'process_id': process_id,
        'chunk_size': len(chunk),
        'sum': sum(chunk),
        'squares': [x**2 for x in chunk]
    }
    time.sleep(0.1)  # Simulate processing time
    return result

def reduce_function(results):
    """Reduce function: combine results from map phase."""
    total_sum = sum(result['sum'] for result in results)
    all_squares = []
    for result in results:
        all_squares.extend(result['squares'])
    
    return {
        'total_sum': total_sum,
        'total_elements': sum(result['chunk_size'] for result in results),
        'sum_of_squares': sum(all_squares),
        'processes_used': len(set(result['process_id'] for result in results))
    }

def chunk_data(data, chunk_size):
    """Split data into chunks."""
    for i in range(0, len(data), chunk_size):
        yield data[i:i + chunk_size]

def mapreduce_example():
    """Demonstrate MapReduce pattern with multiprocessing."""
    # Create large dataset
    data = list(range(1, 10001))  # 1 to 10000
    chunk_size = 1000
    
    # Split data into chunks
    chunks = list(chunk_data(data, chunk_size))
    
    print(f"Processing {len(data)} elements in {len(chunks)} chunks")
    
    # Map phase: process chunks in parallel
    with mp.Pool() as pool:
        map_results = pool.map(map_function, chunks)
    
    # Reduce phase: combine results
    final_result = reduce_function(map_results)
    
    print("MapReduce Results:")
    print(f"  Total sum: {final_result['total_sum']}")
    print(f"  Total elements: {final_result['total_elements']}")
    print(f"  Sum of squares: {final_result['sum_of_squares']}")
    print(f"  Processes used: {final_result['processes_used']}")
    
    # Verify against sequential calculation
    expected_sum = sum(data)
    expected_sum_squares = sum(x**2 for x in data)
    
    print("\nVerification:")
    print(f"  Sum correct: {final_result['total_sum'] == expected_sum}")
    print(f"  Sum of squares correct: {final_result['sum_of_squares'] == expected_sum_squares}")

# Uncomment to test
# mapreduce_example()
```

### Advanced Process Patterns

```python
import multiprocessing as mp
import signal
import time
import logging

class GracefulWorker(mp.Process):
    """Worker process that handles graceful shutdown."""
    
    def __init__(self, task_queue, result_queue, worker_id):
        super().__init__()
        self.task_queue = task_queue
        self.result_queue = result_queue
        self.worker_id = worker_id
        self.shutdown_event = mp.Event()
    
    def run(self):
        """Main worker loop."""
        signal.signal(signal.SIGTERM, self._signal_handler)
        signal.signal(signal.SIGINT, self._signal_handler)
        
        print(f"Worker {self.worker_id} started")
        
        while not self.shutdown_event.is_set():
            try:
                # Get task with timeout
                task = self.task_queue.get(timeout=1)
                
                if task is None:  # Shutdown signal
                    break
                
                # Process task
                result = self._process_task(task)
                self.result_queue.put(result)
                
            except mp.queues.Empty:
                continue
            except Exception as e:
                error_result = {
                    'worker_id': self.worker_id,
                    'error': str(e),
                    'task': task
                }
                self.result_queue.put(error_result)
        
        print(f"Worker {self.worker_id} shutting down gracefully")
    
    def _process_task(self, task):
        """Process a single task."""
        # Simulate work
        time.sleep(task.get('duration', 0.5))
        
        return {
            'worker_id': self.worker_id,
            'task_id': task['id'],
            'result': task['value'] ** 2,
            'timestamp': time.time()
        }
    
    def _signal_handler(self, signum, frame):
        """Handle shutdown signals."""
        print(f"Worker {self.worker_id} received signal {signum}")
        self.shutdown_event.set()
    
    def shutdown(self):
        """Request graceful shutdown."""
        self.shutdown_event.set()

class ProcessManager:
    """Manager for worker processes."""
    
    def __init__(self, num_workers=4):
        self.num_workers = num_workers
        self.task_queue = mp.Queue()
        self.result_queue = mp.Queue()
        self.workers = []
    
    def start_workers(self):
        """Start all worker processes."""
        for i in range(self.num_workers):
            worker = GracefulWorker(
                self.task_queue,
                self.result_queue,
                i
            )
            worker.start()
            self.workers.append(worker)
        
        print(f"Started {len(self.workers)} workers")
    
    def add_task(self, task):
        """Add task to queue."""
        self.task_queue.put(task)
    
    def get_result(self, timeout=None):
        """Get result from queue."""
        try:
            return self.result_queue.get(timeout=timeout)
        except mp.queues.Empty:
            return None
    
    def shutdown(self):
        """Shutdown all workers gracefully."""
        print("Initiating shutdown...")
        
        # Send shutdown signals
        for _ in self.workers:
            self.task_queue.put(None)
        
        # Wait for workers to finish
        for worker in self.workers:
            worker.join(timeout=5)
            if worker.is_alive():
                print(f"Force terminating worker {worker.worker_id}")
                worker.terminate()
                worker.join()
        
        print("All workers shut down")

# Example usage
def process_management_example():
    """Demonstrate process management."""
    manager = ProcessManager(num_workers=3)
    manager.start_workers()
    
    # Add tasks
    tasks = [
        {'id': i, 'value': i, 'duration': 0.2}
        for i in range(10)
    ]
    
    for task in tasks:
        manager.add_task(task)
    
    # Collect results
    results = []
    for _ in range(len(tasks)):
        result = manager.get_result(timeout=5)
        if result:
            results.append(result)
    
    print(f"Processed {len(results)} tasks")
    for result in results:
        if 'error' in result:
            print(f"Error in task {result['task']['id']}: {result['error']}")
        else:
            print(f"Task {result['task_id']}: {result['result']} (Worker {result['worker_id']})")
    
    manager.shutdown()

# Uncomment to test
# process_management_example()
```

---

## Conclusion - Part 1

This concludes Part 1 of the Advanced Python Programming Course. You've learned about:

- **Advanced Functions and Decorators**: Function introspection, complex decorators, closures
- **Generators and Iterators**: Advanced patterns, pipelines, custom iterators
- **Context Managers**: Custom implementations, contextlib usage
- **Advanced OOP**: Multiple inheritance, abstract classes, mixins, properties
- **Metaclasses**: Understanding and creating custom metaclasses
- **Descriptors**: Data validation, caching, property management
- **Advanced Data Structures**: Custom collections, specialized dictionaries
- **Memory Management**: Profiling, optimization, weak references
- **Threading**: Thread-safe data structures, synchronization, patterns
- **Multiprocessing**: Process pools, IPC, MapReduce patterns

### Next Steps

Continue with **Part 2** which covers:
- Asynchronous Programming (async/await)
- Advanced Testing and Debugging
- Performance Optimization
- Design Patterns
- Network Programming
- Database Integration
- Security and Cryptography
- Building and Packaging
- Web Development Frameworks
- Data Science and Machine Learning Integration

These advanced concepts will help you write more efficient, maintainable, and scalable Python applications.