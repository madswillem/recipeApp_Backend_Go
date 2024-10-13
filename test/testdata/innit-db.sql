--
-- PostgreSQL database dump
--

-- Dumped from database version 15.6
-- Dumped by pg_dump version 15.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: diet; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.diet (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    name text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.diet OWNER TO mads;

--
-- Name: ingredient; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.ingredient (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    name text NOT NULL,
    standard_unit text,
    ndb_number bigint,
    category text,
    fdic_id bigint
);


ALTER TABLE public.ingredient OWNER TO mads;

--
-- Name: nutritional_value; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.nutritional_value (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    ingredient_id uuid,
    recipe_id uuid,
    kcal numeric,
    kj numeric,
    fat numeric,
    saturated_fat numeric,
    carbohydrate numeric,
    sugar numeric,
    protein numeric,
    salt numeric,
    nutriscore character(1),
    CONSTRAINT nutritional_value_check CHECK (((((ingredient_id IS NOT NULL))::integer + ((recipe_id IS NOT NULL))::integer) = 1))
);


ALTER TABLE public.nutritional_value OWNER TO mads;

--
-- Name: rating; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.rating (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    ingredient_id uuid,
    recipe_id uuid,
    overall numeric,
    mon numeric,
    tue numeric,
    wed numeric,
    thu numeric,
    fri numeric,
    sat numeric,
    sun numeric,
    win numeric,
    spr numeric,
    sum numeric,
    aut numeric,
    thirtydegree numeric,
    twentiedegree numeric,
    tendegree numeric,
    zerodegree numeric,
    subzerodegree numeric,
    CONSTRAINT rating_check CHECK (((((ingredient_id IS NOT NULL))::integer + ((recipe_id IS NOT NULL))::integer) = 1))
);


ALTER TABLE public.rating OWNER TO mads;

--
-- Name: recipe_ingredient; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.recipe_ingredient (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    recipe_id uuid,
    ingredient_id uuid,
    amount bigint,
    unit text
);


ALTER TABLE public.recipe_ingredient OWNER TO mads;

--
-- Name: recipe_selects_views_log; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.recipe_selects_views_log (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    day timestamp without time zone DEFAULT (now())::timestamp without time zone,
    selects bigint,
    views bigint,
    recipe_id uuid NOT NULL,
    view_change bigint,
    selects_change bigint
);


ALTER TABLE public.recipe_selects_views_log OWNER TO mads;

--
-- Name: recipes; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.recipes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    author uuid NOT NULL,
    name text NOT NULL,
    cuisine text,
    yield smallint,
    yield_unit text,
    prep_time interval,
    cooking_time interval,
    version bigint DEFAULT 0,
    selects bigint DEFAULT 0,
    views bigint DEFAULT 0
);


ALTER TABLE public.recipes OWNER TO mads;

--
-- Name: rel_diet_recipe; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.rel_diet_recipe (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    recipe_id uuid,
    diet_id uuid
);


ALTER TABLE public.rel_diet_recipe OWNER TO mads;

--
-- Name: rel_diet_user; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.rel_diet_user (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    diet_id uuid NOT NULL
);


ALTER TABLE public.rel_diet_user OWNER TO mads;

--
-- Name: step; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.step (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    step text,
    recipe_id uuid NOT NULL,
    technique_id uuid,
    ingredient_id uuid
);


ALTER TABLE public.step OWNER TO mads;

--
-- Name: technique; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public.technique (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    description text NOT NULL
);


ALTER TABLE public.technique OWNER TO mads;

--
-- Name: user; Type: TABLE; Schema: public; Owner: mads
--

CREATE TABLE public."user" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
    cookie text,
    ip text,
    groups jsonb
);


ALTER TABLE public."user" OWNER TO mads;

--
-- Data for Name: diet; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.diet (id, created_at, name, description) FROM stdin;
bbadd945-5557-459f-951e-9ad3ad277059	2024-08-26 20:38:25.856765	Vegetarien	A diet woithout fish and meat
\.


