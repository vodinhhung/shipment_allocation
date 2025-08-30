INSERT INTO vehicle_tab (vehicle_id, max_shipments, min_shipments)
VALUES ("V1", 2300, 0);

INSERT INTO vehicle_tab (vehicle_id, max_shipments, min_shipments)
VALUES ("V2", 1500, 500);

INSERT INTO cost_per_shipment_tab (vehicle_id, zone_id, cost)
VALUES ("V1", "Z1", 5);

INSERT INTO cost_per_shipment_tab (vehicle_id, zone_id, cost)
VALUES ("V1", "Z2", 7);

INSERT INTO cost_per_shipment_tab (vehicle_id, zone_id, cost)
VALUES ("V2", "Z2", 8);