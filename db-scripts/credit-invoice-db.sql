use credit_invoice;

create table customer_credentials (
    customer_id varchar(255),
    credit_account_id int
);

insert into customer_credentials
values ("abc-123-def", 123);