
--
-- Data for Name: assets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.assets (id, name, abbrv, created_at, updated_at) FROM stdin;
54483425-efea-4472-84da-d50f50188e2b	Euro	EUR	1554540743	1554540743
fbc2bbc1-2d5a-4607-a1a4-1e26d284146c	US Dollar	USD	1554540743	1554540743
\.

--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (id, external_source_type, external_account_id, created_at, updated_at) FROM stdin;
5701249e-f33a-45a3-8722-e6917ccff6f0	bill.com	12345667	1543509959	1543509959
6eae6bb8-f7fb-425a-8af8-64adb616b739	bill.com	87654323	1543509959	1543509959
\.

--
-- Data for Name: entities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entities (id, entity_id, entity_type, account_id, created_at, updated_at) FROM stdin;
ca13775e-6134-477c-84f2-d0568ff66849	b582a786-48ec-469b-b655-17cf779b9ce1	project	5701249e-f33a-45a3-8722-e6917ccff6f0	1543509959	1543509959
3cf072df-8f92-4a18-863a-468c084bee73	1e04d5a3-a4bb-4bd8-b61a-bee83a4b57e2	project	6eae6bb8-f7fb-425a-8af8-64adb616b739	1543509959	1543509959
\.

--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, account_id, transaction_category, external_transaction_id, running_balance, created_at, updated_at) FROM stdin;
61b0c143-f1f9-457d-a889-80570b820348	5701249e-f33a-45a3-8722-e6917ccff6f0	random	a04c291f-234567	10200	1543509959	1543509959
fd54832d-d872-428b-a10d-17ddf782b4df	6eae6bb8-f7fb-425a-8af8-64adb616b739	random	a04c291f-234567	1000	1543509959	1543509959
\.

--
-- Data for Name: line_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.line_items (id, transaction_id, amount, asset_id, created_at) FROM stdin;
6b933db0-9316-4ddc-9399-b33ae9592a9f	61b0c143-f1f9-457d-a889-80570b820348	1000	fbc2bbc1-2d5a-4607-a1a4-1e26d284146c	1543509959
7939f305-400e-4e21-93ec-da3bf519f09e	61b0c143-f1f9-457d-a889-80570b820348	8200	fbc2bbc1-2d5a-4607-a1a4-1e26d284146c	1543509960
39cf8486-bde3-413e-aea4-188589b12d18	fd54832d-d872-428b-a10d-17ddf782b4df	-1000	54483425-efea-4472-84da-d50f50188e2b	1543509962
7bb133b8-15ef-4cfa-94cf-1c413b7c5cc1	fd54832d-d872-428b-a10d-17ddf782b4df	2000	54483425-efea-4472-84da-d50f50188e2b	1543509969
8797ac8a-bf4a-463c-a499-ce4a5643f7c5	fd54832d-d872-428b-a10d-17ddf782b4df	-1500	54483425-efea-4472-84da-d50f50188e2b	1543509969
\.
