### 概述

一个针对Go的数据流分析器，支持

1. 基于GoAST的数据流分析引擎（yet another static analyzer / YASA）建设，YASA引擎负责
   引擎算法执行，以GoAST格式数据作为输入，经过核心分析算法模块，生成对应控制流，数据流依
   赖，并且在过程中，引入规则模块，进行具体规则检查。如果有规则匹配，则通过上报模块进行
   问题封装及上报。
2. 查询式语言（YAQL），查询式规则语言，通过声明式的规则描述，帮助业务和安全人员快速
   定制规则，结合统一语言表达方式及底层数据流引擎，帮助使用者通过程序分析底座快速查找和
   定位问题。支持基础的查询规则定义与解析。

### 设计思路

#### 1. **GoAST 数据流分析引擎（YASA）**

- **输入**：GoAST 格式的数据。
- **核心分析算法模块**：负责解析 GoAST，构建控制流图（CFG）和数据流图（DFG）。
- **规则模块**：包含一组规则，用于检测特定的代码模式或潜在问题。
- **上报模块**：当规则匹配时，生成报告并上报问题。

#### 2. **查询式语言（YAQL）**

- **规则定义**：使用声明式语法定义规则。
- **规则执行**：通过 YAQL 查询语言在数据流分析引擎上执行规则。

### 设计组件

1. **GoAST 解析器**
   - 解析 Go 代码，生成 AST。
   - 可以使用现有的解析库，如 `go/ast`。
2. **控制流图（CFG）生成器**
   
   - 从 AST 生成控制流图，表示程序的执行路径。
3. **数据流图（DFG）生成器**
   
   - 从 AST 生成数据流图，表示变量和数据在程序中的流动。
4. **规则模块**
   
   - 包含一组规则，每个规则使用 YAQL 定义。
   - 规则可以是安全检查、性能优化等。
5. **上报模块**
   
   - 当规则匹配时，生成报告并上报问题。
6. **YAQL 查询引擎**
   
   - 解析和执行 YAQL 查询。
   - 与数据流分析引擎集成，执行规则检查。

### 基础使用场景

假设我们有一个简单的 Go 代码片段，我们希望检测未初始化的变量使用问题。

```go
package main

import "fmt"

func main() {
    var x int
    fmt.Println(x)
}
```

#### 分析流程

1. **UAST 解析器**

使用 `go/ast` 解析 Go 代码并生成 UAST。

2. **控制流图（CFG）生成器**

从 UAST 生成控制流图，表示程序的执行路径。

3. **数据流图（DFG）生成器**

从 UAST 生成数据流图，表示变量和数据在程序中的流动。

4. **规则模块**

定义一个规则，检测未初始化的变量使用。

```yaql
RULE UninitializedVariable {
    MATCH {
        VariableDeclaration(var) AND
        VariableUsage(var) AND
        NOT(VariableInitialization(var))
    }
    REPORT {
        message: "Variable {var} is used without initialization."
    }
}
```

5. **上报模块**

当规则匹配时，生成报告并上报问题。

### 场景2：污点分析

污点分析（Taint Analysis）是一种静态分析技术，用于跟踪数据从源（source）到汇（sink）的传播路径，检测潜在的安全漏洞（如SQL注入、跨站脚本攻击等）。要使用你的工具编写查询语言（YAQL）来执行污点分析，你需要定义源、传播和汇的规则。

下面是一个详细的步骤和示例，展示如何编写 YAQL 规则来执行污点分析。

#### 污点分析步骤

1. **定义源（Source）**：确定哪些输入点可以引入污点数据。
2. **定义传播（Propagation）**：定义数据如何在程序中传播。
3. **定义汇（Sink）**：确定哪些输出点可能导致安全问题。
4. **编写规则**：定义源、传播和汇的规则。
5. **执行分析**：运行工具，检测从源到汇的污点传播路径。

#### 示例：检测 SQL 注入漏洞

假设我们有以下 Go 代码：

