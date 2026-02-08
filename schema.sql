CREATE TABLE services (
    service_name VARCHAR(50) PRIMARY KEY,
    display_name VARCHAR(100) NOT NULL,
    base_url VARCHAR(255) NOT NULL
);

CREATE TABLE credentials (
    id INT AUTO_INCREMENT PRIMARY KEY,
    service_name VARCHAR(50) REFERENCES services(service_name),
    user_id INT NOT NULL,
    
    -- OAuth credentials
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP,
);

CREATE TABLE workflows (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    name VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT FALSE,
    user_id INT NOT NULL
);

CREATE TABLE workflow_nodes (
    -- React flow stores ids as strings
    id VARCHAR(255) PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    workflow_id INT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    
    service_name VARCHAR(50) NOT NULL,
    task_name VARCHAR(50) NOT NULL,

    -- listener, action, transfomer
    type VARCHAR(50) NOT NULL,

    config JSONB DEFAULT '{}',
    credential_id INT REFERENCES credentials(id),
    position JSON
);

CREATE TABLE workflow_edges (
    id VARCHAR(255) AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    node_from VARCHAR(255) REFERENCES workflow_nodes(id),
    node_to VARCHAR(255) REFERENCES workflow_nodes(id),
    workflow_id INT NOT NULL
);

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    username VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
);
