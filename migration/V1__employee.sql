CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name, description) VALUES
('HR', 'Human Resources personnel with full access to employee management'),
('Manager', 'Department managers with access to team management'),
('Employee', 'Regular employees with basic access');



CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL UNIQUE,  -- 員工編號
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    department VARCHAR(50),
    position VARCHAR(50),
    hire_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE',  -- ACTIVE, INACTIVE, TERMINATED
    manager_id INT REFERENCES employees(id),  -- 主管ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- 插入員工數據
-- HR (密碼: hr123456)
INSERT INTO employees (
    employee_id,
    first_name,
    last_name,
    email,
    password_hash,  -- bcrypt hash of "hr123456"
    phone,
    department,
    position,
    hire_date,
    status
) VALUES (
    'HR001',
    'John',
    'Doe',
    'hr@example.com',
    '$2a$12$vTaqOtTAvADcpMCnYhSuouGKDyc.z4OA977Bop/AYP/Rwet8RbZNi',
    '1234567890',
    'Human Resources',
    'HR Manager',
    '2023-01-01',
    'ACTIVE'
);

-- 部門經理 (密碼: manager123)
INSERT INTO employees (
    employee_id,
    first_name,
    last_name,
    email,
    password_hash,  -- bcrypt hash of "manager123"
    phone,
    department,
    position,
    hire_date,
    status
) VALUES (
    'MGR001',
    'Jane',
    'Smith',
    'manager@example.com',
    '$2a$12$/sIJgDlLaZrEAA8w8rppY.e.uLLxyt6UgW8EiJXnqWO2g2T8Gg4w6',
    '0987654321',
    'Engineering',
    'Engineering Manager',
    '2023-01-02',
    'ACTIVE'
);

-- 普通員工 (密碼: emp123456)
INSERT INTO employees (
    employee_id,
    first_name,
    last_name,
    email,
    password_hash,  -- bcrypt hash of "emp123456"
    phone,
    department,
    position,
    hire_date,
    status,
    manager_id
) VALUES (
    'EMP001',
    'Bob',
    'Wilson',
    'employee@example.com',
    '$2a$12$uZC52Q0MrU7l3OKtEDUACOj8uNYrdEGhagsSnSAEWWg7ZeDoow4TS',
    '1122334455',
    'Engineering',
    'Software Engineer',
    '2023-01-03',
    'ACTIVE',
    (SELECT id FROM employees WHERE employee_id = 'MGR001')
);


CREATE TABLE employee_roles (
    employee_id INT REFERENCES employees(id),
    role_id INT REFERENCES roles(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (employee_id, role_id)
);

INSERT INTO employee_roles (employee_id, role_id)
SELECT 
    (SELECT id FROM employees WHERE employee_id = 'HR001'),
    (SELECT id FROM roles WHERE name = 'HR');

-- 經理關聯
INSERT INTO employee_roles (employee_id, role_id)
SELECT 
    (SELECT id FROM employees WHERE employee_id = 'MGR001'),
    (SELECT id FROM roles WHERE name = 'Manager');

-- 普通員工關聯
INSERT INTO employee_roles (employee_id, role_id)
SELECT 
    (SELECT id FROM employees WHERE employee_id = 'EMP001'),
    (SELECT id FROM roles WHERE name = 'Employee');


-- 為了支持角色繼承，在casbin_rule表中添加角色繼承規則
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
('g', 'HR001', 'HR'),          -- HR用戶具有HR角色
('g', 'MGR001', 'Manager'),    -- 經理用戶具有Manager角色
('g', 'EMP001', 'Employee');   -- 普通員工具有Employee角色
