CREATE TABLE casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(255) NOT NULL,
    v0 VARCHAR(255) DEFAULT NULL,
    v1 VARCHAR(255) DEFAULT NULL,
    v2 VARCHAR(255) DEFAULT NULL,
    v3 VARCHAR(255) DEFAULT NULL,
    v4 VARCHAR(255) DEFAULT NULL,
    v5 VARCHAR(255) DEFAULT NULL
);

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'HR', '/api/v1/employees', 'GET'),
('p', 'HR', '/api/v1/employees/:id', 'GET'),
('p', 'HR', '/api/v1/employees', 'POST'),
('p', 'HR', '/api/v1/employees/:id', 'PUT'),
('p', 'HR', '/api/v1/employees/:id', 'DELETE'),
('p', 'HR', '/api/v1/attendance/:id', 'GET'),
('p', 'HR', '/api/v1/leave/:id', 'GET'),
('p', 'Manager', '/api/v1/attendance/:id', 'GET'),
('p', 'Manager', '/api/v1/leave/:id', 'GET'),
('p', 'Employee', '/api/v1/auth/me', 'GET'),
('p', 'Employee', '/api/v1/attendance/:id', 'GET'),
('p', 'Employee', '/api/v1/attendance/:id/:date', 'GET');