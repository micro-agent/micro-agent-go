# Python Project Setup Guide

## Table of Contents

1. [Project Structure](#1-project-structure)
2. [Python Environment Management](#2-python-environment-management)
3. [Dependency Management](#3-dependency-management)
4. [Project Configuration Files](#4-project-configuration-files)
5. [Development Tools Setup](#5-development-tools-setup)
6. [Testing Setup](#6-testing-setup)
7. [Code Quality Tools](#7-code-quality-tools)
8. [Documentation](#8-documentation)
9. [Version Control](#9-version-control)
10. [Continuous Integration](#10-continuous-integration)
11. [Packaging and Distribution](#11-packaging-and-distribution)
12. [Project Templates](#12-project-templates)

---

## 1. Project Structure

### Standard Python Project Layout

```
my_project/
├── README.md                 # Project description and setup instructions
├── LICENSE                   # License file
├── pyproject.toml           # Modern Python project configuration
├── requirements.txt         # Production dependencies (alternative to pyproject.toml)
├── requirements-dev.txt     # Development dependencies
├── .gitignore              # Git ignore patterns
├── .env.example            # Environment variables template
├── Makefile                # Common commands (optional)
├── docs/                   # Documentation
│   ├── conf.py
│   └── index.md
├── tests/                  # Test files
│   ├── __init__.py
│   ├── test_main.py
│   └── conftest.py
├── src/                    # Source code (recommended layout)
│   └── my_project/
│       ├── __init__.py
│       ├── main.py
│       ├── config/
│       │   ├── __init__.py
│       │   └── settings.py
│       ├── models/
│       │   ├── __init__.py
│       │   └── user.py
│       ├── services/
│       │   ├── __init__.py
│       │   └── user_service.py
│       └── utils/
│           ├── __init__.py
│           └── helpers.py
└── scripts/                # Utility scripts
    ├── setup.py
    └── deploy.sh
```

### Alternative Flat Layout (for smaller projects)

```
my_project/
├── README.md
├── pyproject.toml
├── requirements.txt
├── .gitignore
├── my_project/             # Source code directly in project root
│   ├── __init__.py
│   ├── main.py
│   ├── config.py
│   └── utils.py
├── tests/
│   ├── __init__.py
│   └── test_main.py
└── docs/
    └── README.md
```

### Creating Project Structure

```bash
# Create project directory
mkdir my_project
cd my_project

# Create subdirectories
mkdir -p src/my_project/{config,models,services,utils}
mkdir -p tests docs scripts

# Create __init__.py files
touch src/my_project/__init__.py
touch src/my_project/config/__init__.py
touch src/my_project/models/__init__.py
touch src/my_project/services/__init__.py
touch src/my_project/utils/__init__.py
touch tests/__init__.py

# Create basic files
touch README.md
touch LICENSE
touch .gitignore
touch pyproject.toml
```

---

## 2. Python Environment Management

### Using venv (Built-in)

```bash
# Create virtual environment
python -m venv venv

# Activate virtual environment
# On Windows:
venv\Scripts\activate
# On macOS/Linux:
source venv/bin/activate

# Deactivate
deactivate

# Remove virtual environment
rm -rf venv  # or rmdir /s venv on Windows
```

### Using conda

```bash
# Install Miniconda or Anaconda first

# Create environment
conda create --name my_project python=3.11

# Activate environment
conda activate my_project

# Deactivate environment
conda deactivate

# List environments
conda env list

# Remove environment
conda env remove --name my_project
```

### Using pyenv (Python Version Management)

```bash
# Install pyenv (macOS with Homebrew)
brew install pyenv

# Install specific Python version
pyenv install 3.11.0

# Set global Python version
pyenv global 3.11.0

# Set local Python version for project
pyenv local 3.11.0

# List available versions
pyenv versions
```

### Using pipenv

```bash
# Install pipenv
pip install pipenv

# Create Pipfile and virtual environment
pipenv install

# Install packages
pipenv install requests
pipenv install pytest --dev  # Development dependency

# Activate shell
pipenv shell

# Run commands in environment
pipenv run python script.py

# Install from Pipfile
pipenv install --dev  # Include dev dependencies
```

### Using poetry (Recommended for modern projects)

```bash
# Install poetry
curl -sSL https://install.python-poetry.org | python3 -

# Create new project
poetry new my_project
cd my_project

# Initialize poetry in existing project
poetry init

# Add dependencies
poetry add requests
poetry add pytest --group dev  # Development dependency

# Install dependencies
poetry install

# Activate shell
poetry shell

# Run commands
poetry run python script.py

# Build and publish
poetry build
poetry publish
```

---

## 3. Dependency Management

### requirements.txt

```txt
# Production dependencies
requests>=2.28.0,<3.0.0
fastapi>=0.95.0
uvicorn[standard]>=0.20.0
pydantic>=1.10.0
sqlalchemy>=2.0.0
alembic>=1.10.0

# Pinned versions for reproducible builds
numpy==1.24.3
pandas==2.0.1
```

### requirements-dev.txt

```txt
# Include production requirements
-r requirements.txt

# Development dependencies
pytest>=7.0.0
pytest-cov>=4.0.0
pytest-asyncio>=0.21.0
black>=23.0.0
isort>=5.12.0
flake8>=6.0.0
mypy>=1.0.0
pre-commit>=3.0.0
sphinx>=6.0.0
sphinx-rtd-theme>=1.2.0
```

### pyproject.toml (Modern approach)

```toml
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "my-project"
version = "0.1.0"
authors = [
  { name="Your Name", email="your.email@example.com" },
]
description = "A sample Python project"
readme = "README.md"
license = {file = "LICENSE"}
requires-python = ">=3.8"
classifiers = [
    "Development Status :: 3 - Alpha",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
]
keywords = ["sample", "python", "package"]

dependencies = [
    "requests>=2.28.0,<3.0.0",
    "fastapi>=0.95.0",
    "uvicorn[standard]>=0.20.0",
    "pydantic>=1.10.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "black>=23.0.0",
    "isort>=5.12.0",
    "flake8>=6.0.0",
    "mypy>=1.0.0",
]
docs = [
    "sphinx>=6.0.0",
    "sphinx-rtd-theme>=1.2.0",
]
test = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "pytest-asyncio>=0.21.0",
]

[project.urls]
"Homepage" = "https://github.com/username/my-project"
"Bug Reports" = "https://github.com/username/my-project/issues"
"Source" = "https://github.com/username/my-project"

[project.scripts]
my-project = "my_project.main:main"

[tool.hatch.build.targets.wheel]
packages = ["src/my_project"]

[tool.hatch.build.targets.sdist]
include = [
    "/src",
    "/tests",
    "/docs",
]
```

### Installing Dependencies

```bash
# Using pip
pip install -r requirements.txt
pip install -r requirements-dev.txt

# Using pip with pyproject.toml
pip install -e .  # Install in editable mode
pip install -e ".[dev]"  # Install with dev dependencies

# Using poetry
poetry install  # Install all dependencies
poetry install --only main  # Install only main dependencies
poetry install --with dev,test  # Install with specific groups

# Using conda
conda install --file requirements.txt
conda env create -f environment.yml  # If using conda environment file
```

---

## 4. Project Configuration Files

### pyproject.toml (Complete example)

```toml
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "my-project"
version = "0.1.0"
description = "A sample Python project"
authors = [{name = "Your Name", email = "your.email@example.com"}]
license = {file = "LICENSE"}
readme = "README.md"
requires-python = ">=3.8"

dependencies = [
    "requests>=2.28.0",
    "pydantic>=1.10.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "black>=23.0.0",
    "isort>=5.12.0",
    "mypy>=1.0.0",
]

# Tool configurations
[tool.black]
line-length = 88
target-version = ['py38', 'py39', 'py310', 'py311']
include = '\.pyi?$'
extend-exclude = '''
/(
  # directories
  \.eggs
  | \.git
  | \.hg
  | \.mypy_cache
  | \.tox
  | \.venv
  | build
  | dist
)/
'''

[tool.isort]
profile = "black"
multi_line_output = 3
line_length = 88
include_trailing_comma = true
force_grid_wrap = 0
use_parentheses = true
ensure_newline_before_comments = true

[tool.mypy]
python_version = "3.8"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true
disallow_incomplete_defs = true
check_untyped_defs = true
disallow_untyped_decorators = true
no_implicit_optional = true
warn_redundant_casts = true
warn_unused_ignores = true
warn_no_return = true
warn_unreachable = true
strict_equality = true

[tool.pytest.ini_options]
minversion = "6.0"
addopts = "-ra -q --strict-markers --strict-config"
testpaths = [
    "tests",
]
markers = [
    "slow: marks tests as slow",
    "integration: marks tests as integration tests",
    "unit: marks tests as unit tests",
]

[tool.coverage.run]
source = ["src"]
omit = [
    "*/tests/*",
    "*/test_*",
    "*/__pycache__/*",
]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "if self.debug:",
    "if settings.DEBUG",
    "raise AssertionError",
    "raise NotImplementedError",
    "if 0:",
    "if __name__ == .__main__.:",
    "class .*\\bProtocol\\):",
    "@(abc\\.)?abstractmethod",
]
```

### setup.cfg (Alternative configuration)

```ini
[metadata]
name = my-project
version = 0.1.0
author = Your Name
author_email = your.email@example.com
description = A sample Python project
long_description = file: README.md
long_description_content_type = text/markdown
url = https://github.com/username/my-project
classifiers =
    Development Status :: 3 - Alpha
    Intended Audience :: Developers
    License :: OSI Approved :: MIT License
    Programming Language :: Python :: 3
    Programming Language :: Python :: 3.8
    Programming Language :: Python :: 3.9
    Programming Language :: Python :: 3.10
    Programming Language :: Python :: 3.11

[options]
package_dir =
    = src
packages = find:
python_requires = >=3.8
install_requires =
    requests>=2.28.0
    pydantic>=1.10.0

[options.packages.find]
where = src

[options.extras_require]
dev =
    pytest>=7.0.0
    black>=23.0.0
    isort>=5.12.0
    mypy>=1.0.0

[flake8]
max-line-length = 88
extend-ignore = E203, W503
exclude = .git,__pycache__,docs/source/conf.py,old,build,dist

[isort]
profile = black
multi_line_output = 3
line_length = 88

[mypy]
python_version = 3.8
warn_return_any = True
warn_unused_configs = True
disallow_untyped_defs = True
```

### .env.example

```bash
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/mydatabase
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=mydatabase
DATABASE_USER=user
DATABASE_PASSWORD=password

# API Keys
API_KEY=your_api_key_here
SECRET_KEY=your_secret_key_here

# Environment
ENVIRONMENT=development
DEBUG=True
LOG_LEVEL=DEBUG

# External Services
REDIS_URL=redis://localhost:6379/0
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Application
PORT=8000
HOST=0.0.0.0
WORKERS=1
```

### .gitignore

```gitignore
# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# C extensions
*.so

# Distribution / packaging
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
share/python-wheels/
*.egg-info/
.installed.cfg
*.egg
MANIFEST

# PyInstaller
*.manifest
*.spec

# Installer logs
pip-log.txt
pip-delete-this-directory.txt

# Unit test / coverage reports
htmlcov/
.tox/
.nox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.py,cover
.hypothesis/
.pytest_cache/
cover/

# Translations
*.mo
*.pot

# Django stuff:
*.log
local_settings.py
db.sqlite3
db.sqlite3-journal

# Flask stuff:
instance/
.webassets-cache

# Scrapy stuff:
.scrapy

# Sphinx documentation
docs/_build/

# PyBuilder
.pybuilder/
target/

# Jupyter Notebook
.ipynb_checkpoints

# IPython
profile_default/
ipython_config.py

# pyenv
.python-version

# pipenv
Pipfile.lock

# poetry
poetry.lock

# pdm
.pdm.toml

# PEP 582
__pypackages__/

# Celery stuff
celerybeat-schedule
celerybeat.pid

# SageMath parsed files
*.sage.py

# Environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/

# Spyder project settings
.spyderproject
.spyproject

# Rope project settings
.ropeproject

# mkdocs documentation
/site

# mypy
.mypy_cache/
.dmypy.json
dmypy.json

# Pyre type checker
.pyre/

# pytype static type analyzer
.pytype/

# Cython debug symbols
cython_debug/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
```

---

## 5. Development Tools Setup

### Code Formatting with Black

```bash
# Install
pip install black

# Format files
black .
black src/
black file.py

# Check without formatting
black --check .

# Configuration in pyproject.toml (shown above)
```

### Import Sorting with isort

```bash
# Install
pip install isort

# Sort imports
isort .
isort src/
isort file.py

# Check without sorting
isort --check-only .

# Configuration in pyproject.toml (shown above)
```

### Linting with flake8

```bash
# Install
pip install flake8

# Run linting
flake8 .
flake8 src/

# Configuration in setup.cfg or .flake8
```

### Type Checking with mypy

```bash
# Install
pip install mypy

# Run type checking
mypy .
mypy src/

# Install type stubs
mypy --install-types

# Configuration in pyproject.toml (shown above)
```

### Pre-commit Hooks

Create `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
      - id: check-json
      - id: check-toml
      - id: check-xml
      - id: debug-statements
      - id: check-docstring-first
      - id: check-merge-conflict
      - id: check-executables-have-shebangs

  - repo: https://github.com/psf/black
    rev: 23.3.0
    hooks:
      - id: black
        language_version: python3

  - repo: https://github.com/pycqa/isort
    rev: 5.12.0
    hooks:
      - id: isort
        args: ["--profile", "black"]

  - repo: https://github.com/pycqa/flake8
    rev: 6.0.0
    hooks:
      - id: flake8

  - repo: https://github.com/pre-commit/mirrors-mypy
    rev: v1.3.0
    hooks:
      - id: mypy
        additional_dependencies: [types-requests]
```

Install and use pre-commit:

```bash
# Install
pip install pre-commit

# Install hooks
pre-commit install

# Run on all files
pre-commit run --all-files

# Update hooks
pre-commit autoupdate
```

---

## 6. Testing Setup

### Basic Test Structure

```python
# tests/conftest.py
import pytest
from pathlib import Path
import tempfile
import os

@pytest.fixture
def temp_dir():
    """Create a temporary directory for tests."""
    with tempfile.TemporaryDirectory() as temp_dir:
        yield Path(temp_dir)

@pytest.fixture
def sample_data():
    """Provide sample data for tests."""
    return {
        "users": [
            {"id": 1, "name": "Alice", "email": "alice@example.com"},
            {"id": 2, "name": "Bob", "email": "bob@example.com"}
        ]
    }

@pytest.fixture(autouse=True)
def setup_test_env():
    """Set up test environment variables."""
    os.environ["TESTING"] = "1"
    yield
    if "TESTING" in os.environ:
        del os.environ["TESTING"]
```

```python
# tests/test_main.py
import pytest
from src.my_project.main import add_numbers, process_data

class TestMathOperations:
    """Test mathematical operations."""
    
    def test_add_numbers_positive(self):
        """Test adding positive numbers."""
        result = add_numbers(2, 3)
        assert result == 5
    
    def test_add_numbers_negative(self):
        """Test adding negative numbers."""
        result = add_numbers(-2, -3)
        assert result == -5
    
    @pytest.mark.parametrize("a,b,expected", [
        (1, 2, 3),
        (0, 0, 0),
        (-1, 1, 0),
        (100, 200, 300),
    ])
    def test_add_numbers_parametrized(self, a, b, expected):
        """Test adding numbers with various inputs."""
        result = add_numbers(a, b)
        assert result == expected

class TestDataProcessing:
    """Test data processing functions."""
    
    def test_process_data_with_sample(self, sample_data):
        """Test data processing with sample data."""
        result = process_data(sample_data["users"])
        assert len(result) == 2
        assert result[0]["name"] == "Alice"
    
    def test_process_data_empty_list(self):
        """Test processing empty data."""
        result = process_data([])
        assert result == []
    
    @pytest.mark.slow
    def test_process_large_dataset(self):
        """Test processing large dataset (marked as slow)."""
        large_data = [{"id": i, "name": f"User{i}"} for i in range(10000)]
        result = process_data(large_data)
        assert len(result) == 10000

# Async tests
@pytest.mark.asyncio
class TestAsyncOperations:
    """Test asynchronous operations."""
    
    async def test_async_function(self):
        """Test async function."""
        from src.my_project.async_module import fetch_data
        result = await fetch_data("test")
        assert result is not None
```

### pytest Configuration

```ini
# pytest.ini
[tool:pytest]
minversion = 6.0
addopts = 
    -ra
    -q
    --strict-markers
    --strict-config
    --cov=src
    --cov-report=term-missing
    --cov-report=html
    --cov-report=xml
testpaths = tests
markers =
    slow: marks tests as slow (deselect with '-m "not slow"')
    integration: marks tests as integration tests
    unit: marks tests as unit tests
    asyncio: marks tests as asyncio tests
filterwarnings =
    error
    ignore::UserWarning
    ignore::DeprecationWarning
```

### Running Tests

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=src --cov-report=html

# Run specific test file
pytest tests/test_main.py

# Run specific test class
pytest tests/test_main.py::TestMathOperations

# Run specific test method
pytest tests/test_main.py::TestMathOperations::test_add_numbers_positive

# Run tests with markers
pytest -m "not slow"  # Skip slow tests
pytest -m "unit"      # Run only unit tests

# Run tests in parallel
pip install pytest-xdist
pytest -n auto  # Use all available CPUs

# Generate coverage report
pytest --cov=src --cov-report=html
# Open htmlcov/index.html in browser
```

---

## 7. Code Quality Tools

### Makefile for Common Commands

```makefile
# Makefile
.PHONY: help install install-dev test test-cov lint format type-check clean build

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install:  ## Install production dependencies
	pip install -e .

install-dev:  ## Install development dependencies
	pip install -e ".[dev,test]"
	pre-commit install

test:  ## Run tests
	pytest

test-cov:  ## Run tests with coverage
	pytest --cov=src --cov-report=html --cov-report=term

lint:  ## Run linting
	flake8 src tests
	isort --check-only src tests
	black --check src tests

format:  ## Format code
	isort src tests
	black src tests

type-check:  ## Run type checking
	mypy src

clean:  ## Clean build artifacts
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	rm -rf .pytest_cache/
	rm -rf .coverage
	rm -rf htmlcov/
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete

build:  ## Build package
	python -m build

quality:  ## Run all quality checks
	make lint
	make type-check
	make test-cov

ci:  ## Run CI pipeline locally
	make quality
	make build
```

### tox Configuration

Create `tox.ini`:

```ini
[tox]
envlist = py38,py39,py310,py311,lint,type-check
isolated_build = true

[testenv]
deps = 
    pytest
    pytest-cov
    pytest-asyncio
commands = pytest {posargs}

[testenv:lint]
deps = 
    flake8
    black
    isort
commands = 
    flake8 src tests
    black --check src tests
    isort --check-only src tests

[testenv:type-check]
deps = 
    mypy
    types-requests
commands = mypy src

[testenv:coverage]
deps = 
    pytest
    pytest-cov
commands = 
    pytest --cov=src --cov-report=html --cov-report=term

[testenv:docs]
deps = 
    sphinx
    sphinx-rtd-theme
commands = 
    sphinx-build -b html docs docs/_build/html
```

### GitHub Actions Workflow

Create `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: [3.8, 3.9, "3.10", "3.11"]

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}
    
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -e ".[dev,test]"
    
    - name: Lint with flake8
      run: |
        flake8 src tests
    
    - name: Check formatting with black
      run: |
        black --check src tests
    
    - name: Check imports with isort
      run: |
        isort --check-only src tests
    
    - name: Type check with mypy
      run: |
        mypy src
    
    - name: Test with pytest
      run: |
        pytest --cov=src --cov-report=xml
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.xml
        flags: unittests
        name: codecov-umbrella

  build:
    runs-on: ubuntu-latest
    needs: test
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: "3.11"
    
    - name: Install build dependencies
      run: |
        python -m pip install --upgrade pip
        pip install build
    
    - name: Build package
      run: python -m build
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: dist
        path: dist/
```

---

## 8. Documentation

### README.md Template

```markdown
# My Project

[![CI](https://github.com/username/my-project/workflows/CI/badge.svg)](https://github.com/username/my-project/actions)
[![Coverage](https://codecov.io/gh/username/my-project/branch/main/graph/badge.svg)](https://codecov.io/gh/username/my-project)
[![PyPI version](https://badge.fury.io/py/my-project.svg)](https://badge.fury.io/py/my-project)
[![Python versions](https://img.shields.io/pypi/pyversions/my-project.svg)](https://pypi.org/project/my-project/)

A brief description of what this project does and who it's for.

## Features

- ✨ Feature 1
- 🚀 Feature 2
- 🔒 Feature 3

## Installation

### Using pip

```bash
pip install my-project
```

### Development Installation

```bash
git clone https://github.com/username/my-project.git
cd my-project
pip install -e ".[dev]"
```

## Quick Start

```python
from my_project import main_function

result = main_function("hello")
print(result)
```

## Usage

### Basic Usage

```python
import my_project

# Example usage
client = my_project.Client(api_key="your-key")
data = client.fetch_data()
```

### Advanced Usage

```python
from my_project.advanced import AdvancedClient

# Advanced features
client = AdvancedClient(
    api_key="your-key",
    timeout=30,
    retries=3
)

async def main():
    async with client:
        data = await client.async_fetch_data()
        processed = client.process(data)
        return processed
```

## Configuration

Create a `.env` file:

```bash
API_KEY=your_api_key
DATABASE_URL=postgresql://user:pass@localhost/db
DEBUG=True
```

## API Reference

### Class: `Client`

Main client class for interacting with the API.

#### Methods

- `fetch_data(query: str) -> Dict[str, Any]`: Fetch data from API
- `process(data: Dict) -> ProcessedData`: Process raw data

## Development

### Setup

```bash
git clone https://github.com/username/my-project.git
cd my-project
make install-dev
```

### Running Tests

```bash
make test
make test-cov  # With coverage
```

### Code Quality

```bash
make lint       # Check code quality
make format     # Format code
make type-check # Type checking
make quality    # Run all checks
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for details.

## Support

- 📧 Email: support@example.com
- 💬 Discord: [Join our community](https://discord.gg/example)
- 📝 Issues: [GitHub Issues](https://github.com/username/my-project/issues)
```

### Sphinx Documentation

Create `docs/conf.py`:

```python
# Configuration file for Sphinx documentation

import os
import sys
sys.path.insert(0, os.path.abspath('../src'))

# Project information
project = 'My Project'
copyright = '2024, Your Name'
author = 'Your Name'
release = '0.1.0'

# General configuration
extensions = [
    'sphinx.ext.autodoc',
    'sphinx.ext.viewcode',
    'sphinx.ext.napoleon',
    'sphinx.ext.intersphinx',
    'sphinx.ext.todo',
]

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# HTML output options
html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']

# Napoleon settings (for Google/NumPy style docstrings)
napoleon_google_docstring = True
napoleon_numpy_docstring = True
napoleon_include_init_with_doc = False
napoleon_include_private_with_doc = False

# Autodoc settings
autodoc_default_options = {
    'members': True,
    'undoc-members': True,
    'show-inheritance': True,
}

# Intersphinx configuration
intersphinx_mapping = {
    'python': ('https://docs.python.org/3', None),
    'requests': ('https://requests.readthedocs.io/en/latest', None),
}
```

Create `docs/index.rst`:

```rst
Welcome to My Project's documentation!
======================================

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   installation
   quickstart
   api
   examples
   contributing

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
```

Build documentation:

```bash
cd docs
make html
# Open _build/html/index.html
```

---

## 9. Version Control

### Git Configuration

```bash
# Set up Git (first time)
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Initialize repository
git init
git add .
git commit -m "Initial commit"

# Add remote repository
git remote add origin https://github.com/username/my-project.git
git push -u origin main
```

### Branching Strategy

```bash
# Create and switch to development branch
git checkout -b develop

# Create feature branch
git checkout -b feature/new-feature

# Work on feature, then merge
git checkout develop
git merge feature/new-feature
git branch -d feature/new-feature

# Create release branch
git checkout -b release/v1.0.0

# Merge to main and tag
git checkout main
git merge release/v1.0.0
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin main --tags
```

### Semantic Versioning

Use semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Conventional Commits

Use conventional commit messages:

```bash
git commit -m "feat: add new authentication method"
git commit -m "fix: resolve database connection issue"
git commit -m "docs: update API documentation"
git commit -m "refactor: improve error handling"
git commit -m "test: add unit tests for user service"
git commit -m "chore: update dependencies"
```

---

## 10. Continuous Integration

### GitHub Actions (Complete Example)

`.github/workflows/ci.yml`:

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main, develop ]

env:
  PYTHON_VERSION: "3.11"

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        python-version: [3.8, 3.9, "3.10", "3.11"]
        exclude:
          - os: windows-latest
            python-version: 3.8
          - os: macos-latest
            python-version: 3.8

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}

    - name: Cache pip dependencies
      uses: actions/cache@v3
      with:
        path: ~/.cache/pip
        key: ${{ runner.os }}-pip-${{ hashFiles('**/pyproject.toml') }}
        restore-keys: |
          ${{ runner.os }}-pip-

    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -e ".[dev,test]"

    - name: Run linting
      run: |
        flake8 src tests
        black --check src tests
        isort --check-only src tests

    - name: Run type checking
      run: mypy src

    - name: Run tests
      run: |
        pytest --cov=src --cov-report=xml --cov-report=term

    - name: Upload coverage to Codecov
      if: matrix.os == 'ubuntu-latest' && matrix.python-version == '3.11'
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.xml
        flags: unittests

  security:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: ${{ env.PYTHON_VERSION }}

    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install safety bandit

    - name: Run safety check
      run: safety check --json

    - name: Run bandit security check
      run: bandit -r src/

  build:
    runs-on: ubuntu-latest
    needs: [test, security]
    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: ${{ env.PYTHON_VERSION }}

    - name: Install build dependencies
      run: |
        python -m pip install --upgrade pip
        pip install build twine

    - name: Build package
      run: python -m build

    - name: Check package
      run: twine check dist/*

    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: dist
        path: dist/

  publish:
    if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
    runs-on: ubuntu-latest
    needs: [build]
    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Download build artifacts
      uses: actions/download-artifact@v3
      with:
        name: dist
        path: dist/

    - name: Publish to PyPI
      uses: pypa/gh-action-pypi-publish@release/v1
      with:
        user: __token__
        password: ${{ secrets.PYPI_API_TOKEN }}
```

---

## 11. Packaging and Distribution

### Building Packages

```bash
# Install build tools
pip install build twine

# Build source distribution and wheel
python -m build

# Check package
twine check dist/*

# Upload to TestPyPI (for testing)
twine upload --repository testpypi dist/*

# Upload to PyPI
twine upload dist/*
```

### PyPI Configuration

Create `~/.pypirc`:

```ini
[distutils]
index-servers =
    pypi
    testpypi

[pypi]
username = __token__
password = pypi-your-api-token

[testpypi]
repository = https://test.pypi.org/legacy/
username = __token__
password = pypi-your-test-api-token
```

### Docker Support

Create `Dockerfile`:

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first for better caching
COPY pyproject.toml ./
RUN pip install --no-cache-dir -e .

# Copy source code
COPY src/ ./src/

# Create non-root user
RUN adduser --disabled-password --gecos '' appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8000

# Run application
CMD ["python", "-m", "my_project.main"]
```

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://postgres:password@db:5432/myapp
      - REDIS_URL=redis://redis:6379/0
    depends_on:
      - db
      - redis
    volumes:
      - ./src:/app/src

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## 12. Project Templates

### Using Cookiecutter

```bash
# Install cookiecutter
pip install cookiecutter

# Create project from template
cookiecutter https://github.com/audreyfeldroy/cookiecutter-pypackage

# Or use a modern template
cookiecutter https://github.com/TezRomacH/python-package-template
```

### Custom Project Template Script

Create `create_project.py`:

```python
#!/usr/bin/env python3
"""
Script to create a new Python project with best practices.
"""

import os
import sys
from pathlib import Path
import argparse

def create_project_structure(project_name: str, author: str, email: str):
    """Create project directory structure."""
    
    project_path = Path(project_name)
    if project_path.exists():
        print(f"Error: Directory {project_name} already exists!")
        return False
    
    # Create directories
    directories = [
        project_path,
        project_path / "src" / project_name,
        project_path / "src" / project_name / "config",
        project_path / "src" / project_name / "models",
        project_path / "src" / project_name / "services",
        project_path / "src" / project_name / "utils",
        project_path / "tests",
        project_path / "docs",
        project_path / "scripts",
        project_path / ".github" / "workflows",
    ]
    
    for directory in directories:
        directory.mkdir(parents=True, exist_ok=True)
    
    # Create __init__.py files
    init_files = [
        project_path / "src" / project_name / "__init__.py",
        project_path / "src" / project_name / "config" / "__init__.py",
        project_path / "src" / project_name / "models" / "__init__.py",
        project_path / "src" / project_name / "services" / "__init__.py",
        project_path / "src" / project_name / "utils" / "__init__.py",
        project_path / "tests" / "__init__.py",
    ]
    
    for init_file in init_files:
        init_file.touch()
    
    # Create basic files with content
    create_pyproject_toml(project_path, project_name, author, email)
    create_readme(project_path, project_name)
    create_gitignore(project_path)
    create_main_module(project_path, project_name)
    create_test_file(project_path, project_name)
    create_github_workflow(project_path)
    
    print(f"✅ Project {project_name} created successfully!")
    print(f"📁 Location: {project_path.absolute()}")
    print(f"\nNext steps:")
    print(f"1. cd {project_name}")
    print(f"2. python -m venv venv")
    print(f"3. source venv/bin/activate  # or venv\\Scripts\\activate on Windows")
    print(f"4. pip install -e \".[dev]\"")
    print(f"5. git init && git add . && git commit -m \"Initial commit\"")
    
    return True

def create_pyproject_toml(project_path: Path, project_name: str, author: str, email: str):
    """Create pyproject.toml file."""
    content = f'''[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "{project_name}"
version = "0.1.0"
authors = [
  {{ name="{author}", email="{email}" }},
]
description = "A Python project created with best practices"
readme = "README.md"
license = {{file = "LICENSE"}}
requires-python = ">=3.8"
classifiers = [
    "Development Status :: 3 - Alpha",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
]

dependencies = [
    "click>=8.0.0",
    "rich>=13.0.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "black>=23.0.0",
    "isort>=5.12.0",
    "flake8>=6.0.0",
    "mypy>=1.0.0",
    "pre-commit>=3.0.0",
]

[project.urls]
"Homepage" = "https://github.com/{author.lower()}/{project_name}"
"Bug Reports" = "https://github.com/{author.lower()}/{project_name}/issues"
"Source" = "https://github.com/{author.lower()}/{project_name}"

[project.scripts]
{project_name} = "{project_name}.main:main"

[tool.hatch.build.targets.wheel]
packages = ["src/{project_name}"]

[tool.black]
line-length = 88
target-version = ['py38', 'py39', 'py310', 'py311']

[tool.isort]
profile = "black"
multi_line_output = 3
line_length = 88

[tool.mypy]
python_version = "3.8"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true

[tool.pytest.ini_options]
minversion = "6.0"
addopts = "-ra -q --strict-markers --strict-config"
testpaths = ["tests"]
'''
    (project_path / "pyproject.toml").write_text(content)

def create_readme(project_path: Path, project_name: str):
    """Create README.md file."""
    content = f'''# {project_name.title()}

A Python project created with best practices.

## Installation

```bash
pip install {project_name}
```

## Development

```bash
git clone https://github.com/username/{project_name}.git
cd {project_name}
python -m venv venv
source venv/bin/activate  # or venv\\Scripts\\activate on Windows
pip install -e ".[dev]"
```

## Usage

```python
from {project_name} import hello

print(hello("World"))
```

## Testing

```bash
pytest
```

## License

MIT License
'''
    (project_path / "README.md").write_text(content)

def create_gitignore(project_path: Path):
    """Create .gitignore file."""
    content = '''# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
*.egg-info/
.installed.cfg
*.egg

# Virtual environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/

# IDE
.vscode/
.idea/
*.swp
*.swo

# Testing
.pytest_cache/
.coverage
htmlcov/
.tox/

# OS
.DS_Store
Thumbs.db
'''
    (project_path / ".gitignore").write_text(content)

def create_main_module(project_path: Path, project_name: str):
    """Create main module file."""
    content = f'''"""
Main module for {project_name}.
"""

def hello(name: str) -> str:
    """Return a greeting message.
    
    Args:
        name: Name to greet
        
    Returns:
        Greeting message
    """
    return f"Hello, {{name}}!"

def main():
    """Main entry point."""
    import click
    
    @click.command()
    @click.argument("name", default="World")
    def cli(name):
        """Simple CLI for {project_name}."""
        click.echo(hello(name))
    
    cli()

if __name__ == "__main__":
    main()
'''
    (project_path / "src" / project_name / "main.py").write_text(content)

def create_test_file(project_path: Path, project_name: str):
    """Create test file."""
    content = f'''"""
Tests for {project_name}.
"""

import pytest
from {project_name}.main import hello

def test_hello():
    """Test hello function."""
    result = hello("Test")
    assert result == "Hello, Test!"

def test_hello_default():
    """Test hello function with default."""
    result = hello("World")
    assert result == "Hello, World!"

@pytest.mark.parametrize("name,expected", [
    ("Alice", "Hello, Alice!"),
    ("Bob", "Hello, Bob!"),
    ("", "Hello, !"),
])
def test_hello_parametrized(name, expected):
    """Test hello function with various inputs."""
    result = hello(name)
    assert result == expected
'''
    (project_path / "tests" / "test_main.py").write_text(content)

def create_github_workflow(project_path: Path):
    """Create GitHub Actions workflow."""
    content = '''name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: [3.8, 3.9, "3.10", "3.11"]

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}
    
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -e ".[dev]"
    
    - name: Lint
      run: |
        flake8 src tests
        black --check src tests
        isort --check-only src tests
    
    - name: Type check
      run: mypy src
    
    - name: Test
      run: pytest --cov=src --cov-report=xml
'''
    (project_path / ".github" / "workflows" / "ci.yml").write_text(content)

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description="Create a new Python project")
    parser.add_argument("project_name", help="Name of the project")
    parser.add_argument("--author", default="Your Name", help="Author name")
    parser.add_argument("--email", default="your.email@example.com", help="Author email")
    
    args = parser.parse_args()
    
    # Validate project name
    if not args.project_name.replace("_", "").replace("-", "").isalnum():
        print("Error: Project name must contain only letters, numbers, hyphens, and underscores")
        return 1
    
    success = create_project_structure(args.project_name, args.author, args.email)
    return 0 if success else 1

if __name__ == "__main__":
    sys.exit(main())
```

Use the script:

```bash
python create_project.py my_awesome_project --author "John Doe" --email "john@example.com"
```

---

## Summary

This guide covers the complete setup of a professional Python project including:

1. **Project Structure**: Organized layout with clear separation of concerns
2. **Environment Management**: Virtual environments, version control
3. **Dependency Management**: Modern tools like Poetry, pip-tools
4. **Configuration**: pyproject.toml, setup.cfg, environment variables
5. **Development Tools**: Formatting, linting, type checking
6. **Testing**: Comprehensive test setup with pytest
7. **Code Quality**: Pre-commit hooks, CI/CD pipelines
8. **Documentation**: README, Sphinx documentation
9. **Version Control**: Git best practices, branching strategies
10. **CI/CD**: GitHub Actions, automated testing and deployment
11. **Packaging**: Building and distributing packages
12. **Templates**: Project scaffolding and automation

Following these practices will help you create maintainable, scalable, and professional Python projects that are easy to collaborate on and deploy.