--
-- Data for Name: ingredient; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.ingredient (id, created_at, name, standard_unit, ndb_number, category, fdic_id) FROM stdin;
84eb6da1-25b9-40ec-97a1-c0db1844ca54	2024-06-15 23:34:15.856578	tomato	pice	100261	Vegetables and Vegetable Products	1999634
8d7de19b-30f3-4cfd-ae93-c33a8f19a18d	2024-06-15 23:36:56.172512	salt	g	2047	Spices and Herbs	746775
69332cc2-7b6f-42aa-be4d-c2ac2f2954c0	2024-07-01 15:30:32.231656	Spaghetti	g	0	Pasta by Shape & Type	2099117
5e8cd4c6-51aa-42aa-ac24-ac3997c73341	2024-07-02 14:11:08.757873	Pancetta	g	0	Pepperoni, Salami & Cold Cuts	2098421
ea3f9073-6a75-4625-80d1-19dc42aca7ef	2024-07-02 14:14:10.412584	Egg	piece	1123	Dairy and Egg Products	748967
db630404-6115-4ca1-91cd-f9ed8981676f	2024-07-02 14:16:37.660621	Parmesan cheese	g	1032	Dairy and Egg Products	325036
567e990a-20cf-4f85-974f-38189c0bb64b	2024-07-02 14:18:15.185307	Garlic	g	11215	Vegetables and Vegetable Products	1104647
c1ac47d6-2126-48a4-ad73-75da637dee65	2024-07-02 14:23:16.532751	Black pepper	g	0	Spices and Herbes	2157235
\.


--
-- Data for Name: nutritional_value; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.nutritional_value (id, created_at, ingredient_id, recipe_id, kcal, kj, fat, saturated_fat, carbohydrate, sugar, protein, salt, nutriscore) FROM stdin;
002e607b-8b82-4b29-87a9-bfe50ee30433	2024-06-16 11:37:49.95385	84eb6da1-25b9-40ec-97a1-c0db1844ca54	\N	22	92	0.42	\N	3.84	\N	0.7	\N	\N
\.


--
-- Data for Name: rating; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.rating (id, created_at, ingredient_id, recipe_id, overall, mon, tue, wed, thu, fri, sat, sun, win, spr, sum, aut, thirtydegree, twentiedegree, tendegree, zerodegree, subzerodegree) FROM stdin;
489a1860-c881-4057-8607-b6826abfb2cd	2024-07-01 15:30:32.231656	69332cc2-7b6f-42aa-be4d-c2ac2f2954c0	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
97f5f257-d319-4fc7-b939-063858e1ae04	2024-07-02 14:11:08.757873	5e8cd4c6-51aa-42aa-ac24-ac3997c73341	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
e0ba6c46-36ae-450f-ac62-5a133d51df1f	2024-07-02 14:14:10.412584	ea3f9073-6a75-4625-80d1-19dc42aca7ef	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
58081d5d-73ec-4078-b8c4-58db512c8a2f	2024-07-02 14:16:37.660621	db630404-6115-4ca1-91cd-f9ed8981676f	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
bf64206b-53dc-4b9a-a5e5-2f14fafdd4e4	2024-07-02 14:18:15.185307	567e990a-20cf-4f85-974f-38189c0bb64b	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
8e98f6f5-aace-44a7-a09f-e3a71ca12329	2024-07-02 14:23:16.532751	c1ac47d6-2126-48a4-ad73-75da637dee65	\N	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
b0fe7647-5756-4c6c-b576-afc740b8f025	2024-07-24 15:49:43.879625	\N	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	1238.675	1000	1000	1000	1000	1000	1331	1771.6	1000	1000	2358.1	1000	1000	1000	2358.1	1000	1000
eb3ce3eb-e90a-4bc9-a3d7-744e05b7e201	2024-09-01 20:32:48.395312	\N	c4ef5707-1577-4f8c-99ef-0f492e82b895	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000	1000
\.


