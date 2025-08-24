# Advanced Python Programming Course - Part 2

## Table of Contents

1. [Asynchronous Programming](#1-asynchronous-programming)
2. [Advanced Testing and Debugging](#2-advanced-testing-and-debugging)
3. [Performance Optimization](#3-performance-optimization)
4. [Design Patterns](#4-design-patterns)
5. [Network Programming](#5-network-programming)
6. [Database Integration](#6-database-integration)
7. [Security and Cryptography](#7-security-and-cryptography)
8. [Building and Packaging](#8-building-and-packaging)
9. [Web Development Frameworks](#9-web-development-frameworks)
10. [Data Science Integration](#10-data-science-integration)

---

## 1. Asynchronous Programming

### Understanding Async/Await

```python
import asyncio
import aiohttp
import time
from typing import List, Dict, Any

async def simple_async_function():
    """Basic async function."""
    print("Starting async operation")
    await asyncio.sleep(1)  # Non-blocking sleep
    print("Async operation completed")
    return "result"

async def fetch_data(url: str, session: aiohttp.ClientSession) -> Dict[str, Any]:
    """Fetch data from URL asynchronously."""
    try:
        async with session.get(url) as response:
            data = await response.json()
            return {
                'url': url,
                'status': response.status,
                'data': data
            }
    except Exception as e:
        return {
            'url': url,
            'error': str(e)
        }

async def fetch_multiple_urls(urls: List[str]) -> List[Dict[str, Any]]:
    """Fetch multiple URLs concurrently."""
    async with aiohttp.ClientSession() as session:
        tasks = [fetch_data(url, session) for url in urls]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        return results

# Example usage
async def async_example():
    urls = [
        'https://httpbin.org/delay/1',
        'https://httpbin.org/delay/2',
        'https://httpbin.org/json',
        'https://httpbin.org/status/200'
    ]
    
    start_time = time.time()
    results = await fetch_multiple_urls(urls)
    end_time = time.time()
    
    print(f"Fetched {len(results)} URLs in {end_time - start_time:.2f} seconds")
    
    for result in results:
        if isinstance(result, Exception):
            print(f"Error: {result}")
        elif 'error' in result:
            print(f"Failed to fetch {result['url']}: {result['error']}")
        else:
            print(f"Successfully fetched {result['url']}: Status {result['status']}")

# Run the example
# asyncio.run(async_example())
```

### Async Context Managers and Iterators

```python
import asyncio
import aiofiles
from typing import AsyncGenerator, AsyncIterator

class AsyncContextManager:
    """Custom async context manager."""
    
    def __init__(self, resource_name: str):
        self.resource_name = resource_name
        self.resource = None
    
    async def __aenter__(self):
        print(f"Acquiring resource: {self.resource_name}")
        await asyncio.sleep(0.1)  # Simulate async setup
        self.resource = f"Resource-{self.resource_name}"
        return self.resource
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        print(f"Releasing resource: {self.resource_name}")
        await asyncio.sleep(0.1)  # Simulate async cleanup
        self.resource = None
        return False

class AsyncFileProcessor:
    """Async file processing example."""
    
    def __init__(self, filename: str):
        self.filename = filename
    
    async def __aiter__(self) -> AsyncIterator[str]:
        async with aiofiles.open(self.filename, 'r') as file:
            async for line in file:
                yield line.strip()
    
    async def process_lines(self) -> List[str]:
        processed_lines = []
        async for line in self:
            # Simulate async processing
            await asyncio.sleep(0.01)
            processed_lines.append(line.upper())
        return processed_lines

async def async_generator_example() -> AsyncGenerator[int, None]:
    """Async generator that yields numbers."""
    for i in range(10):
        await asyncio.sleep(0.1)  # Simulate async work
        yield i ** 2

# Usage examples
async def context_manager_example():
    async with AsyncContextManager("database") as resource:
        print(f"Using {resource}")
        await asyncio.sleep(0.5)

async def generator_example():
    print("Async generator results:")
    async for value in async_generator_example():
        print(f"Generated: {value}")

# Run examples
# asyncio.run(context_manager_example())
# asyncio.run(generator_example())
```

### Async Patterns and Synchronization

```python
import asyncio
from asyncio import Queue, Event, Semaphore, Lock
from typing import Optional
import random

class AsyncProducerConsumer:
    """Async producer-consumer pattern."""
    
    def __init__(self, queue_size: int = 10):
        self.queue: Queue = Queue(maxsize=queue_size)
        self.shutdown_event = Event()
    
    async def producer(self, producer_id: int, num_items: int):
        """Async producer."""
        for i in range(num_items):
            if self.shutdown_event.is_set():
                break
            
            item = f"Producer-{producer_id}-Item-{i}"
            await self.queue.put(item)
            print(f"Produced: {item}")
            await asyncio.sleep(random.uniform(0.1, 0.3))
        
        await self.queue.put(None)  # Shutdown signal
    
    async def consumer(self, consumer_id: int):
        """Async consumer."""
        while not self.shutdown_event.is_set():
            try:
                item = await asyncio.wait_for(self.queue.get(), timeout=1.0)
                
                if item is None:
                    print(f"Consumer-{consumer_id}: Shutdown signal received")
                    break
                
                # Simulate processing
                await asyncio.sleep(random.uniform(0.2, 0.5))
                print(f"Consumer-{consumer_id} processed: {item}")
                self.queue.task_done()
                
            except asyncio.TimeoutError:
                print(f"Consumer-{consumer_id}: Timeout, checking shutdown")
    
    async def run_simulation(self):
        """Run producer-consumer simulation."""
        # Create tasks
        tasks = [
            asyncio.create_task(self.producer(1, 5)),
            asyncio.create_task(self.producer(2, 5)),
            asyncio.create_task(self.consumer(1)),
            asyncio.create_task(self.consumer(2)),
        ]
        
        # Wait for producers to finish
        await asyncio.gather(*tasks[:2])
        
        # Wait for queue to be empty
        await self.queue.join()
        
        # Signal shutdown
        self.shutdown_event.set()
        
        # Wait for consumers to finish
        await asyncio.gather(*tasks[2:], return_exceptions=True)

class RateLimiter:
    """Async rate limiter using semaphore."""
    
    def __init__(self, max_calls: int, time_period: float):
        self.max_calls = max_calls
        self.time_period = time_period
        self.semaphore = Semaphore(max_calls)
        self.call_times = []
        self.lock = Lock()
    
    async def acquire(self):
        """Acquire rate limit permission."""
        await self.semaphore.acquire()
        
        async with self.lock:
            now = asyncio.get_event_loop().time()
            
            # Remove old calls outside time window
            self.call_times = [
                call_time for call_time in self.call_times
                if now - call_time < self.time_period
            ]
            
            # If we have room, record this call
            if len(self.call_times) < self.max_calls:
                self.call_times.append(now)
                return
            
            # Wait until we can make a call
            sleep_time = self.time_period - (now - self.call_times[0])
            if sleep_time > 0:
                await asyncio.sleep(sleep_time)
            
            self.call_times.append(now)
    
    def release(self):
        """Release rate limit permission."""
        self.semaphore.release()
    
    async def __aenter__(self):
        await self.acquire()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        self.release()

async def rate_limited_api_call(url: str, rate_limiter: RateLimiter) -> str:
    """Make rate-limited API call."""
    async with rate_limiter:
        # Simulate API call
        await asyncio.sleep(0.1)
        return f"Response from {url}"

# Example usage
async def async_patterns_example():
    # Producer-Consumer
    pc = AsyncProducerConsumer()
    await pc.run_simulation()
    
    # Rate limiting
    rate_limiter = RateLimiter(max_calls=3, time_period=2.0)
    
    tasks = [
        rate_limited_api_call(f"https://api.example.com/endpoint{i}", rate_limiter)
        for i in range(10)
    ]
    
    start_time = asyncio.get_event_loop().time()
    results = await asyncio.gather(*tasks)
    end_time = asyncio.get_event_loop().time()
    
    print(f"Completed {len(results)} rate-limited calls in {end_time - start_time:.2f} seconds")

# asyncio.run(async_patterns_example())
```

### Async Web Scraping and API Integration

```python
import asyncio
import aiohttp
from dataclasses import dataclass
from typing import List, Optional, Dict, Any
import json
from urllib.parse import urljoin, urlparse

@dataclass
class ScrapingResult:
    url: str
    status_code: Optional[int] = None
    content: Optional[str] = None
    error: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None

class AsyncWebScraper:
    """Advanced async web scraper."""
    
    def __init__(self, 
                 max_concurrent: int = 10,
                 delay_between_requests: float = 0.1,
                 timeout: int = 30):
        self.max_concurrent = max_concurrent
        self.delay_between_requests = delay_between_requests
        self.timeout = aiohttp.ClientTimeout(total=timeout)
        self.semaphore = asyncio.Semaphore(max_concurrent)
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def __aenter__(self):
        self.session = aiohttp.ClientSession(timeout=self.timeout)
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()
    
    async def scrape_url(self, url: str) -> ScrapingResult:
        """Scrape a single URL."""
        async with self.semaphore:
            try:
                await asyncio.sleep(self.delay_between_requests)
                
                async with self.session.get(url) as response:
                    content = await response.text()
                    
                    return ScrapingResult(
                        url=url,
                        status_code=response.status,
                        content=content,
                        metadata={
                            'content_type': response.headers.get('content-type'),
                            'content_length': len(content),
                            'response_time': response.headers.get('x-response-time')
                        }
                    )
            
            except Exception as e:
                return ScrapingResult(
                    url=url,
                    error=str(e)
                )
    
    async def scrape_multiple(self, urls: List[str]) -> List[ScrapingResult]:
        """Scrape multiple URLs concurrently."""
        tasks = [self.scrape_url(url) for url in urls]
        return await asyncio.gather(*tasks, return_exceptions=True)
    
    async def scrape_with_pagination(self, 
                                   base_url: str, 
                                   max_pages: int = 10) -> List[ScrapingResult]:
        """Scrape with pagination support."""
        results = []
        page = 1
        
        while page <= max_pages:
            url = f"{base_url}?page={page}"
            result = await self.scrape_url(url)
            
            if result.error or result.status_code != 200:
                break
            
            results.append(result)
            
            # Check if there's a next page (simplified logic)
            if not self._has_next_page(result.content):
                break
            
            page += 1
        
        return results
    
    def _has_next_page(self, content: str) -> bool:
        """Check if there's a next page (implement based on site structure)."""
        return 'next' in content.lower() and page < 10  # Simplified

class APIClient:
    """Async API client with authentication and error handling."""
    
    def __init__(self, base_url: str, api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def __aenter__(self):
        headers = {}
        if self.api_key:
            headers['Authorization'] = f'Bearer {self.api_key}'
        
        self.session = aiohttp.ClientSession(
            headers=headers,
            timeout=aiohttp.ClientTimeout(total=30)
        )
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()
    
    async def get(self, endpoint: str, params: Optional[Dict] = None) -> Dict[str, Any]:
        """Make GET request to API."""
        url = urljoin(self.base_url, endpoint)
        
        async with self.session.get(url, params=params) as response:
            if response.status == 200:
                return await response.json()
            else:
                error_text = await response.text()
                raise aiohttp.ClientError(
                    f"API request failed: {response.status} - {error_text}"
                )
    
    async def post(self, endpoint: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Make POST request to API."""
        url = urljoin(self.base_url, endpoint)
        
        async with self.session.post(url, json=data) as response:
            if response.status in [200, 201]:
                return await response.json()
            else:
                error_text = await response.text()
                raise aiohttp.ClientError(
                    f"API request failed: {response.status} - {error_text}"
                )
    
    async def batch_requests(self, requests: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """Make multiple API requests concurrently."""
        tasks = []
        
        for req in requests:
            method = req.get('method', 'GET').upper()
            endpoint = req['endpoint']
            
            if method == 'GET':
                task = self.get(endpoint, req.get('params'))
            elif method == 'POST':
                task = self.post(endpoint, req.get('data', {}))
            else:
                raise ValueError(f"Unsupported method: {method}")
            
            tasks.append(task)
        
        return await asyncio.gather(*tasks, return_exceptions=True)

# Example usage
async def web_scraping_example():
    urls = [
        'https://httpbin.org/json',
        'https://httpbin.org/uuid',
        'https://httpbin.org/ip',
        'https://httpbin.org/user-agent'
    ]
    
    async with AsyncWebScraper(max_concurrent=3) as scraper:
        results = await scraper.scrape_multiple(urls)
        
        for result in results:
            if isinstance(result, Exception):
                print(f"Exception: {result}")
            elif result.error:
                print(f"Error scraping {result.url}: {result.error}")
            else:
                print(f"Successfully scraped {result.url}: {result.status_code}")

async def api_client_example():
    async with APIClient('https://httpbin.org') as client:
        # Single requests
        response1 = await client.get('/json')
        print(f"JSON response: {response1}")
        
        # Batch requests
        batch_requests = [
            {'method': 'GET', 'endpoint': '/uuid'},
            {'method': 'GET', 'endpoint': '/ip'},
            {'method': 'POST', 'endpoint': '/post', 'data': {'key': 'value'}}
        ]
        
        batch_results = await client.batch_requests(batch_requests)
        
        for i, result in enumerate(batch_results):
            if isinstance(result, Exception):
                print(f"Batch request {i} failed: {result}")
            else:
                print(f"Batch request {i} succeeded")

# Run examples
# asyncio.run(web_scraping_example())
# asyncio.run(api_client_example())
```

---

## 2. Advanced Testing and Debugging

### Advanced Testing Techniques

```python
import pytest
import unittest.mock
from unittest.mock import Mock, MagicMock, patch, PropertyMock
from contextlib import contextmanager
import tempfile
import os
from typing import Generator

# Test fixtures and parametrization
@pytest.fixture
def sample_data():
    """Pytest fixture providing test data."""
    return {
        'users': [
            {'id': 1, 'name': 'Alice', 'email': 'alice@example.com'},
            {'id': 2, 'name': 'Bob', 'email': 'bob@example.com'}
        ],
        'products': [
            {'id': 1, 'name': 'Laptop', 'price': 999.99},
            {'id': 2, 'name': 'Phone', 'price': 699.99}
        ]
    }

@pytest.fixture
def temp_directory():
    """Fixture providing temporary directory."""
    with tempfile.TemporaryDirectory() as temp_dir:
        yield temp_dir

class DatabaseService:
    """Example service for testing."""
    
    def __init__(self, connection_string: str):
        self.connection_string = connection_string
        self.connected = False
    
    def connect(self):
        """Connect to database."""
        if not self.connection_string:
            raise ValueError("No connection string provided")
        self.connected = True
        return True
    
    def get_user(self, user_id: int):
        """Get user from database."""
        if not self.connected:
            raise RuntimeError("Not connected to database")
        
        # Simulate database call
        if user_id == 1:
            return {'id': 1, 'name': 'Alice', 'email': 'alice@example.com'}
        return None
    
    def create_user(self, user_data: dict):
        """Create user in database."""
        if not self.connected:
            raise RuntimeError("Not connected to database")
        
        if not user_data.get('name') or not user_data.get('email'):
            raise ValueError("Name and email are required")
        
        # Simulate user creation
        return {'id': 123, **user_data}

class TestDatabaseService:
    """Test class with various testing techniques."""
    
    def test_connection_success(self):
        """Test successful database connection."""
        service = DatabaseService("sqlite:///:memory:")
        assert service.connect() is True
        assert service.connected is True
    
    def test_connection_failure(self):
        """Test database connection failure."""
        service = DatabaseService("")
        
        with pytest.raises(ValueError, match="No connection string provided"):
            service.connect()
    
    @pytest.mark.parametrize("user_id,expected", [
        (1, {'id': 1, 'name': 'Alice', 'email': 'alice@example.com'}),
        (999, None)
    ])
    def test_get_user_parametrized(self, user_id, expected):
        """Parametrized test for user retrieval."""
        service = DatabaseService("sqlite:///:memory:")
        service.connect()
        
        result = service.get_user(user_id)
        assert result == expected
    
    def test_get_user_not_connected(self):
        """Test getting user when not connected."""
        service = DatabaseService("sqlite:///:memory:")
        
        with pytest.raises(RuntimeError, match="Not connected to database"):
            service.get_user(1)
    
    @patch.object(DatabaseService, 'connect')
    def test_with_mock_patch(self, mock_connect):
        """Test using patch decorator."""
        mock_connect.return_value = True
        
        service = DatabaseService("test://")
        result = service.connect()
        
        assert result is True
        mock_connect.assert_called_once()
    
    def test_with_context_manager_mock(self):
        """Test using patch as context manager."""
        service = DatabaseService("test://")
        
        with patch.object(service, 'connected', True):
            user = service.get_user(1)
            assert user is not None

# Advanced mocking techniques
class EmailService:
    """Email service for testing mocks."""
    
    def __init__(self, smtp_server: str):
        self.smtp_server = smtp_server
        self.connection = None
    
    def connect(self):
        """Connect to SMTP server."""
        # Simulate connection
        self.connection = f"Connection to {self.smtp_server}"
        return True
    
    def send_email(self, to: str, subject: str, body: str):
        """Send email."""
        if not self.connection:
            raise RuntimeError("Not connected to SMTP server")
        
        # Simulate sending
        return {
            'message_id': '12345',
            'status': 'sent',
            'to': to,
            'subject': subject
        }

class TestAdvancedMocking:
    """Advanced mocking examples."""
    
    def test_mock_with_side_effects(self):
        """Test mock with side effects."""
        service = EmailService("smtp.example.com")
        
        # Mock with side effect function
        def mock_send_email(to, subject, body):
            if '@invalid' in to:
                raise ValueError("Invalid email address")
            return {'message_id': '12345', 'status': 'sent'}
        
        with patch.object(service, 'send_email', side_effect=mock_send_email):
            # Valid email
            result = service.send_email('test@example.com', 'Test', 'Body')
            assert result['status'] == 'sent'
            
            # Invalid email
            with pytest.raises(ValueError):
                service.send_email('test@invalid', 'Test', 'Body')
    
    def test_mock_with_multiple_side_effects(self):
        """Test mock with sequence of side effects."""
        service = EmailService("smtp.example.com")
        
        side_effects = [
            RuntimeError("Connection failed"),
            True,
            True
        ]
        
        with patch.object(service, 'connect', side_effect=side_effects):
            # First call raises exception
            with pytest.raises(RuntimeError):
                service.connect()
            
            # Subsequent calls return True
            assert service.connect() is True
            assert service.connect() is True
    
    def test_magic_mock_features(self):
        """Test MagicMock special features."""
        mock_service = MagicMock()
        
        # Configure return values
        mock_service.connect.return_value = True
        mock_service.send_email.return_value = {'status': 'sent'}
        
        # Test the mock
        assert mock_service.connect() is True
        result = mock_service.send_email('test@example.com', 'Subject', 'Body')
        assert result['status'] == 'sent'
        
        # Verify calls
        mock_service.connect.assert_called_once()
        mock_service.send_email.assert_called_once_with(
            'test@example.com', 'Subject', 'Body'
        )
        
        # Test call count and arguments
        assert mock_service.send_email.call_count == 1
        args, kwargs = mock_service.send_email.call_args
        assert args[0] == 'test@example.com'
    
    def test_property_mock(self):
        """Test mocking properties."""
        service = EmailService("smtp.example.com")
        
        with patch.object(EmailService, 'connection', new_callable=PropertyMock) as mock_connection:
            mock_connection.return_value = "Mocked connection"
            
            # The property returns our mocked value
            assert service.connection == "Mocked connection"
            
            # Verify the property was accessed
            mock_connection.assert_called()

# Custom test utilities
class TestHelpers:
    """Custom test utilities and helpers."""
    
    @contextmanager
    def assert_raises_with_message(self, exception_class, message_pattern):
        """Custom assertion for exception messages."""
        try:
            yield
        except exception_class as e:
            if message_pattern not in str(e):
                raise AssertionError(
                    f"Expected '{message_pattern}' in exception message, "
                    f"but got '{str(e)}'"
                )
        else:
            raise AssertionError(f"Expected {exception_class.__name__} to be raised")
    
    def assert_dict_subset(self, subset, superset):
        """Assert that subset is contained in superset."""
        for key, value in subset.items():
            assert key in superset, f"Key '{key}' not found in superset"
            assert superset[key] == value, f"Value mismatch for key '{key}'"
    
    def test_custom_assertions(self):
        """Test custom assertion helpers."""
        # Test exception assertion
        with self.assert_raises_with_message(ValueError, "Invalid"):
            raise ValueError("Invalid input provided")
        
        # Test dictionary subset assertion
        superset = {'a': 1, 'b': 2, 'c': 3, 'd': 4}
        subset = {'a': 1, 'c': 3}
        
        self.assert_dict_subset(subset, superset)

# Performance testing
class TestPerformance:
    """Performance testing examples."""
    
    def test_function_performance(self, benchmark):
        """Test function performance using pytest-benchmark."""
        def fibonacci(n):
            if n <= 1:
                return n
            return fibonacci(n-1) + fibonacci(n-2)
        
        # Benchmark the function
        result = benchmark(fibonacci, 20)
        assert result == 6765
    
    @pytest.mark.timeout(5)
    def test_with_timeout(self):
        """Test with timeout to catch infinite loops."""
        import time
        time.sleep(1)  # This should pass
        # time.sleep(10)  # This would fail due to timeout
    
    def test_memory_usage(self):
        """Test memory usage of operations."""
        import tracemalloc
        
        tracemalloc.start()
        
        # Operation to test
        large_list = [i for i in range(100000)]
        
        current, peak = tracemalloc.get_traced_memory()
        tracemalloc.stop()
        
        # Assert memory usage is reasonable (example threshold)
        assert peak < 10 * 1024 * 1024  # Less than 10MB
```

### Advanced Debugging Techniques

```python
import pdb
import traceback
import logging
import functools
from typing import Any, Callable
import sys
import inspect

# Custom debugger and profiling decorators
def debug_calls(func: Callable) -> Callable:
    """Decorator to debug function calls."""
    
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        # Get caller information
        frame = inspect.currentframe().f_back
        caller_info = f"{frame.f_code.co_filename}:{frame.f_lineno}"
        
        print(f"DEBUG: Calling {func.__name__} from {caller_info}")
        print(f"DEBUG: Args: {args}")
        print(f"DEBUG: Kwargs: {kwargs}")
        
        try:
            result = func(*args, **kwargs)
            print(f"DEBUG: {func.__name__} returned: {result}")
            return result
        except Exception as e:
            print(f"DEBUG: {func.__name__} raised {type(e).__name__}: {e}")
            raise
    
    return wrapper

def trace_execution(func: Callable) -> Callable:
    """Decorator to trace function execution."""
    
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        print(f"TRACE: Entering {func.__name__}")
        
        # Set up tracing
        def trace_calls(frame, event, arg):
            if event == 'line':
                filename = frame.f_code.co_filename
                lineno = frame.f_lineno
                line = linecache.getline(filename, lineno).strip()
                print(f"TRACE: {filename}:{lineno} {line}")
            return trace_calls
        
        old_trace = sys.gettrace()
        sys.settrace(trace_calls)
        
        try:
            result = func(*args, **kwargs)
            return result
        finally:
            sys.settrace(old_trace)
            print(f"TRACE: Exiting {func.__name__}")
    
    return wrapper

class DebugContext:
    """Context manager for debugging."""
    
    def __init__(self, name: str, verbose: bool = True):
        self.name = name
        self.verbose = verbose
        self.start_time = None
    
    def __enter__(self):
        if self.verbose:
            print(f"DEBUG: Entering context '{self.name}'")
        
        self.start_time = time.time()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        end_time = time.time()
        duration = end_time - self.start_time
        
        if exc_type is not None:
            print(f"DEBUG: Context '{self.name}' failed after {duration:.4f}s")
            print(f"DEBUG: Exception: {exc_type.__name__}: {exc_val}")
            
            # Print detailed traceback
            print("DEBUG: Traceback:")
            traceback.print_tb(exc_tb)
        else:
            if self.verbose:
                print(f"DEBUG: Context '{self.name}' completed in {duration:.4f}s")
        
        return False  # Don't suppress exceptions

# Custom logging configuration
def setup_advanced_logging():
    """Set up advanced logging configuration."""
    
    class ColoredFormatter(logging.Formatter):
        """Colored log formatter."""
        
        COLORS = {
            'DEBUG': '\033[94m',    # Blue
            'INFO': '\033[92m',     # Green
            'WARNING': '\033[93m',  # Yellow
            'ERROR': '\033[91m',    # Red
            'CRITICAL': '\033[95m', # Magenta
        }
        RESET = '\033[0m'
        
        def format(self, record):
            log_color = self.COLORS.get(record.levelname, '')
            record.levelname = f"{log_color}{record.levelname}{self.RESET}"
            return super().format(record)
    
    # Create logger
    logger = logging.getLogger('advanced_debug')
    logger.setLevel(logging.DEBUG)
    
    # Console handler with colors
    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.DEBUG)
    
    colored_formatter = ColoredFormatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s'
    )
    console_handler.setFormatter(colored_formatter)
    
    # File handler for detailed logs
    file_handler = logging.FileHandler('debug.log')
    file_handler.setLevel(logging.DEBUG)
    
    file_formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(funcName)s - %(message)s'
    )
    file_handler.setFormatter(file_formatter)
    
    logger.addHandler(console_handler)
    logger.addHandler(file_handler)
    
    return logger

# Interactive debugging tools
class InteractiveDebugger:
    """Interactive debugging utilities."""
    
    @staticmethod
    def inspect_object(obj: Any, depth: int = 2):
        """Inspect object attributes and methods."""
        print(f"Inspecting {type(obj).__name__}: {obj}")
        print(f"ID: {id(obj)}")
        print(f"Type: {type(obj)}")
        
        if hasattr(obj, '__dict__'):
            print("Attributes:")
            for name, value in obj.__dict__.items():
                if depth > 0 and hasattr(value, '__dict__'):
                    print(f"  {name}: {type(value).__name__} (nested object)")
                    InteractiveDebugger.inspect_object(value, depth - 1)
                else:
                    print(f"  {name}: {value}")
        
        print("Methods:")
        for name in dir(obj):
            if not name.startswith('_'):
                attr = getattr(obj, name)
                if callable(attr):
                    print(f"  {name}(): {attr.__doc__ or 'No docstring'}")
    
    @staticmethod
    def breakpoint_with_locals():
        """Set breakpoint and print local variables."""
        frame = inspect.currentframe().f_back
        local_vars = frame.f_locals
        
        print("=== BREAKPOINT ===")
        print("Local variables:")
        for name, value in local_vars.items():
            print(f"  {name}: {value}")
        
        print("Enter 'c' to continue, 'q' to quit debugging")
        while True:
            command = input("(debug) ").strip().lower()
            if command == 'c':
                break
            elif command == 'q':
                sys.exit(0)
            elif command.startswith('p '):
                # Print variable
                var_name = command[2:]
                if var_name in local_vars:
                    print(f"{var_name}: {local_vars[var_name]}")
                else:
                    print(f"Variable '{var_name}' not found")
            else:
                print("Commands: c (continue), q (quit), p <var> (print variable)")

# Example usage of debugging tools
@debug_calls
def example_function(x: int, y: int) -> int:
    """Example function for debugging."""
    with DebugContext("calculation"):
        result = x * y + 10
        
        if result > 100:
            # Custom breakpoint
            InteractiveDebugger.breakpoint_with_locals()
        
        return result

class ExampleClass:
    """Example class for debugging."""
    
    def __init__(self, name: str):
        self.name = name
        self.value = 0
    
    @trace_execution
    def process_data(self, data: list) -> int:
        """Process data and return sum."""
        total = 0
        for item in data:
            total += item
            self.value += item
        
        return total

# Example debugging session
def debugging_example():
    """Example of debugging techniques."""
    logger = setup_advanced_logging()
    
    logger.info("Starting debugging example")
    
    try:
        # Test function debugging
        result = example_function(10, 15)
        logger.debug(f"Function result: {result}")
        
        # Test class debugging
        obj = ExampleClass("test")
        InteractiveDebugger.inspect_object(obj)
        
        data_result = obj.process_data([1, 2, 3, 4, 5])
        logger.info(f"Data processing result: {data_result}")
        
    except Exception as e:
        logger.error(f"Error in debugging example: {e}")
        logger.debug("Full traceback:", exc_info=True)

# Uncomment to run debugging example
# debugging_example()
```

---

## 3. Performance Optimization

### Profiling and Benchmarking

```python
import cProfile
import pstats
import time
import timeit
from functools import wraps
from typing import Callable, Any, Dict
import memory_profiler
import line_profiler

def profile_time(func: Callable) -> Callable:
    """Decorator to profile function execution time."""
    
    @wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.perf_counter()
        result = func(*args, **kwargs)
        end_time = time.perf_counter()
        
        print(f"{func.__name__} took {end_time - start_time:.6f} seconds")
        return result
    
    return wrapper

def profile_memory(func: Callable) -> Callable:
    """Decorator to profile memory usage."""
    
    @wraps(func)
    def wrapper(*args, **kwargs):
        import tracemalloc
        
        tracemalloc.start()
        result = func(*args, **kwargs)
        current, peak = tracemalloc.get_traced_memory()
        tracemalloc.stop()
        
        print(f"{func.__name__} used {current / 1024 / 1024:.2f} MB current, "
              f"{peak / 1024 / 1024:.2f} MB peak")
        return result
    
    return wrapper

class PerformanceProfiler:
    """Advanced performance profiling utilities."""
    
    def __init__(self):
        self.profiler = cProfile.Profile()
        self.timings: Dict[str, float] = {}
    
    def profile_function(self, func: Callable, *args, **kwargs) -> Any:
        """Profile a function call."""
        self.profiler.enable()
        result = func(*args, **kwargs)
        self.profiler.disable()
        return result
    
    def get_stats(self, sort_by: str = 'cumulative') -> str:
        """Get profiling statistics."""
        stats = pstats.Stats(self.profiler)
        stats.sort_stats(sort_by)
        
        # Capture output to string
        import io
        output = io.StringIO()
        stats.print_stats(output)
        return output.getvalue()
    
    def time_function(self, func: Callable, *args, **kwargs) -> Dict[str, Any]:
        """Time function execution with detailed statistics."""
        times = []
        
        # Run multiple times for statistical accuracy
        for _ in range(10):
            start = time.perf_counter()
            result = func(*args, **kwargs)
            end = time.perf_counter()
            times.append(end - start)
        
        return {
            'min_time': min(times),
            'max_time': max(times),
            'avg_time': sum(times) / len(times),
            'total_time': sum(times),
            'result': result
        }
    
    def compare_implementations(self, implementations: Dict[str, Callable], 
                              *args, **kwargs) -> Dict[str, Dict]:
        """Compare multiple implementations."""
        results = {}
        
        for name, func in implementations.items():
            print(f"Testing {name}...")
            results[name] = self.time_function(func, *args, **kwargs)
        
        return results

# Example optimization scenarios
class DataProcessor:
    """Example class with different optimization levels."""
    
    @profile_time
    @profile_memory
    def process_data_naive(self, data: list) -> list:
        """Naive implementation."""
        result = []
        for item in data:
            if item % 2 == 0:
                result.append(item * item)
        return result
    
    @profile_time
    @profile_memory
    def process_data_list_comp(self, data: list) -> list:
        """Optimized with list comprehension."""
        return [item * item for item in data if item % 2 == 0]
    
    @profile_time
    @profile_memory
    def process_data_generator(self, data: list):
        """Memory-optimized with generator."""
        return (item * item for item in data if item % 2 == 0)
    
    @profile_time
    @profile_memory
    def process_data_numpy(self, data: list):
        """Highly optimized with NumPy (if available)."""
        try:
            import numpy as np
            arr = np.array(data)
            even_mask = arr % 2 == 0
            return arr[even_mask] ** 2
        except ImportError:
            return self.process_data_list_comp(data)

# Caching and memoization
from functools import lru_cache, cache
import functools

class MemoizeDict:
    """Custom memoization with dictionary."""
    
    def __init__(self, max_size: int = 128):
        self.max_size = max_size
        self.cache = {}
        self.access_order = []
    
    def __call__(self, func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            # Create cache key
            key = str(args) + str(sorted(kwargs.items()))
            
            if key in self.cache:
                # Move to end (LRU)
                self.access_order.remove(key)
                self.access_order.append(key)
                return self.cache[key]
            
            # Compute result
            result = func(*args, **kwargs)
            
            # Add to cache
            if len(self.cache) >= self.max_size:
                # Remove least recently used
                oldest_key = self.access_order.pop(0)
                del self.cache[oldest_key]
            
            self.cache[key] = result
            self.access_order.append(key)
            
            return result
        
        wrapper.cache_info = lambda: {
            'cache_size': len(self.cache),
            'max_size': self.max_size
        }
        wrapper.cache_clear = lambda: self.cache.clear()
        
        return wrapper

# Example functions to optimize
def fibonacci_naive(n: int) -> int:
    """Naive recursive Fibonacci."""
    if n <= 1:
        return n
    return fibonacci_naive(n - 1) + fibonacci_naive(n - 2)

@lru_cache(maxsize=None)
def fibonacci_cached(n: int) -> int:
    """Cached recursive Fibonacci."""
    if n <= 1:
        return n
    return fibonacci_cached(n - 1) + fibonacci_cached(n - 2)

def fibonacci_iterative(n: int) -> int:
    """Iterative Fibonacci (most efficient)."""
    if n <= 1:
        return n
    
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    
    return b

@MemoizeDict(max_size=64)
def fibonacci_custom_cache(n: int) -> int:
    """Fibonacci with custom cache."""
    if n <= 1:
        return n
    return fibonacci_custom_cache(n - 1) + fibonacci_custom_cache(n - 2)

# Performance optimization example
def performance_optimization_example():
    """Example of performance optimization techniques."""
    profiler = PerformanceProfiler()
    
    # Test data processing
    processor = DataProcessor()
    test_data = list(range(100000))
    
    print("=== Data Processing Comparison ===")
    implementations = {
        'naive': processor.process_data_naive,
        'list_comp': processor.process_data_list_comp,
        'numpy': processor.process_data_numpy
    }
    
    results = profiler.compare_implementations(implementations, test_data)
    
    for name, stats in results.items():
        print(f"{name}: {stats['avg_time']:.6f}s average")
    
    # Test Fibonacci implementations
    print("\n=== Fibonacci Comparison ===")
    n = 30
    
    fib_implementations = {
        'naive': fibonacci_naive,
        'cached': fibonacci_cached,
        'iterative': fibonacci_iterative,
        'custom_cache': fibonacci_custom_cache
    }
    
    fib_results = profiler.compare_implementations(fib_implementations, n)
    
    for name, stats in fib_results.items():
        print(f"{name}: {stats['avg_time']:.6f}s average, result: {stats['result']}")
    
    # Memory usage comparison
    print("\n=== Memory Usage Analysis ===")
    
    # Large list creation
    def create_large_list():
        return [i ** 2 for i in range(1000000)]
    
    def create_large_generator():
        return (i ** 2 for i in range(1000000))
    
    print("Large list:")
    list_result = profiler.time_function(create_large_list)
    print(f"Time: {list_result['avg_time']:.6f}s")
    
    print("Large generator:")
    gen_result = profiler.time_function(create_large_generator)
    print(f"Time: {gen_result['avg_time']:.6f}s")

# Uncomment to run performance optimization example
# performance_optimization_example()
```

### Algorithmic Optimization

```python
import bisect
from collections import defaultdict, Counter, deque
from typing import List, Dict, Tuple, Set
import heapq

class AlgorithmicOptimizations:
    """Examples of algorithmic optimizations."""
    
    @staticmethod
    def find_two_sum_naive(nums: List[int], target: int) -> Tuple[int, int]:
        """Naive O(n²) solution for two sum problem."""
        for i in range(len(nums)):
            for j in range(i + 1, len(nums)):
                if nums[i] + nums[j] == target:
                    return (i, j)
        return (-1, -1)
    
    @staticmethod
    def find_two_sum_optimized(nums: List[int], target: int) -> Tuple[int, int]:
        """Optimized O(n) solution using hash map."""
        seen = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in seen:
                return (seen[complement], i)
            seen[num] = i
        return (-1, -1)
    
    @staticmethod
    def find_element_naive(sorted_list: List[int], target: int) -> int:
        """Naive O(n) linear search."""
        for i, item in enumerate(sorted_list):
            if item == target:
                return i
        return -1
    
    @staticmethod
    def find_element_optimized(sorted_list: List[int], target: int) -> int:
        """Optimized O(log n) binary search."""
        index = bisect.bisect_left(sorted_list, target)
        if index < len(sorted_list) and sorted_list[index] == target:
            return index
        return -1
    
    @staticmethod
    def group_anagrams_naive(words: List[str]) -> List[List[str]]:
        """Naive O(n² * m log m) anagram grouping."""
        groups = []
        used = set()
        
        for i, word in enumerate(words):
            if i in used:
                continue
            
            group = [word]
            sorted_word = ''.join(sorted(word))
            used.add(i)
            
            for j, other_word in enumerate(words[i+1:], i+1):
                if j not in used and ''.join(sorted(other_word)) == sorted_word:
                    group.append(other_word)
                    used.add(j)
            
            groups.append(group)
        
        return groups
    
    @staticmethod
    def group_anagrams_optimized(words: List[str]) -> List[List[str]]:
        """Optimized O(n * m log m) anagram grouping."""
        groups = defaultdict(list)
        
        for word in words:
            # Use sorted string as key
            key = ''.join(sorted(word))
            groups[key].append(word)
        
        return list(groups.values())
    
    @staticmethod
    def find_top_k_naive(nums: List[int], k: int) -> List[int]:
        """Naive O(n log n) solution for top k elements."""
        return sorted(nums, reverse=True)[:k]
    
    @staticmethod
    def find_top_k_optimized(nums: List[int], k: int) -> List[int]:
        """Optimized O(n log k) solution using min heap."""
        if k == 0:
            return []
        
        # Use min heap of size k
        heap = []
        
        for num in nums:
            if len(heap) < k:
                heapq.heappush(heap, num)
            elif num > heap[0]:
                heapq.heapreplace(heap, num)
        
        return sorted(heap, reverse=True)

class DataStructureOptimizations:
    """Examples of data structure optimizations."""
    
    def __init__(self):
        self.cache = {}
        self.graph = defaultdict(list)
    
    def build_frequency_map_naive(self, items: List[str]) -> Dict[str, int]:
        """Naive frequency counting."""
        freq_map = {}
        for item in items:
            if item in freq_map:
                freq_map[item] += 1
            else:
                freq_map[item] = 1
        return freq_map
    
    def build_frequency_map_optimized(self, items: List[str]) -> Dict[str, int]:
        """Optimized frequency counting."""
        return dict(Counter(items))
    
    def find_shortest_path_naive(self, graph: Dict[str, List[str]], 
                                start: str, end: str) -> List[str]:
        """Naive DFS path finding."""
        def dfs(current, target, path, visited):
            if current == target:
                return path + [current]
            
            visited.add(current)
            
            for neighbor in graph.get(current, []):
                if neighbor not in visited:
                    result = dfs(neighbor, target, path + [current], visited.copy())
                    if result:
                        return result
            
            return None
        
        return dfs(start, end, [], set()) or []
    
    def find_shortest_path_optimized(self, graph: Dict[str, List[str]], 
                                   start: str, end: str) -> List[str]:
        """Optimized BFS path finding."""
        if start == end:
            return [start]
        
        queue = deque([(start, [start])])
        visited = {start}
        
        while queue:
            current, path = queue.popleft()
            
            for neighbor in graph.get(current, []):
                if neighbor == end:
                    return path + [neighbor]
                
                if neighbor not in visited:
                    visited.add(neighbor)
                    queue.append((neighbor, path + [neighbor]))
        
        return []
    
    def process_stream_naive(self, data_stream) -> List[int]:
        """Naive stream processing with list."""
        results = []
        running_sum = 0
        
        for item in data_stream:
            running_sum += item
            # Keep only items above average
            if results:
                avg = sum(results) / len(results)
                results = [x for x in results if x > avg]
            results.append(running_sum)
        
        return results
    
    def process_stream_optimized(self, data_stream) -> List[int]:
        """Optimized stream processing with deque."""
        results = deque(maxlen=1000)  # Bounded memory
        running_sum = 0
        sum_total = 0
        count = 0
        
        for item in data_stream:
            running_sum += item
            
            # Efficient rolling average
            if len(results) == results.maxlen:
                # Remove oldest from average calculation
                oldest = results[0]
                sum_total -= oldest
                count -= 1
            
            results.append(running_sum)
            sum_total += running_sum
            count += 1
            
            # Efficient filtering (only when needed)
            if count > 10:  # Only filter after some items
                avg = sum_total / count
                # Update deque in place where possible
                filtered_results = deque(
                    (x for x in results if x > avg), 
                    maxlen=results.maxlen
                )
                results = filtered_results
                sum_total = sum(results)
                count = len(results)
        
        return list(results)

# Space optimization techniques
class SpaceOptimizations:
    """Examples of space optimization techniques."""
    
    @staticmethod
    def fibonacci_space_naive(n: int) -> int:
        """Naive space usage - stores all values."""
        if n <= 1:
            return n
        
        fib = [0] * (n + 1)
        fib[1] = 1
        
        for i in range(2, n + 1):
            fib[i] = fib[i-1] + fib[i-2]
        
        return fib[n]
    
    @staticmethod
    def fibonacci_space_optimized(n: int) -> int:
        """Space optimized - only stores last two values."""
        if n <= 1:
            return n
        
        prev2, prev1 = 0, 1
        
        for _ in range(2, n + 1):
            current = prev1 + prev2
            prev2, prev1 = prev1, current
        
        return prev1
    
    @staticmethod
    def matrix_multiply_naive(A: List[List[int]], B: List[List[int]]) -> List[List[int]]:
        """Naive matrix multiplication."""
        rows_A, cols_A = len(A), len(A[0])
        rows_B, cols_B = len(B), len(B[0])
        
        # Create full result matrix
        result = [[0 for _ in range(cols_B)] for _ in range(rows_A)]
        
        for i in range(rows_A):
            for j in range(cols_B):
                for k in range(cols_A):
                    result[i][j] += A[i][k] * B[k][j]
        
        return result
    
    @staticmethod
    def matrix_multiply_optimized(A: List[List[int]], B: List[List[int]]) -> List[List[int]]:
        """Space-optimized matrix multiplication with generators."""
        rows_A, cols_A = len(A), len(A[0])
        cols_B = len(B[0])
        
        def compute_row(i):
            return [
                sum(A[i][k] * B[k][j] for k in range(cols_A))
                for j in range(cols_B)
            ]
        
        # Generate rows on demand
        return [compute_row(i) for i in range(rows_A)]

# Performance testing
def optimization_performance_test():
    """Test performance of various optimizations."""
    import random
    
    # Test two sum
    nums = [random.randint(1, 1000) for _ in range(1000)]
    target = 500
    
    print("=== Two Sum Performance ===")
    
    # Time naive version
    start = time.perf_counter()
    result1 = AlgorithmicOptimizations.find_two_sum_naive(nums, target)
    naive_time = time.perf_counter() - start
    
    # Time optimized version
    start = time.perf_counter()
    result2 = AlgorithmicOptimizations.find_two_sum_optimized(nums, target)
    optimized_time = time.perf_counter() - start
    
    print(f"Naive: {naive_time:.6f}s")
    print(f"Optimized: {optimized_time:.6f}s")
    print(f"Speedup: {naive_time / optimized_time:.2f}x")
    
    # Test anagram grouping
    words = ['eat', 'tea', 'tan', 'ate', 'nat', 'bat'] * 100
    
    print("\n=== Anagram Grouping Performance ===")
    
    start = time.perf_counter()
    groups1 = AlgorithmicOptimizations.group_anagrams_naive(words)
    naive_time = time.perf_counter() - start
    
    start = time.perf_counter()
    groups2 = AlgorithmicOptimizations.group_anagrams_optimized(words)
    optimized_time = time.perf_counter() - start
    
    print(f"Naive: {naive_time:.6f}s")
    print(f"Optimized: {optimized_time:.6f}s")
    print(f"Speedup: {naive_time / optimized_time:.2f}x")
    
    # Test space optimization
    print("\n=== Space Optimization ===")
    n = 1000
    
    import sys
    
    # Measure memory usage (simplified)
    start = time.perf_counter()
    result1 = SpaceOptimizations.fibonacci_space_naive(n)
    naive_time = time.perf_counter() - start
    
    start = time.perf_counter()
    result2 = SpaceOptimizations.fibonacci_space_optimized(n)
    optimized_time = time.perf_counter() - start
    
    print(f"Naive Fibonacci: {naive_time:.6f}s")
    print(f"Optimized Fibonacci: {optimized_time:.6f}s")
    print(f"Results match: {result1 == result2}")

# Uncomment to run optimization tests
# optimization_performance_test()
```

---

## 4. Design Patterns

### Creational Patterns

```python
from abc import ABC, abstractmethod
from typing import Dict, Any, Optional, Type
import copy

# Singleton Pattern
class Singleton:
    """Thread-safe Singleton implementation."""
    
    _instances = {}
    _lock = threading.Lock()
    
    def __new__(cls, *args, **kwargs):
        if cls not in cls._instances:
            with cls._lock:
                if cls not in cls._instances:
                    cls._instances[cls] = super().__new__(cls)
        return cls._instances[cls]

class DatabaseConnection(Singleton):
    """Example singleton database connection."""
    
    def __init__(self, connection_string: str = "default"):
        if not hasattr(self, 'initialized'):
            self.connection_string = connection_string
            self.initialized = True
    
    def query(self, sql: str) -> str:
        return f"Executing: {sql}"

# Factory Pattern
class Animal(ABC):
    """Abstract animal class."""
    
    @abstractmethod
    def make_sound(self) -> str:
        pass
    
    @abstractmethod
    def get_type(self) -> str:
        pass

class Dog(Animal):
    def make_sound(self) -> str:
        return "Woof!"
    
    def get_type(self) -> str:
        return "Dog"

class Cat(Animal):
    def make_sound(self) -> str:
        return "Meow!"
    
    def get_type(self) -> str:
        return "Cat"

class AnimalFactory:
    """Factory for creating animals."""
    
    _animals: Dict[str, Type[Animal]] = {
        'dog': Dog,
        'cat': Cat
    }
    
    @classmethod
    def create_animal(cls, animal_type: str) -> Animal:
        animal_class = cls._animals.get(animal_type.lower())
        if not animal_class:
            raise ValueError(f"Unknown animal type: {animal_type}")
        return animal_class()
    
    @classmethod
    def register_animal(cls, animal_type: str, animal_class: Type[Animal]):
        """Register new animal type."""
        cls._animals[animal_type.lower()] = animal_class
    
    @classmethod
    def get_available_types(cls) -> list:
        return list(cls._animals.keys())

# Abstract Factory Pattern
class UIElement(ABC):
    """Abstract UI element."""
    
    @abstractmethod
    def render(self) -> str:
        pass

class WindowsButton(UIElement):
    def render(self) -> str:
        return "[Windows Button]"

class MacButton(UIElement):
    def render(self) -> str:
        return "(Mac Button)"

class WindowsTextBox(UIElement):
    def render(self) -> str:
        return "[Windows TextBox]"

class MacTextBox(UIElement):
    def render(self) -> str:
        return "(Mac TextBox)"

class UIFactory(ABC):
    """Abstract UI factory."""
    
    @abstractmethod
    def create_button(self) -> UIElement:
        pass
    
    @abstractmethod
    def create_textbox(self) -> UIElement:
        pass

class WindowsUIFactory(UIFactory):
    def create_button(self) -> UIElement:
        return WindowsButton()
    
    def create_textbox(self) -> UIElement:
        return WindowsTextBox()

class MacUIFactory(UIFactory):
    def create_button(self) -> UIElement:
        return MacButton()
    
    def create_textbox(self) -> UIElement:
        return MacTextBox()

# Builder Pattern
class Computer:
    """Product class for builder pattern."""
    
    def __init__(self):
        self.cpu = None
        self.memory = None
        self.storage = None
        self.graphics = None
        self.os = None
    
    def __str__(self):
        specs = []
        for attr, value in self.__dict__.items():
            if value:
                specs.append(f"{attr}: {value}")
        return f"Computer({', '.join(specs)})"

class ComputerBuilder:
    """Builder for computer configuration."""
    
    def __init__(self):
        self.computer = Computer()
    
    def add_cpu(self, cpu: str) -> 'ComputerBuilder':
        self.computer.cpu = cpu
        return self
    
    def add_memory(self, memory: str) -> 'ComputerBuilder':
        self.computer.memory = memory
        return self
    
    def add_storage(self, storage: str) -> 'ComputerBuilder':
        self.computer.storage = storage
        return self
    
    def add_graphics(self, graphics: str) -> 'ComputerBuilder':
        self.computer.graphics = graphics
        return self
    
    def add_os(self, os: str) -> 'ComputerBuilder':
        self.computer.os = os
        return self
    
    def build(self) -> Computer:
        return self.computer

class ComputerDirector:
    """Director for predefined computer configurations."""
    
    @staticmethod
    def build_gaming_computer() -> Computer:
        return (ComputerBuilder()
                .add_cpu("Intel i9")
                .add_memory("32GB DDR4")
                .add_storage("1TB NVMe SSD")
                .add_graphics("RTX 4080")
                .add_os("Windows 11")
                .build())
    
    @staticmethod
    def build_office_computer() -> Computer:
        return (ComputerBuilder()
                .add_cpu("Intel i5")
                .add_memory("16GB DDR4")
                .add_storage("512GB SSD")
                .add_os("Windows 11")
                .build())

# Prototype Pattern
class Prototype(ABC):
    """Abstract prototype class."""
    
    @abstractmethod
    def clone(self) -> 'Prototype':
        pass

class Document(Prototype):
    """Document class implementing prototype pattern."""
    
    def __init__(self, title: str, content: str, metadata: Dict[str, Any]):
        self.title = title
        self.content = content
        self.metadata = metadata
        self.creation_date = datetime.now()
    
    def clone(self) -> 'Document':
        # Deep copy to ensure independent instances
        return Document(
            title=self.title + " (Copy)",
            content=self.content,
            metadata=copy.deepcopy(self.metadata)
        )
    
    def __str__(self):
        return f"Document(title='{self.title}', created={self.creation_date})"
```

### Structural Patterns

```python
from typing import List, Union

# Adapter Pattern
class EuropeanSocket:
    """European electrical socket."""
    
    def voltage(self) -> int:
        return 230
    
    def live(self) -> int:
        return 1
    
    def neutral(self) -> int:
        return -1
    
    def earth(self) -> int:
        return 0

class USASocket:
    """USA electrical socket."""
    
    def voltage(self) -> int:
        return 120
    
    def live(self) -> int:
        return 1
    
    def neutral(self) -> int:
        return -1

class SocketAdapter:
    """Adapter to use European socket with USA appliances."""
    
    def __init__(self, european_socket: EuropeanSocket):
        self.european_socket = european_socket
    
    def voltage(self) -> int:
        # Convert European voltage to USA voltage
        return self.european_socket.voltage() // 2
    
    def live(self) -> int:
        return self.european_socket.live()
    
    def neutral(self) -> int:
        return self.european_socket.neutral()

# Decorator Pattern
class Coffee(ABC):
    """Abstract coffee component."""
    
    @abstractmethod
    def get_description(self) -> str:
        pass
    
    @abstractmethod
    def get_cost(self) -> float:
        pass

class SimpleCoffee(Coffee):
    """Basic coffee implementation."""
    
    def get_description(self) -> str:
        return "Simple coffee"
    
    def get_cost(self) -> float:
        return 2.0

class CoffeeDecorator(Coffee):
    """Base decorator for coffee."""
    
    def __init__(self, coffee: Coffee):
        self.coffee = coffee
    
    def get_description(self) -> str:
        return self.coffee.get_description()
    
    def get_cost(self) -> float:
        return self.coffee.get_cost()

class Milk(CoffeeDecorator):
    """Milk decorator."""
    
    def get_description(self) -> str:
        return self.coffee.get_description() + ", milk"
    
    def get_cost(self) -> float:
        return self.coffee.get_cost() + 0.5

class Sugar(CoffeeDecorator):
    """Sugar decorator."""
    
    def get_description(self) -> str:
        return self.coffee.get_description() + ", sugar"
    
    def get_cost(self) -> float:
        return self.coffee.get_cost() + 0.2

class WhippedCream(CoffeeDecorator):
    """Whipped cream decorator."""
    
    def get_description(self) -> str:
        return self.coffee.get_description() + ", whipped cream"
    
    def get_cost(self) -> float:
        return self.coffee.get_cost() + 0.7

# Facade Pattern
class AudioSystem:
    def turn_on(self): return "Audio system on"
    def set_volume(self, level): return f"Volume set to {level}"
    def turn_off(self): return "Audio system off"

class VideoSystem:
    def turn_on(self): return "Video system on"
    def set_resolution(self, res): return f"Resolution set to {res}"
    def turn_off(self): return "Video system off"

class LightingSystem:
    def dim_lights(self): return "Lights dimmed"
    def turn_off_lights(self): return "Lights off"

class HomeTheaterFacade:
    """Facade for home theater system."""
    
    def __init__(self):
        self.audio = AudioSystem()
        self.video = VideoSystem()
        self.lighting = LightingSystem()
    
    def watch_movie(self) -> List[str]:
        """Simplified interface for watching a movie."""
        return [
            self.lighting.dim_lights(),
            self.audio.turn_on(),
            self.audio.set_volume(8),
            self.video.turn_on(),
            self.video.set_resolution("4K"),
            "Movie started"
        ]
    
    def end_movie(self) -> List[str]:
        """Simplified interface for ending movie."""
        return [
            self.video.turn_off(),
            self.audio.turn_off(),
            self.lighting.turn_off_lights(),
            "Movie ended"
        ]

# Composite Pattern
class FileSystemItem(ABC):
    """Abstract file system item."""
    
    @abstractmethod
    def get_size(self) -> int:
        pass
    
    @abstractmethod
    def get_name(self) -> str:
        pass

class File(FileSystemItem):
    """File implementation."""
    
    def __init__(self, name: str, size: int):
        self.name = name
        self.size = size
    
    def get_size(self) -> int:
        return self.size
    
    def get_name(self) -> str:
        return self.name

class Directory(FileSystemItem):
    """Directory implementation (composite)."""
    
    def __init__(self, name: str):
        self.name = name
        self.items: List[FileSystemItem] = []
    
    def add_item(self, item: FileSystemItem):
        self.items.append(item)
    
    def remove_item(self, item: FileSystemItem):
        self.items.remove(item)
    
    def get_size(self) -> int:
        return sum(item.get_size() for item in self.items)
    
    def get_name(self) -> str:
        return self.name
    
    def list_items(self, indent: int = 0) -> str:
        result = "  " * indent + f"📁 {self.name}/ ({self.get_size()} bytes)\n"
        for item in self.items:
            if isinstance(item, Directory):
                result += item.list_items(indent + 1)
            else:
                result += "  " * (indent + 1) + f"📄 {item.get_name()} ({item.get_size()} bytes)\n"
        return result
```

### Behavioral Patterns

```python
from typing import List, Dict, Any, Callable

# Observer Pattern
class Observer(ABC):
    """Abstract observer."""
    
    @abstractmethod
    def update(self, message: str, data: Any = None):
        pass

class Subject:
    """Subject that notifies observers."""
    
    def __init__(self):
        self._observers: List[Observer] = []
        self._state = None
    
    def attach(self, observer: Observer):
        self._observers.append(observer)
    
    def detach(self, observer: Observer):
        self._observers.remove(observer)
    
    def notify(self, message: str, data: Any = None):
        for observer in self._observers:
            observer.update(message, data)
    
    def set_state(self, state: Any):
        self._state = state
        self.notify("state_changed", state)

class EmailNotifier(Observer):
    """Email notification observer."""
    
    def __init__(self, email: str):
        self.email = email
    
    def update(self, message: str, data: Any = None):
        print(f"📧 Email to {self.email}: {message} - {data}")

class SMSNotifier(Observer):
    """SMS notification observer."""
    
    def __init__(self, phone: str):
        self.phone = phone
    
    def update(self, message: str, data: Any = None):
        print(f"📱 SMS to {self.phone}: {message} - {data}")

# Strategy Pattern
class PaymentStrategy(ABC):
    """Abstract payment strategy."""
    
    @abstractmethod
    def pay(self, amount: float) -> str:
        pass

class CreditCardPayment(PaymentStrategy):
    def __init__(self, card_number: str):
        self.card_number = card_number
    
    def pay(self, amount: float) -> str:
        return f"Paid ${amount} using Credit Card ending in {self.card_number[-4:]}"

class PayPalPayment(PaymentStrategy):
    def __init__(self, email: str):
        self.email = email
    
    def pay(self, amount: float) -> str:
        return f"Paid ${amount} using PayPal account {self.email}"

class CryptoPayment(PaymentStrategy):
    def __init__(self, wallet_address: str):
        self.wallet_address = wallet_address
    
    def pay(self, amount: float) -> str:
        return f"Paid ${amount} using Crypto wallet {self.wallet_address[:8]}..."

class PaymentContext:
    """Context that uses payment strategies."""
    
    def __init__(self, strategy: PaymentStrategy):
        self.strategy = strategy
    
    def set_strategy(self, strategy: PaymentStrategy):
        self.strategy = strategy
    
    def process_payment(self, amount: float) -> str:
        return self.strategy.pay(amount)

# Command Pattern
class Command(ABC):
    """Abstract command."""
    
    @abstractmethod
    def execute(self):
        pass
    
    @abstractmethod
    def undo(self):
        pass

class Light:
    """Receiver for light commands."""
    
    def __init__(self, location: str):
        self.location = location
        self.is_on = False
    
    def turn_on(self):
        self.is_on = True
        return f"{self.location} light is ON"
    
    def turn_off(self):
        self.is_on = False
        return f"{self.location} light is OFF"

class LightOnCommand(Command):
    """Command to turn light on."""
    
    def __init__(self, light: Light):
        self.light = light
    
    def execute(self):
        return self.light.turn_on()
    
    def undo(self):
        return self.light.turn_off()

class LightOffCommand(Command):
    """Command to turn light off."""
    
    def __init__(self, light: Light):
        self.light = light
    
    def execute(self):
        return self.light.turn_off()
    
    def undo(self):
        return self.light.turn_on()

class RemoteControl:
    """Invoker for commands."""
    
    def __init__(self):
        self.commands: Dict[str, Command] = {}
        self.last_command: Optional[Command] = None
    
    def set_command(self, slot: str, command: Command):
        self.commands[slot] = command
    
    def press_button(self, slot: str) -> str:
        if slot in self.commands:
            self.last_command = self.commands[slot]
            return self.commands[slot].execute()
        return f"No command set for slot {slot}"
    
    def press_undo(self) -> str:
        if self.last_command:
            return self.last_command.undo()
        return "No command to undo"

# State Pattern
class State(ABC):
    """Abstract state."""
    
    @abstractmethod
    def handle_request(self, context: 'VendingMachine') -> str:
        pass

class IdleState(State):
    def handle_request(self, context: 'VendingMachine') -> str:
        context.set_state(WaitingForMoneyState())
        return "Please insert money"

class WaitingForMoneyState(State):
    def handle_request(self, context: 'VendingMachine') -> str:
        if context.money_inserted >= context.item_price:
            context.set_state(DispensingState())
            return "Money received. Dispensing item..."
        else:
            return f"Please insert ${context.item_price - context.money_inserted:.2f} more"

class DispensingState(State):
    def handle_request(self, context: 'VendingMachine') -> str:
        change = context.money_inserted - context.item_price
        context.money_inserted = 0
        context.set_state(IdleState())
        
        if change > 0:
            return f"Item dispensed. Change: ${change:.2f}"
        else:
            return "Item dispensed"

class VendingMachine:
    """Vending machine with state pattern."""
    
    def __init__(self, item_price: float):
        self.item_price = item_price
        self.money_inserted = 0
        self.state = IdleState()
    
    def set_state(self, state: State):
        self.state = state
    
    def insert_money(self, amount: float):
        self.money_inserted += amount
    
    def request_item(self) -> str:
        return self.state.handle_request(self)

# Example usage of design patterns
def design_patterns_example():
    """Demonstrate various design patterns."""
    
    print("=== Singleton Pattern ===")
    db1 = DatabaseConnection("mysql://localhost")
    db2 = DatabaseConnection("postgres://localhost")
    print(f"Same instance: {db1 is db2}")
    print(f"Connection string: {db1.connection_string}")
    
    print("\n=== Factory Pattern ===")
    dog = AnimalFactory.create_animal("dog")
    cat = AnimalFactory.create_animal("cat")
    print(f"Dog: {dog.make_sound()}")
    print(f"Cat: {cat.make_sound()}")
    
    print("\n=== Builder Pattern ===")
    gaming_pc = ComputerDirector.build_gaming_computer()
    office_pc = ComputerDirector.build_office_computer()
    print(f"Gaming PC: {gaming_pc}")
    print(f"Office PC: {office_pc}")
    
    print("\n=== Decorator Pattern ===")
    coffee = SimpleCoffee()
    coffee = Milk(coffee)
    coffee = Sugar(coffee)
    coffee = WhippedCream(coffee)
    print(f"Order: {coffee.get_description()}")
    print(f"Cost: ${coffee.get_cost():.2f}")
    
    print("\n=== Observer Pattern ===")
    news_agency = Subject()
    email_subscriber = EmailNotifier("user@example.com")
    sms_subscriber = SMSNotifier("+1234567890")
    
    news_agency.attach(email_subscriber)
    news_agency.attach(sms_subscriber)
    news_agency.set_state("Breaking news: Python 3.13 released!")
    
    print("\n=== Strategy Pattern ===")
    payment = PaymentContext(CreditCardPayment("1234567890123456"))
    print(payment.process_payment(100.0))
    
    payment.set_strategy(PayPalPayment("user@example.com"))
    print(payment.process_payment(50.0))
    
    print("\n=== Command Pattern ===")
    living_room_light = Light("Living Room")
    bedroom_light = Light("Bedroom")
    
    remote = RemoteControl()
    remote.set_command("1", LightOnCommand(living_room_light))
    remote.set_command("2", LightOffCommand(living_room_light))
    remote.set_command("3", LightOnCommand(bedroom_light))
    
    print(remote.press_button("1"))  # Turn on living room light
    print(remote.press_button("3"))  # Turn on bedroom light
    print(remote.press_undo())       # Undo last command
    
    print("\n=== State Pattern ===")
    vending_machine = VendingMachine(1.50)
    print(vending_machine.request_item())  # Please insert money
    
    vending_machine.insert_money(2.00)
    print(vending_machine.request_item())  # Item dispensed with change

# Uncomment to run design patterns example
# design_patterns_example()
```

---

## Conclusion - Part 2

This concludes Part 2 of the Advanced Python Programming Course. You've learned about:

- **Asynchronous Programming**: async/await, context managers, patterns, web scraping
- **Advanced Testing and Debugging**: Test techniques, mocking, profiling, debugging tools
- **Performance Optimization**: Profiling, algorithmic optimization, space optimization
- **Design Patterns**: Creational, structural, and behavioral patterns

### Next Steps for Continued Learning

1. **Practice**: Implement these patterns and techniques in real projects
2. **Explore Libraries**: Django/Flask, NumPy/Pandas, TensorFlow/PyTorch
3. **Study Frameworks**: FastAPI, Celery, SQLAlchemy
4. **Learn DevOps**: Docker, CI/CD, cloud deployment
5. **Contribute**: Open source projects, code reviews

### Additional Topics to Explore

- **Microservices Architecture**
- **Database Design and ORM Patterns**
- **API Design and Documentation**
- **Security Best Practices**
- **Machine Learning Integration**
- **Web Development Frameworks**
- **Data Engineering Pipelines**
- **Containerization and Orchestration**

The advanced Python concepts covered in both parts will help you build scalable, maintainable, and efficient applications. Remember to practice regularly and apply these concepts to real-world projects to solidify your understanding.