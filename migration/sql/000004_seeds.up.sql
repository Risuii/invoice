-- Insert statements for Customers
insert into customers (customer_id, name, address) values ('0a7fb210-a232-4233-83bb-5af9cb37a0fe', 'discovery design', 'tangerang');
insert into customers (customer_id, name, address) values ('f822f341-2c63-4984-973c-b4b1e0b6739c', 'barrington publisher', 'jakarta');
insert into customers (customer_id, name, address) values ('1bb6c946-9ab5-43d8-828c-45b8ada531d1', 'sinar terang', 'bandung');
insert into customers (customer_id, name, address) values ('4c3b319f-6851-4aa6-b9ef-19765bed3ae4', 'gelap redup', 'padang');
insert into customers (customer_id, name, address) values ('84accf73-6296-4502-bcd0-b18098276e11', 'jaya maju', 'sabang');


-- Insert statements for Invoices
insert into Invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) values ('0001', '01-24-2023', 'service payment', 3, '0a7fb210-a232-4233-83bb-5af9cb37a0fe', '01-24-2024', 'Paid', 300, 400, 500);
insert into Invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) values ('0002', '02-25-2023', 'service payment', 3, 'f822f341-2c63-4984-973c-b4b1e0b6739c', '02-25-2024', 'Unpaid', 100, 200, 300);
insert into Invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) values ('0003', '03-26-2023', 'service payment', 2, '1bb6c946-9ab5-43d8-828c-45b8ada531d1', '03-26-2024', 'Paid', 600, 700, 800);
insert into Invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) values ('0004', '04-27-2023', 'service payment', 2, '4c3b319f-6851-4aa6-b9ef-19765bed3ae4', '04-27-2024', 'Unpaid', 700, 800, 900);
insert into Invoices (invoice_id, issue_date, subject, total_items, customer_id, due_date, status, sub_total, tax, grand_total) values ('0005', '05-26-2023', 'service payment', 2, '84accf73-6296-4502-bcd0-b18098276e11', '05-28-2024', 'Paid', 100, 300, 500);

-- Insert statements for Items
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0001', 'ca0f91e7-bc88-4149-9271-8f75d85f8cd0', 'design', 'service', 2, 300, 600);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0001', '2862126f-d7e2-4ae9-b006-bc62b52ea40a', 'development', 'service', 5, 100, 500);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0001', 'c84078a9-cfd8-4706-86ca-bb6a323c590f', 'meetings', 'service', 3, 200, 600);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0002', 'ffbb2961-66b0-4b1d-81e9-59e16bba3dd7', 'printer', 'hardware', 4, 200, 800);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0002', 'c7020df5-47dd-4744-bb7b-0b008e6ecd65', 'monitor', 'hardware', 5, 200, 1000);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0004', 'd2eb3301-ca71-490e-80d1-98f09786c2b1', 'design', 'service', 3, 300, 900);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0004', '8f04b997-fcfc-4753-86d6-8e8eace7579f', 'development', 'service', 6, 100, 600);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0005', '97d6b57f-de35-4cd0-974d-bea119cf1cdc', 'meetings', 'service', 4, 200, 800);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0005', '3bf16061-19f7-42e6-9299-ff2c1fa1d82b', 'development', 'service', 7, 100, 700);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0003', 'dc74af5e-adfd-4495-92c8-c95fc8d84923', 'printer', 'hardware', 6, 200, 1200);
insert into items (invoice_id, item_id, name, type, quantity, unit_price, amount) values ('0003', '0f0cbbab-b885-4e5c-be47-6f7504c2b8fb', 'monitor', 'service', 7, 200, 1400);
