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
VALUES ('apnic'), ('arin'), ('iana'), ('lacnic'), ('ripencc');

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
VALUES ('allocated'), ('assigned');

CREATE TABLE ip_ranges_a (
    id SERIAL PRIMARY KEY,
    rir INT NOT NULL REFERENCES rirs(id) ON DELETE RESTRICT,
    country_code CHAR(2),
    version_ip INT NOT NULL REFERENCES ip_versions(id) ON DELETE RESTRICT,
    start_ip INET NOT NULL,
    end_ip INET NOT NULL,
    status INT NOT NULL REFERENCES ip_range_statuses(id) ON DELETE RESTRICT,
    created_at DATE NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(start_ip, end_ip)
);
CREATE INDEX idx_ip_ranges_a_rir ON ip_ranges_a (rir);
CREATE INDEX idx_ip_ranges_a_version_ip on ip_ranges_a (version_ip);
CREATE INDEX idx_ip_ranges_a_start_ip ON ip_ranges_a (start_ip);
CREATE INDEX idx_ip_ranges_a_end_ip ON ip_ranges_a (end_ip);
CREATE INDEX idx_ip_ranges_a_status ON ip_ranges_a (status);

CREATE OR REPLACE VIEW ip_ranges AS
SELECT * FROM ip_ranges_a;