```go
package main

import (
    "database/sql"
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    userInput := r.URL.Query().Get("input")
    query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)
    db, _ := sql.Open("mysql", "user:password@/dbname")
    db.Query(query)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

#### 规则定义

1. **定义源（Source）**

源是用户输入，例如 HTTP 请求参数。

```yaql
RULE Source {
    MATCH {
        FunctionCall("r.URL.Query().Get", input)
    }
    TAG {
        taint: "source"
    }
}
```

2. **定义传播（Propagation）**

传播是数据在程序中的传递，例如赋值和函数调用。

```yaql
RULE Propagation {
    MATCH {
        Assignment(var1, var2) OR
        FunctionCall("fmt.Sprintf", var1, var2)
    }
    PROPAGATE {
        taint: "source"
    }
}
```

3. **定义汇（Sink）**

汇是潜在的安全问题点，例如 SQL 查询。

```yaql
RULE Sink {
    MATCH {
        FunctionCall("db.Query", query)
    }
    CHECK {
        taint: "source"
    }
    REPORT {
        message: "Potential SQL Injection vulnerability with tainted data in query: {query}"
    }
}
```

#### 规则解释

- **Source 规则**：匹配 HTTP 请求参数的获取，并标记为污点源。
- **Propagation 规则**：匹配变量赋值和 `fmt.Sprintf` 函数调用，传播污点标记。
- **Sink 规则**：匹配 SQL 查询的执行，检查是否包含污点数据，并报告潜在的 SQL 注入漏洞。

#### 执行分析

1. **解析 Go 代码**：使用 UAST 解析器生成 UAST。
2. **生成 CFG 和 DFG**：使用 CFG 和 DFG 生成器生成控制流图和数据流图。
3. **执行 YAQL 规则**：使用 YAQL 查询引擎执行定义的规则。
4. **上报问题**：当检测到污点数据从源到汇的传播路径时，生成并上报问题。



### 基础功能

#### 示例 1：检测未初始化的变量使用

我们希望检测未初始化的变量使用问题，如前面提到的示例。

```go
package main

import "fmt"

func main() {
    var x int
    fmt.Println(x)
}
```

规则：

```yaql
RULE UninitializedVariable {
    MATCH {
        VariableDeclaration(var) AND
        VariableUsage(var) AND
        NOT(VariableInitialization(var))
    }
    REPORT {
        message: "Variable {var} is used without initialization."
    }
}
```

#### 示例 2：检测未关闭的文件句柄

我们希望检测打开文件后未关闭的情况。

```go
package main

import "os"

func main() {
    file, err := os.Open("example.txt")
    if err != nil {
        return
    }
    // file.Close() is missing
}
```

规则：

```yaql
RULE UnclosedFileHandle {
    MATCH {
        FunctionCall("os.Open", args) AND
        NOT(FunctionCall("file.Close", _))
    }
    REPORT {
        message: "File opened but not closed."
    }
}
```

#### 示例 3：检测无效的类型断言

我们希望检测无效的类型断言。

```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"
    s := i.(int)  // invalid type assertion
    fmt.Println(s)
}
```

规则：

```yaql
RULE InvalidTypeAssertion {
    MATCH {
        TypeAssertion(variable, targetType) AND
        NOT(TypeCheck(variable, targetType))
    }
    REPORT {
        message: "Invalid type assertion: {variable} to {targetType}."
    }
}
```

#### 示例 4：检测死代码

我们希望检测死代码，即永远不会执行的代码。

```go
package main

func main() {
    return
    fmt.Println("This is dead code")
}
```

规则：

```yaql
RULE DeadCode {
    MATCH {
        UnreachableCode(code)
    }
    REPORT {
        message: "Unreachable code detected."
    }
}
```

#### 示例 5：检测未处理的错误

我们希望检测函数调用后未处理的错误。

```go
package main

import "os"

func main() {
    file, _ := os.Open("example.txt")
    // error is ignored
}
```

规则：

```yaql
RULE UnhandledError {
    MATCH {
        FunctionCallWithIgnoredError(funcName, args)
    }
    REPORT {
        message: "Error from function {funcName} is ignored."
    }
}
```

#### 示例 6：检测未使用的变量

我们希望检测声明但未使用的变量。

```go
package main

func main() {
    var unusedVar int
}
```

规则：

```yaql
RULE UnusedVariable {
    MATCH {
        VariableDeclaration(var) AND
        NOT(VariableUsage(var))
    }
    REPORT {
        message: "Variable {var} is declared but not used."
    }
}
```

#### 示例 7：检测潜在的SQL注入

我们希望检测潜在的 SQL 注入漏洞。

```go
package main

import "database/sql"

func main() {
    db, _ := sql.Open("mysql", "user:password@/dbname")
    userInput := "some input"
    db.Query("SELECT * FROM users WHERE name = '" + userInput + "'")
}
```

规则：

```yaql
RULE PotentialSQLInjection {
    MATCH {
        SQLQuery(query) AND
        ConcatenationInQuery(query, userInput)
    }
    REPORT {
        message: "Potential SQL injection detected in query: {query}."
    }
}
```