--
-- Data for Name: recipe_ingredient; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.recipe_ingredient (id, created_at, recipe_id, ingredient_id, amount, unit) FROM stdin;
185ae84d-4fe5-4328-ba1d-7af4434cb521	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	69332cc2-7b6f-42aa-be4d-c2ac2f2954c0	400	g
2de9c1c6-cc35-4038-8fbc-17029984f1d8	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	5e8cd4c6-51aa-42aa-ac24-ac3997c73341	150	g
c2f50f80-71dd-4374-a856-bf417a26a5eb	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	ea3f9073-6a75-4625-80d1-19dc42aca7ef	4	large
ed5fbfb6-2d2d-4467-82cf-9e7a97924724	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	db630404-6115-4ca1-91cd-f9ed8981676f	100	g
a4dd3925-e377-4380-8f0c-797d266b40e4	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	567e990a-20cf-4f85-974f-38189c0bb64b	2	cloves
69842c21-5832-4c64-9d27-2ffb8abd4617	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	8d7de19b-30f3-4cfd-ae93-c33a8f19a18d	1	tsp
07c807a0-15c2-4db5-8ca6-836499825c46	2024-07-24 15:49:43.879625	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	c1ac47d6-2126-48a4-ad73-75da637dee65	1	tsp
3d88c865-916c-406c-bb28-4252484a8744	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	69332cc2-7b6f-42aa-be4d-c2ac2f2954c0	400	g
652abb1a-9014-47de-b5c9-c8c43721b686	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	5e8cd4c6-51aa-42aa-ac24-ac3997c73341	150	g
6c16fc93-514d-4ffb-98f7-194c1d36257b	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	ea3f9073-6a75-4625-80d1-19dc42aca7ef	4	large
d973f0cd-f55f-401d-ba65-7f9bafe5240e	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	db630404-6115-4ca1-91cd-f9ed8981676f	100	g
9cb0edb9-c4dc-4499-8b4e-da5c7daad99d	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	567e990a-20cf-4f85-974f-38189c0bb64b	2	cloves
2125e6f4-04f0-4a2b-a54c-7dbde1b6641e	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	8d7de19b-30f3-4cfd-ae93-c33a8f19a18d	1	tsp
0daaf384-c8bd-419f-b31d-2d3442295a2a	2024-09-01 20:32:48.395312	c4ef5707-1577-4f8c-99ef-0f492e82b895	c1ac47d6-2126-48a4-ad73-75da637dee65	1	tsp
\.


--
-- Data for Name: recipe_selects_views_log; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.recipe_selects_views_log (id, day, selects, views, recipe_id, view_change, selects_change) FROM stdin;
\.


--
-- Data for Name: recipes; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.recipes (id, created_at, author, name, cuisine, yield, yield_unit, prep_time, cooking_time, version, selects, views) FROM stdin;
aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	2024-07-24 15:49:43.879625	f85a98f8-2572-420a-9ae5-2c997ad96b6d	Classic Spaghetti Carbonara	italian	500		01:00:00	01:00:00	9	9	0
c4ef5707-1577-4f8c-99ef-0f492e82b895	2024-09-01 20:32:48.395312	f85a98f8-2572-420a-9ae5-2c997ad96b6d	Classic Spaghetti Carbonara 2	italian	500		01:00:00	01:00:00	0	0	0
\.


--
-- Data for Name: rel_diet_recipe; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.rel_diet_recipe (id, recipe_id, diet_id) FROM stdin;
a8e9c8b8-c857-49f3-b700-1b87592ae51f	c4ef5707-1577-4f8c-99ef-0f492e82b895	bbadd945-5557-459f-951e-9ad3ad277059
\.


--
-- Data for Name: rel_diet_user; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.rel_diet_user (id, user_id, diet_id) FROM stdin;
cbf679ff-539f-4078-ac6f-af7c9beac8e5	f85a98f8-2572-420a-9ae5-2c997ad96b6d	bbadd945-5557-459f-951e-9ad3ad277059
\.


