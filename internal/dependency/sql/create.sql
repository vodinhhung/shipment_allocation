CREATE TABLE zone_tab (
    id INT AUTO_INCREMENT PRIMARY KEY,
    no_shipments INT
);

CREATE TABLE vehicle_tab (
    id INT AUTO_INCREMENT PRIMARY KEY,
    max_shipments INT,
    min_shipments INT
);

CREATE TABLE cost_per_shipment_tab (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id INT,
    zone_id INT,
    cost INT
);