CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE(name)
);

CREATE TABLE rirs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
INSERT INTO rirs (name) 
VALUES ('apnic'), ('arin'), ('afrinic'), ('lacnic'), ('ripencc');

CREATE TABLE ip_versions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
INSERT INTO ip_versions (name)
VALUES ('ipv4'), ('ipv6');

CREATE TABLE ip_range_statuses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
INSERT INTO ip_range_statuses (name)
VALUES ('allocated'), ('assigned'), ('available'), ('reserved'), ('unknown');

CREATE TABLE apnic (
    id SERIAL PRIMARY KEY,
    country_code CHAR(2),
    ip_version_id INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    quantity INT NOT NULL,
    status_id INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    status_changed_at DATE,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_apnic_ip_version_id on apnic (ip_version_id);
CREATE INDEX idx_apnic_start_ip ON apnic (start_ip);
CREATE INDEX idx_apnic_end_ip ON apnic (end_ip);
CREATE INDEX idx_apnic_status_id ON apnic (status_id);

CREATE TABLE arin (
    id SERIAL PRIMARY KEY,
    country_code CHAR(2),
    ip_version_id INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    quantity INT NOT NULL,
    status_id INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    status_changed_at DATE,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_arin_ip_version_id on arin (ip_version_id);
CREATE INDEX idx_arin_start_ip ON arin (start_ip);
CREATE INDEX idx_arin_end_ip ON arin (end_ip);
CREATE INDEX idx_arin_status_id ON arin (status_id);

CREATE TABLE afrinic (
    id SERIAL PRIMARY KEY,
    country_code CHAR(2),
    ip_version_id INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    quantity INT NOT NULL,
    status_id INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    status_changed_at DATE,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_afrinic_ip_version_id on afrinic (ip_version_id);
CREATE INDEX idx_afrinic_start_ip ON afrinic (start_ip);
CREATE INDEX idx_afrinic_end_ip ON afrinic (end_ip);
CREATE INDEX idx_afrinic_status_id ON afrinic (status_id);

CREATE TABLE lacnic (
    id SERIAL PRIMARY KEY,
    country_code CHAR(2),
    ip_version_id INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    quantity INT NOT NULL,
    status_id INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    status_changed_at DATE,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_lacnic_ip_version_id on lacnic (ip_version_id);
CREATE INDEX idx_lacnic_start_ip ON lacnic (start_ip);
CREATE INDEX idx_lacnic_end_ip ON lacnic (end_ip);
CREATE INDEX idx_lacnic_status_id ON lacnic (status_id);

CREATE TABLE ripencc (
    id SERIAL PRIMARY KEY,
    country_code CHAR(2),
    ip_version_id INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    quantity INT NOT NULL,
    status_id INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    status_changed_at DATE,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_ripencc_ip_version_id on ripencc (ip_version_id);
CREATE INDEX idx_ripencc_start_ip ON ripencc (start_ip);
CREATE INDEX idx_ripencc_end_ip ON ripencc (end_ip);
CREATE INDEX idx_ripencc_status_id ON ripencc (status_id);

CREATE MATERIALIZED VIEW IF NOT EXISTS ip_ranges AS
    SELECT 'apnic_' || apnic.id as id, rirs.name as rir_name, apnic.country_code, ip_versions.name as ip_version_name, apnic.start_ip, apnic.end_ip, apnic.quantity, ip_range_statuses.name as status_name, apnic.status_changed_at 
    FROM apnic
        JOIN rirs ON rirs.id = 1
        JOIN ip_versions ON ip_versions.id = apnic.ip_version_id
        JOIN ip_range_statuses ON apnic.status_id = ip_range_statuses.id
UNION ALL
    SELECT 'arin_' || arin.id as id, rirs.name as rir_name, arin.country_code, ip_versions.name as ip_version_name, arin.start_ip, arin.end_ip, arin.quantity, ip_range_statuses.name as status_name, arin.status_changed_at
    FROM arin 
        JOIN rirs ON rirs.id = 2
        JOIN ip_versions ON ip_versions.id = arin.ip_version_id
        JOIN ip_range_statuses ON arin.status_id = ip_range_statuses.id
UNION ALL
    SELECT 'afrinic_' || afrinic.id as id, rirs.name as rir_name, afrinic.country_code, ip_versions.name as ip_version_name, afrinic.start_ip, afrinic.end_ip, afrinic.quantity, ip_range_statuses.name as status_name, afrinic.status_changed_at
    FROM afrinic 
        JOIN rirs ON rirs.id = 3
        JOIN ip_versions ON ip_versions.id = afrinic.ip_version_id
        JOIN ip_range_statuses ON afrinic.status_id = ip_range_statuses.id
UNION ALL
    SELECT 'lacnic_' || lacnic.id as id, rirs.name as rir_name, lacnic.country_code, ip_versions.name as ip_version_name, lacnic.start_ip, lacnic.end_ip, lacnic.quantity, ip_range_statuses.name as status_name, lacnic.status_changed_at
    FROM lacnic
        JOIN rirs ON rirs.id = 4
        JOIN ip_versions ON ip_versions.id = lacnic.ip_version_id
        JOIN ip_range_statuses ON lacnic.status_id = ip_range_statuses.id
UNION ALL
    SELECT 'ripencc_' || ripencc.id as id, rirs.name as rir_name, ripencc.country_code, ip_versions.name as ip_version_name, ripencc.start_ip, ripencc.end_ip, ripencc.quantity, ip_range_statuses.name as status_name, ripencc.status_changed_at
    FROM ripencc 
        JOIN rirs ON rirs.id = 5
        JOIN ip_versions ON ip_versions.id = ripencc.ip_version_id
        JOIN ip_range_statuses ON ripencc.status_id = ip_range_statuses.id;