--
-- Data for Name: step; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.step (id, created_at, step, recipe_id, technique_id, ingredient_id) FROM stdin;
705897bb-6ec9-4d5f-adfc-0a7b4fa471dc	2024-07-24 15:49:43.879625	Cook the spaghetti according to package directions until al dente. Reserve 1 cup of pasta water, then drain the pasta.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
13b29b7b-8ce8-44ba-90ae-c243c98da031	2024-07-24 15:49:43.879625	While the pasta cooks, heat a large skillet over medium heat and add the pancetta. Cook until crispy, then remove from heat and set aside.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
e8148c4f-6203-49e9-b50a-aa3a8545e808	2024-07-24 15:49:43.879625	In a bowl, whisk together the eggs and grated Parmesan cheese until well combined.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
926126e2-463b-436c-b574-5fbb646f82c8	2024-07-24 15:49:43.879625	Return the skillet with pancetta to low heat. Add the minced garlic and cook until fragrant, about 1 minute.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
f508edbc-4c63-4f7e-949c-2f89422d7ad9	2024-07-24 15:49:43.879625	Add the cooked pasta to the skillet and toss to combine with the pancetta and garlic.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
5e800272-2815-4220-ae2b-dc08c4ffc80b	2024-07-24 15:49:43.879625	Remove the skillet from heat and quickly pour in the egg and cheese mixture, tossing rapidly to create a creamy sauce. If the sauce is too thick, add a little reserved pasta water until desired consistency is reached.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
539c001b-aa6d-4e20-9721-9a34eef5cccc	2024-07-24 15:49:43.879625	Season with salt and freshly ground black pepper to taste. Serve immediately with extra Parmesan cheese on top, if desired.	aa85daf1-dbc5-462d-a6fe-3fbb358b08dd	\N	\N
6617a65c-bc28-4066-bc94-669ba5af86d7	2024-09-01 20:32:48.395312	Cook the spaghetti according to package directions until al dente. Reserve 1 cup of pasta water, then drain the pasta.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
02e95737-52bb-4a52-b234-e1bda0d214ab	2024-09-01 20:32:48.395312	While the pasta cooks, heat a large skillet over medium heat and add the pancetta. Cook until crispy, then remove from heat and set aside.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
0b271eb6-498d-4d4f-a33a-a9dcbf7605cc	2024-09-01 20:32:48.395312	In a bowl, whisk together the eggs and grated Parmesan cheese until well combined.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
c9910490-db84-44a6-bc04-39e602968c7c	2024-09-01 20:32:48.395312	Return the skillet with pancetta to low heat. Add the minced garlic and cook until fragrant, about 1 minute.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
567a04dd-92b7-4e00-a75a-2271199d9f47	2024-09-01 20:32:48.395312	Add the cooked pasta to the skillet and toss to combine with the pancetta and garlic.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
463158ac-cd78-4294-ab01-f86c41fac8e5	2024-09-01 20:32:48.395312	Remove the skillet from heat and quickly pour in the egg and cheese mixture, tossing rapidly to create a creamy sauce. If the sauce is too thick, add a little reserved pasta water until desired consistency is reached.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
ba1b9c28-5f7c-496a-84c0-ba0de4ad482d	2024-09-01 20:32:48.395312	Season with salt and freshly ground black pepper to taste. Serve immediately with extra Parmesan cheese on top, if desired.	c4ef5707-1577-4f8c-99ef-0f492e82b895	\N	\N
\.


--
-- Data for Name: technique; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public.technique (id, created_at, description) FROM stdin;
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: mads
--

COPY public."user" (id, created_at, cookie, ip, groups) FROM stdin;
f85a98f8-2572-420a-9ae5-2c997ad96b6d	2024-07-21 22:21:31.536743	6uZEqNNvlGeQOCO9fIvY	127.0.0.1	\N
\.


--
-- Name: ingredient con_unique_name; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.ingredient
    ADD CONSTRAINT con_unique_name UNIQUE (name);


--
-- Name: diet diet_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.diet
    ADD CONSTRAINT diet_pkey PRIMARY KEY (id);


--
-- Name: diet diet_unique; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.diet
    ADD CONSTRAINT diet_unique UNIQUE (name);


--
-- Name: ingredient ingredient_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.ingredient
    ADD CONSTRAINT ingredient_pkey PRIMARY KEY (id);


--
-- Name: nutritional_value nutritional_value_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.nutritional_value
    ADD CONSTRAINT nutritional_value_pkey PRIMARY KEY (id);


--
-- Name: nutritional_value nutritional_value_unique_fk_ingredient; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.nutritional_value
    ADD CONSTRAINT nutritional_value_unique_fk_ingredient UNIQUE (ingredient_id);


--
-- Name: nutritional_value nutritional_value_unique_fk_recipe; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.nutritional_value
    ADD CONSTRAINT nutritional_value_unique_fk_recipe UNIQUE (recipe_id);


--
-- Name: rating rating_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rating
    ADD CONSTRAINT rating_pkey PRIMARY KEY (id);


--
-- Name: rating rating_unique_fk_ingredient; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rating
    ADD CONSTRAINT rating_unique_fk_ingredient UNIQUE (ingredient_id);


