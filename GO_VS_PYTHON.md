# Go vs Python: Key Differences for This Project

## Type System

### Python (Dynamic Typing)
```python
config = {}  # Type inferred at runtime
config["api_secret"] = "secret"
```

### Go (Static Typing)
```go
type Config struct {
    APISecret string `json:"api_secret"`
}
var config Config
config.APISecret = "secret"
```

**Key Difference**: Go requires you to define types upfront. This catches errors at compile time instead of runtime.

---

## Error Handling

### Python (Exceptions)
```python
try:
    data = read_file(path)
    config = json.loads(data)
except Exception as e:
    print(f"Error: {e}")
```

### Go (Explicit Returns)
```go
data, err := ioutil.ReadFile(path)
if err != nil {
    return fmt.Errorf("failed to read: %w", err)
}

var config Config
if err := json.Unmarshal(data, &config); err != nil {
    return fmt.Errorf("failed to parse: %w", err)
}
```

**Key Difference**: Go doesn't have exceptions. Functions return errors that you must explicitly check. This is more verbose but makes error paths obvious.

---

## Pointers

### Python
```python
def modify_config(config):
    config["value"] = 123
    # Changes the original dict
```

### Go
```go
func modifyConfig(config *Config) {
    config.Value = 123
    // * means pointer - changes original
}

func readConfig(config Config) {
    config.Value = 123
    // No * - this is a copy, original unchanged
}
```

**Key Difference**: Go makes the difference between pass-by-value and pass-by-reference explicit using `*`. This helps you understand when you're copying data vs. modifying the original.

---

## Interfaces

### Python (Explicit)
```python
from abc import ABC, abstractmethod

class Workflow(ABC):
    @abstractmethod
    def execute(self):
        pass

class MyWorkflow(Workflow):  # Must explicitly inherit
    def execute(self):
        print("executing")
```

### Go (Implicit)
```go
type Workflow interface {
    Execute(config *Config, args []string) error
}

type MyWorkflow struct{}

func (w *MyWorkflow) Execute(config *Config, args []string) error {
    fmt.Println("executing")
    return nil
}
// Automatically satisfies Workflow interface!
```

**Key Difference**: Go interfaces are satisfied implicitly. If a type has the right methods, it implements the interface automatically. This is called "duck typing at compile time."

---

## Package/Module System

### Python
```python
# mymodule.py
def my_function():
    pass

# main.py
import mymodule
mymodule.my_function()
```

### Go
```go
// mymodule.go
package main  // All files in same directory share package

func MyFunction() {  // Capital = exported/public
    privateFunction()
}

func privateFunction() {  // lowercase = private
}

// main.go
package main
// No import needed - same package!
MyFunction()
```

**Key Difference**: 
1. All `.go` files in a directory are in the same package
2. Capitalization determines visibility (Capital = public, lowercase = private)
3. No `__init__.py` needed

---

## Method Receivers

### Python
```python
class MyClass:
    def my_method(self, arg):
        self.value = arg
```

### Go
```go
type MyType struct {
    value string
}

func (m *MyType) MyMethod(arg string) {
    m.value = arg  // * allows modification
}

func (m MyType) ReadOnly() string {
    // No * - receives a copy
    return m.value
}
```

**Key Difference**: Go doesn't have classes. Instead, you define types and attach methods using receivers. Use `*` for methods that modify, no `*` for read-only methods.

---

## Slices vs Lists

### Python
```python
items = [1, 2, 3]
items.append(4)
subset = items[1:3]  # [2, 3]
for item in items:
    print(item)
```

### Go
```go
items := []int{1, 2, 3}
items = append(items, 4)  // Returns new slice
subset := items[1:3]      // [2, 3]
for _, item := range items {
    fmt.Println(item)
}
```

**Key Difference**: Go slices are similar to Python lists but:
1. `append()` returns a new slice (may or may not reuse memory)
2. Slices are typed (`[]int`, `[]string`, etc.)
3. `range` gives you index and value (use `_` to ignore index)

---

