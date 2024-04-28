BEGIN;

CREATE TYPE status_type AS ENUM ('Paid', 'Unpaid');

CREATE TABLE public.Invoices (
    id bigint NOT NULL,
    invoice_id VARCHAR(10) NOT NULL UNIQUE,
    issue_date timestamp with time zone,
    subject character varying(255),
    total_items INT,
    customer_id UUID NOT NULL,
    due_date timestamp with time zone,
    status status_type,
    sub_total numeric,
    tax numeric,
    grand_total numeric,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);

CREATE SEQUENCE public.invoices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.invoices_id_seq OWNED BY public.invoices.id;

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT customer_id FOREIGN KEY (customer_id) REFERENCES public.customers(customer_id);

COMMIT;