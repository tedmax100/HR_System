CREATE TYPE attendance_event_type AS ENUM ('CHECK_IN', 'CHECK_OUT');

CREATE TABLE attendance_events (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL REFERENCES employees(employee_id),
    event_type attendance_event_type NOT NULL,
    event_time TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(100),
    device_info VARCHAR(255),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_attendance_employee_time ON attendance_events(employee_id, event_time);
