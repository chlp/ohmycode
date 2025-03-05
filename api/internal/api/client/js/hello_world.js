const helloWorldButton = document.getElementById('hello-world-button');
helloWorldButton.onclick = () => {
    const helloWorldType = languages[app.lang].helloWorld;
    if (typeof helloWorldType === 'undefined') {
        return;
    }
    if (file.writer_id !== '' && file.writer_id !== app.id) {
        return;
    }

    let helloWorldCode = '';
    switch (helloWorldType) {
        case 'go':
            helloWorldCode = `package main

import "fmt"

func main() {
    fmt.Println("Hello, OhMyCode!")
}`;
            break;
        case 'java':
            helloWorldCode = `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, OhMyCode!");
    }
}`;
            break;
        case 'json':
            helloWorldCode = `{"name":"Aleksei","age":30,"profession":"Developer","skills":["Python","Docker","AWS"]}`;
            break;
        case 'markdown':
            helloWorldCode = `# Project Documentation

## Introduction

This document describes the project structure, main objectives, and the technologies used. It includes code examples, lists, and tables for better information presentation.

---

## Project Goals

The main goals of the project:

- Create a high-performance system.
- Ensure scalability.
- Guarantee data security.

### Additional Goals

- Optimize the development process.
- Increase unit test coverage to 95%.

---

## Tasks

### Short-term Tasks

1. Write the project specification.
2. Create a prototype.
3. Test the basic functionality.

### Long-term Tasks

- Optimize algorithms.
- Integrate with external APIs.
- Support multi-user mode.

---

## Technologies

The project uses the following technologies:

- **Programming Languages**:
  - Python
  - JavaScript (Node.js)
- **Databases**:
  - PostgreSQL
  - Redis
- **Services**:
  - AWS Lambda
  - Docker

---

| No. | Name | Age | Profession |
| --- | --- | --- | --- |
| 1 | Ivan Ivanov | 29 | Developer |
| 2 | Peter Peterson | 34 | Designer |
| 3 | Black White | 25 | Tester |
`;
            break;
        case 'mysql':
            helloWorldCode = `create table employees
(
    id            bigint       not null primary key,
    first_name       varchar(255) null,
    last_name       varchar(255) null,
    salary int       null,
    department_id int       null
);

create table departments
(
    id         bigint       not null primary key,
    name       varchar(255) null
);

#SELECT * FROM

INSERT INTO employees (id, first_name, last_name, salary, department_id) values (1, 'Aa', 'FamilyA', 1000, 1);
INSERT INTO employees (id, first_name, last_name, salary, department_id) values (2, 'Ab', 'FamilyA', 100, 1);
INSERT INTO employees (id, first_name, last_name, salary, department_id) values (3, 'Cc', 'Hello, OhMyCode!', 3000, 3);

INSERT INTO departments (id, name) values (1, 'D_A');
INSERT INTO departments (id, name) values (2, 'D_B');
INSERT INTO departments (id, name) values (3, 'D_C');

EXPLAIN SELECT departments.* FROM departments
LEFT JOIN employees ON departments.id = employees.department_id
WHERE employees.id IS NOT NULL;

select last_name from employees where first_name like 'C%';`;
            break;
        case 'php':
            helloWorldCode = `<?php

echo 'Hello, OhMyCode!';`;
            break;
        case 'postgres':
            helloWorldCode = `create table employees
(
    id            bigint       not null primary key,
    first_name       varchar(255) null,
    last_name       varchar(255) null,
    salary int       null,
    department_id int       null
);

create table departments
(
    id         bigint       not null primary key,
    name       varchar(255) null
);

INSERT INTO employees (id, first_name, last_name, salary, department_id) values (1, 'Aa', 'FamilyA', 1000, 1);
INSERT INTO employees (id, first_name, last_name, salary, department_id) values (2, 'Ab', 'FamilyA', 100, 1);
INSERT INTO employees (id, first_name, last_name, salary, department_id) values (3, 'Cc', 'Hello, OhMyCode!', 3000, 3);

INSERT INTO departments (id, name) values (1, 'D_A');
INSERT INTO departments (id, name) values (2, 'D_B');
INSERT INTO departments (id, name) values (3, 'D_C');

EXPLAIN SELECT departments.* FROM departments
LEFT JOIN employees ON departments.id = employees.department_id
WHERE employees.id IS NOT NULL;

select last_name from employees where first_name like 'C%';`;
            break;
        default:
            console.error('no helloWorld code');
    }

    if (helloWorldCode !== '') {
        const content = helloWorldCode + "\r\n" + contentCodeMirror.getValue();
        contentMarkdownBlock.innerHTML = marked.parse(content);
        contentCodeMirror.setValue(content);
    }
}
