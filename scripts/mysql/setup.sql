create table delivery_order (id varchar(36) PRIMARY KEY, distance MEDIUMINT unsigned,  status ENUM("UNASSIGNED", "TAKEN", "SUCCESS"));
