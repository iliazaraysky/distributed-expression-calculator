--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

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
-- Name: requests; Type: TABLE; Schema: public; Owner: calculator_db_user
--

CREATE TABLE public.requests (
    id integer NOT NULL,
    unique_id uuid,
    query_text text,
    creation_time timestamp with time zone,
    completion_time timestamp with time zone,
    server_name character varying(50),
    username character varying(50),
    result text,
    execution_time interval,
    status text
);


ALTER TABLE public.requests OWNER TO calculator_db_user;

--
-- Name: requests_id_seq; Type: SEQUENCE; Schema: public; Owner: calculator_db_user
--

CREATE SEQUENCE public.requests_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.requests_id_seq OWNER TO calculator_db_user;

--
-- Name: requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: calculator_db_user
--

ALTER SEQUENCE public.requests_id_seq OWNED BY public.requests.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: calculator_db_user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character(50) NOT NULL,
    password character(255) NOT NULL,
    email character(100),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO calculator_db_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: calculator_db_user
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO calculator_db_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: calculator_db_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: workers; Type: TABLE; Schema: public; Owner: calculator_db_user
--

CREATE TABLE public.workers (
    id integer NOT NULL,
    name character varying,
    timer_setup_date timestamp with time zone,
    status character varying,
    last_task uuid,
    timeout integer
);


ALTER TABLE public.workers OWNER TO calculator_db_user;

--
-- Name: workers_id_seq; Type: SEQUENCE; Schema: public; Owner: calculator_db_user
--

CREATE SEQUENCE public.workers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workers_id_seq OWNER TO calculator_db_user;

--
-- Name: workers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: calculator_db_user
--

ALTER SEQUENCE public.workers_id_seq OWNED BY public.workers.id;


--
-- Name: requests id; Type: DEFAULT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.requests ALTER COLUMN id SET DEFAULT nextval('public.requests_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: workers id; Type: DEFAULT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.workers ALTER COLUMN id SET DEFAULT nextval('public.workers_id_seq'::regclass);


--
-- Data for Name: requests; Type: TABLE DATA; Schema: public; Owner: calculator_db_user
--

COPY public.requests (id, unique_id, query_text, creation_time, completion_time, server_name, username, result, execution_time, status) FROM stdin;
1	abb1f26f-2014-43f6-9f88-22fd6bf7aedc	2 + 1 + 3 + 1	2024-02-11 21:17:24.764203+00	2024-02-11 21:17:36.870463+00	worker3	adm	7	00:00:12.10626	Done
2	70d9df9f-84ec-426a-b408-31fe7fe8380f	2 + 2 + 1 + 12	2024-02-11 21:17:34.333891+00	2024-02-11 21:17:48.561868+00	worker2	adm	17	3952215:46:40	Done
3	a58338de-a36c-4013-aa52-c51f755da11f	25 + 2 + 1 + 12	2024-02-11 21:17:39.952909+00	2024-02-11 21:18:01.004218+00	worker1	adm	40	00:00:21.051309	Done
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: calculator_db_user
--

COPY public.users (id, login, password, email, created_at) FROM stdin;
1	test_user                                         	longPassswordForTest                                                                                                                                                                                                                                           	test_user@distributed-expression-calculator.com                                                     	2024-04-21 18:36:53.831993
\.


--
-- Data for Name: workers; Type: TABLE DATA; Schema: public; Owner: calculator_db_user
--

COPY public.workers (id, name, timer_setup_date, status, last_task, timeout) FROM stdin;
1	worker1	2024-02-18 19:12:39.952909+00	ready	a58338de-a36c-4013-aa52-c51f755da11f	10
2	worker2	2024-02-18 18:17:39.952909+00	ready	70d9df9f-84ec-426a-b408-31fe7fe8380f	11
3	worker3	2024-02-18 19:17:39.952909+00	ready	abb1f26f-2014-43f6-9f88-22fd6bf7aedc	12
\.


--
-- Name: requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: calculator_db_user
--

SELECT pg_catalog.setval('public.requests_id_seq', 3, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: calculator_db_user
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: workers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: calculator_db_user
--

SELECT pg_catalog.setval('public.workers_id_seq', 3, true);


--
-- Name: requests requests_pkey; Type: CONSTRAINT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: workers workers_pkey; Type: CONSTRAINT; Schema: public; Owner: calculator_db_user
--

ALTER TABLE ONLY public.workers
    ADD CONSTRAINT workers_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