--
-- Name: rating rating_unique_fk_recipe; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rating
    ADD CONSTRAINT rating_unique_fk_recipe UNIQUE (recipe_id);


--
-- Name: recipe_ingredient recipe_ingredient_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipe_ingredient
    ADD CONSTRAINT recipe_ingredient_pkey PRIMARY KEY (id);


--
-- Name: recipe_selects_views_log recipe_selected_view_log_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipe_selects_views_log
    ADD CONSTRAINT recipe_selected_view_log_pkey PRIMARY KEY (id);


--
-- Name: recipes recipes_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_pkey PRIMARY KEY (id);


--
-- Name: rel_diet_recipe rel_diet_recipe_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_recipe
    ADD CONSTRAINT rel_diet_recipe_pkey PRIMARY KEY (id);


--
-- Name: rel_diet_user rel_diet_user_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_user
    ADD CONSTRAINT rel_diet_user_pkey PRIMARY KEY (id);


--
-- Name: step step_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.step
    ADD CONSTRAINT step_pkey PRIMARY KEY (id);


--
-- Name: technique technique_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.technique
    ADD CONSTRAINT technique_pkey PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: fki_fk_recipe_user; Type: INDEX; Schema: public; Owner: mads
--

CREATE INDEX fki_fk_recipe_user ON public.recipes USING btree (author);


--
-- Name: rel_diet_recipe fk_diet_rel_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_recipe
    ADD CONSTRAINT fk_diet_rel_recipe FOREIGN KEY (diet_id) REFERENCES public.diet(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rel_diet_user fk_diet_rel_user; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_user
    ADD CONSTRAINT fk_diet_rel_user FOREIGN KEY (diet_id) REFERENCES public.diet(id);


--
-- Name: recipe_selects_views_log fk_log_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipe_selects_views_log
    ADD CONSTRAINT fk_log_recipe FOREIGN KEY (recipe_id) REFERENCES public.recipes(id);


--
-- Name: nutritional_value fk_nutritional_value_ingredient; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.nutritional_value
    ADD CONSTRAINT fk_nutritional_value_ingredient FOREIGN KEY (ingredient_id) REFERENCES public.ingredient(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: nutritional_value fk_nutritional_value_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.nutritional_value
    ADD CONSTRAINT fk_nutritional_value_recipe FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: rating fk_rating_ingredient; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rating
    ADD CONSTRAINT fk_rating_ingredient FOREIGN KEY (ingredient_id) REFERENCES public.ingredient(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rating fk_rating_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rating
    ADD CONSTRAINT fk_rating_recipe FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: recipe_ingredient fk_recipe_ingredient_ingredient; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipe_ingredient
    ADD CONSTRAINT fk_recipe_ingredient_ingredient FOREIGN KEY (ingredient_id) REFERENCES public.ingredient(id) MATCH FULL ON UPDATE CASCADE;


--
-- Name: recipe_ingredient fk_recipe_ingredient_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipe_ingredient
    ADD CONSTRAINT fk_recipe_ingredient_recipe FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rel_diet_recipe fk_recipe_rel_diet; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_recipe
    ADD CONSTRAINT fk_recipe_rel_diet FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: recipes fk_recipe_user; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT fk_recipe_user FOREIGN KEY (author) REFERENCES public."user"(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: step fk_step_recipe; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.step
    ADD CONSTRAINT fk_step_recipe FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: step fk_step_recipe_ingredient; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.step
    ADD CONSTRAINT fk_step_recipe_ingredient FOREIGN KEY (ingredient_id) REFERENCES public.recipe_ingredient(id) ON DELETE SET NULL NOT VALID;


--
-- Name: step fk_step_technique; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.step
    ADD CONSTRAINT fk_step_technique FOREIGN KEY (technique_id) REFERENCES public.technique(id) ON DELETE SET NULL NOT VALID;


--
-- Name: rel_diet_user fk_user_rel_diet; Type: FK CONSTRAINT; Schema: public; Owner: mads
--

ALTER TABLE ONLY public.rel_diet_user
    ADD CONSTRAINT fk_user_rel_diet FOREIGN KEY (user_id) REFERENCES public."user"(id) NOT VALID;


--
-- PostgreSQL database dump complete
--

