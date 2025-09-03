CREATE TABLE vehicle_tab (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id VARCHAR(255),
    max_shipments INT,
    min_shipments INT
);

CREATE TABLE cost_per_shipment_tab (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id VARCHAR(255),
    zone_id VARCHAR(255),
    cost DOUBLE
);