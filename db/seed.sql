
--
-- Data for Name: entities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entities (id, entity_id, entity_type, created_at) FROM stdin;
ca13775e-6134-477c-84f2-d0568ff66849	b582a786-48ec-469b-b655-17cf779b9ce1	project	1543509959
3cf072df-8f92-4a18-863a-468c084bee73	1e04d5a3-a4bb-4bd8-b61a-bee83a4b57e2	project	1543509959
\.

--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (id, external_source_type, entity_id, external_account_id, created_at) FROM stdin;
5701249e-f33a-45a3-8722-e6917ccff6f0	bill.com	ca13775e-6134-477c-84f2-d0568ff66849	12345667	1543509959
6eae6bb8-f7fb-425a-8af8-64adb616b739	bill.com	3cf072df-8f92-4a18-863a-468c084bee73	87654323	1543509959
\.

--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, account_id, transaction_category, external_transaction_id, running_balance, created_at) FROM stdin;
fd54832d-d872-428b-a10d-17ddf782b4df	6eae6bb8-f7fb-425a-8af8-64adb616b739	random	a04c291f-234567	0	1543509959
\.

--
-- Data for Name: line_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.line_items (id, transaction_id, amount, created_at) FROM stdin;
39cf8486-bde3-413e-aea4-188589b12d18	fd54832d-d872-428b-a10d-17ddf782b4df	-1000	1543509962
7bb133b8-15ef-4cfa-94cf-1c413b7c5cc1	fd54832d-d872-428b-a10d-17ddf782b4df	2000	1543509969
8797ac8a-bf4a-463c-a499-ce4a5643f7c5	fd54832d-d872-428b-a10d-17ddf782b4df	-1500	1543509969
\.