## Maps vs Dictionaries

### Python
```python
data = {"key": "value"}
data["new_key"] = "new_value"

if "key" in data:
    print(data["key"])

for key, value in data.items():
    print(key, value)
```

### Go
```go
data := make(map[string]string)
data["key"] = "value"
data["new_key"] = "new_value"

// "comma ok" idiom
if value, ok := data["key"]; ok {
    fmt.Println(value)
}

for key, value := range data {
    fmt.Println(key, value)
}
```

**Key Difference**: Go maps use the "comma ok" idiom to check existence. Missing keys return the zero value (not nil/None) and `ok` is false.

---

## JSON Handling

### Python
```python
import json

data = {"name": "Alice", "age": 30}
json_str = json.dumps(data)
parsed = json.loads(json_str)
```

### Go
```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

person := Person{Name: "Alice", Age: 30}
jsonBytes, err := json.Marshal(person)
// jsonBytes is []byte, not string

var parsed Person
err = json.Unmarshal(jsonBytes, &parsed)
```

**Key Difference**: 
1. Go requires struct tags to map JSON fields
2. Returns bytes, not strings
3. Must pass pointer to Unmarshal for it to modify struct
4. Must check errors explicitly

---

## Goroutines vs Threading

### Python
```python
import threading

def task():
    print("running")

thread = threading.Thread(target=task)
thread.start()
thread.join()
```

### Go
```go
func task() {
    fmt.Println("running")
}

go task()  // Launch goroutine
// Use channels to wait/communicate
```

**Key Difference**: Goroutines are much lighter weight than threads. The `go` keyword is all you need to run something concurrently. Use channels for communication between goroutines.

---

## Variable Declaration

### Python
```python
name = "Alice"
age = 30
items = []
```

### Go
```go
// Multiple ways to declare:
var name string = "Alice"  // Explicit type
var age = 30               // Type inferred
items := []string{}        // Short declaration (inside functions)

// Multiple variables:
var (
    host   = "localhost"
    port   = 8080
    active = true
)
```

**Key Difference**: Go has several declaration styles. `:=` is the most common inside functions. You must use `var` at package level.

---

## nil vs None

### Python
```python
value = None
if value is None:
    print("empty")
```

### Go
```go
var value *string  // Pointer, defaults to nil
if value == nil {
    fmt.Println("empty")
}
```

**Key Difference**: 
- Go's `nil` is similar to Python's `None`
- Only pointers, slices, maps, channels, interfaces, and functions can be `nil`
- Zero values differ by type: `0` for numbers, `""` for strings, `false` for bools

---

## defer Statement

### Python
```python
try:
    f = open("file.txt")
    # do stuff
finally:
    f.close()
```

### Go
```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()  // Runs when function returns

// Do stuff with f
// defer ensures Close() happens even if we return early
```

**Key Difference**: `defer` schedules a function call to run when the surrounding function returns. It's more flexible than try/finally.

---

## Compilation

### Python
```bash
python main.py  # Interpreted, runs immediately
```

### Go
```bash
go run main.go              # Compile and run (quick for development)
go build -o myapp          # Compile to binary
./myapp                     # Run the binary

# Cross-compile
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
```

**Key Difference**: Go compiles to a single binary with no dependencies. This makes deployment much simpler than Python (no need for virtual environments, dependencies, or even the Go toolchain on the target system).

---

## Quick Reference Table

| Feature | Python | Go |
|---------|--------|-----|
| **Error handling** | Exceptions | Explicit error returns |
| **Types** | Dynamic | Static |
| **Classes** | Yes | No (use structs + methods) |
| **Interfaces** | Explicit inheritance | Implicit satisfaction |
| **Concurrency** | Threading, asyncio | Goroutines, channels |
| **Package visibility** | `_` prefix | Capitalization |
| **Memory management** | Garbage collection | Garbage collection |
| **Compilation** | Interpreted | Compiled |
| **null/nil** | None | nil (typed) |
| **Multiple return** | Tuple unpacking | Built-in |

