--
-- PostgreSQL database dump
--

-- Dumped from database version 15.0
-- Dumped by pg_dump version 15.0

-- Started on 2022-11-19 03:04:10 UTC

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

--
-- TOC entry 6 (class 2615 OID 16566)
-- Name: goradd_unit; Type: SCHEMA; Schema: -; Owner: root
--

CREATE SCHEMA goradd_unit;


ALTER SCHEMA goradd_unit OWNER TO root;

--
-- TOC entry 862 (class 1247 OID 16576)
-- Name: unsupported_types_type_enum; Type: TYPE; Schema: goradd_unit; Owner: root
--

CREATE TYPE goradd_unit.unsupported_types_type_enum AS ENUM (
    'a',
    'b',
    'c'
);


ALTER TYPE goradd_unit.unsupported_types_type_enum OWNER TO root;

--
-- TOC entry 859 (class 1247 OID 16568)
-- Name: unsupported_types_type_set; Type: TYPE; Schema: goradd_unit; Owner: root
--

CREATE TYPE goradd_unit.unsupported_types_type_set AS ENUM (
    'a',
    'b',
    'c'
);


ALTER TYPE goradd_unit.unsupported_types_type_set OWNER TO root;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 16583)
-- Name: double_index; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.double_index (
    id integer NOT NULL,
    field_int integer NOT NULL,
    field_string character varying(50) NOT NULL
);


ALTER TABLE goradd_unit.double_index OWNER TO root;

--
-- TOC entry 219 (class 1259 OID 16587)
-- Name: forward_cascade; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_cascade (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id integer
);


ALTER TABLE goradd_unit.forward_cascade OWNER TO root;

--
-- TOC entry 218 (class 1259 OID 16586)
-- Name: forward_cascade_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_cascade_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_cascade_id_seq OWNER TO root;

--
-- TOC entry 3459 (class 0 OID 0)
-- Dependencies: 218
-- Name: forward_cascade_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_cascade_id_seq OWNED BY goradd_unit.forward_cascade.id;


--
-- TOC entry 221 (class 1259 OID 16592)
-- Name: forward_cascade_unique; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_cascade_unique (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id integer
);


ALTER TABLE goradd_unit.forward_cascade_unique OWNER TO root;

--
-- TOC entry 220 (class 1259 OID 16591)
-- Name: forward_cascade_unique_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_cascade_unique_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_cascade_unique_id_seq OWNER TO root;

--
-- TOC entry 3460 (class 0 OID 0)
-- Dependencies: 220
-- Name: forward_cascade_unique_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_cascade_unique_id_seq OWNED BY goradd_unit.forward_cascade_unique.id;


--
-- TOC entry 223 (class 1259 OID 16597)
-- Name: forward_null; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_null (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id integer
);


ALTER TABLE goradd_unit.forward_null OWNER TO root;

--
-- TOC entry 222 (class 1259 OID 16596)
-- Name: forward_null_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_null_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_null_id_seq OWNER TO root;

--
-- TOC entry 3461 (class 0 OID 0)
-- Dependencies: 222
-- Name: forward_null_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_null_id_seq OWNED BY goradd_unit.forward_null.id;


--
-- TOC entry 225 (class 1259 OID 16602)
-- Name: forward_null_unique; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_null_unique (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id integer
);


ALTER TABLE goradd_unit.forward_null_unique OWNER TO root;

--
-- TOC entry 224 (class 1259 OID 16601)
-- Name: forward_null_unique_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_null_unique_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_null_unique_id_seq OWNER TO root;

--
-- TOC entry 3462 (class 0 OID 0)
-- Dependencies: 224
-- Name: forward_null_unique_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_null_unique_id_seq OWNED BY goradd_unit.forward_null_unique.id;


--
-- TOC entry 227 (class 1259 OID 16607)
-- Name: forward_restrict; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_restrict (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id bigint NOT NULL
);


ALTER TABLE goradd_unit.forward_restrict OWNER TO root;

--
-- TOC entry 226 (class 1259 OID 16606)
-- Name: forward_restrict_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_restrict_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_restrict_id_seq OWNER TO root;

--
-- TOC entry 3463 (class 0 OID 0)
-- Dependencies: 226
-- Name: forward_restrict_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_restrict_id_seq OWNED BY goradd_unit.forward_restrict.id;


--
-- TOC entry 229 (class 1259 OID 16612)
-- Name: forward_restrict_unique; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.forward_restrict_unique (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    reverse_id integer
);


ALTER TABLE goradd_unit.forward_restrict_unique OWNER TO root;

--
-- TOC entry 228 (class 1259 OID 16611)
-- Name: forward_restrict_unique_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.forward_restrict_unique_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.forward_restrict_unique_id_seq OWNER TO root;

--
-- TOC entry 3464 (class 0 OID 0)
-- Dependencies: 228
-- Name: forward_restrict_unique_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.forward_restrict_unique_id_seq OWNED BY goradd_unit.forward_restrict_unique.id;


--
-- TOC entry 231 (class 1259 OID 16617)
-- Name: reverse; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.reverse (
    id integer NOT NULL,
    name character varying(100) NOT NULL
);


ALTER TABLE goradd_unit.reverse OWNER TO root;

--
-- TOC entry 230 (class 1259 OID 16616)
-- Name: reverse_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.reverse_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.reverse_id_seq OWNER TO root;

--
-- TOC entry 3465 (class 0 OID 0)
-- Dependencies: 230
-- Name: reverse_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.reverse_id_seq OWNED BY goradd_unit.reverse.id;


--
-- TOC entry 232 (class 1259 OID 16621)
-- Name: two_key; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.two_key (
    server character varying(50) NOT NULL,
    directory character varying(50) NOT NULL,
    file_name character varying(50) NOT NULL
);


ALTER TABLE goradd_unit.two_key OWNER TO root;

--
-- TOC entry 234 (class 1259 OID 16625)
-- Name: type_test; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.type_test (
    id bigint NOT NULL,
    date date,
    "time" time without time zone,
    date_time timestamp with time zone,
    ts timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    test_int integer DEFAULT 5,
    test_float double precision,
    test_double double precision NOT NULL,
    test_text text,
    test_bit boolean,
    test_varchar character varying(10) DEFAULT NULL::character varying,
    test_blob bytea NOT NULL
);


ALTER TABLE goradd_unit.type_test OWNER TO root;

--
-- TOC entry 233 (class 1259 OID 16624)
-- Name: type_test_id_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.type_test_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.type_test_id_seq OWNER TO root;

--
-- TOC entry 3466 (class 0 OID 0)
-- Dependencies: 233
-- Name: type_test_id_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.type_test_id_seq OWNED BY goradd_unit.type_test.id;


--
-- TOC entry 236 (class 1259 OID 16635)
-- Name: unsupported_types; Type: TABLE; Schema: goradd_unit; Owner: root
--

CREATE TABLE goradd_unit.unsupported_types (
    type_set goradd_unit.unsupported_types_type_set[] NOT NULL,
    type_enum goradd_unit.unsupported_types_type_enum NOT NULL,
    type_decimal numeric(10,4) NOT NULL,
    type_double double precision NOT NULL,
    type_geo point NOT NULL,
    type_tiny_blob bytea NOT NULL,
    type_medium_blob bytea NOT NULL,
    type_varbinary bytea NOT NULL,
    type_longtext text NOT NULL,
    type_binary bytea NOT NULL,
    type_small smallint NOT NULL,
    type_medium integer NOT NULL,
    type_big bigint NOT NULL,
    type_polygon polygon NOT NULL,
    type_serial bigint NOT NULL,
    type_unsigned bigint NOT NULL,
    type_multfk1 character varying(50) NOT NULL,
    type_multifk2 character varying(50) NOT NULL
);


ALTER TABLE goradd_unit.unsupported_types OWNER TO root;

--
-- TOC entry 235 (class 1259 OID 16634)
-- Name: unsupported_types_type_serial_seq; Type: SEQUENCE; Schema: goradd_unit; Owner: root
--

CREATE SEQUENCE goradd_unit.unsupported_types_type_serial_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goradd_unit.unsupported_types_type_serial_seq OWNER TO root;

--
-- TOC entry 3467 (class 0 OID 0)
-- Dependencies: 235
-- Name: unsupported_types_type_serial_seq; Type: SEQUENCE OWNED BY; Schema: goradd_unit; Owner: root
--

ALTER SEQUENCE goradd_unit.unsupported_types_type_serial_seq OWNED BY goradd_unit.unsupported_types.type_serial;


--
-- TOC entry 3245 (class 2604 OID 25093)
-- Name: forward_cascade id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_cascade_id_seq'::regclass);


--
-- TOC entry 3246 (class 2604 OID 25112)
-- Name: forward_cascade_unique id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade_unique ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_cascade_unique_id_seq'::regclass);


--
-- TOC entry 3247 (class 2604 OID 25131)
-- Name: forward_null id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_null_id_seq'::regclass);


--
-- TOC entry 3248 (class 2604 OID 25150)
-- Name: forward_null_unique id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null_unique ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_null_unique_id_seq'::regclass);


--
-- TOC entry 3249 (class 2604 OID 25169)
-- Name: forward_restrict id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_restrict_id_seq'::regclass);


--
-- TOC entry 3250 (class 2604 OID 25177)
-- Name: forward_restrict_unique id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict_unique ALTER COLUMN id SET DEFAULT nextval('goradd_unit.forward_restrict_unique_id_seq'::regclass);


--
-- TOC entry 3251 (class 2604 OID 25197)
-- Name: reverse id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.reverse ALTER COLUMN id SET DEFAULT nextval('goradd_unit.reverse_id_seq'::regclass);


--
-- TOC entry 3252 (class 2604 OID 16628)
-- Name: type_test id; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.type_test ALTER COLUMN id SET DEFAULT nextval('goradd_unit.type_test_id_seq'::regclass);


--
-- TOC entry 3256 (class 2604 OID 16638)
-- Name: unsupported_types type_serial; Type: DEFAULT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.unsupported_types ALTER COLUMN type_serial SET DEFAULT nextval('goradd_unit.unsupported_types_type_serial_seq'::regclass);


--
-- TOC entry 3434 (class 0 OID 16583)
-- Dependencies: 217
-- Data for Name: double_index; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.double_index (id, field_int, field_string) FROM stdin;
\.


--
-- TOC entry 3436 (class 0 OID 16587)
-- Dependencies: 219
-- Data for Name: forward_cascade; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_cascade (id, name, reverse_id) FROM stdin;
\.


--
-- TOC entry 3438 (class 0 OID 16592)
-- Dependencies: 221
-- Data for Name: forward_cascade_unique; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_cascade_unique (id, name, reverse_id) FROM stdin;
\.


--
-- TOC entry 3440 (class 0 OID 16597)
-- Dependencies: 223
-- Data for Name: forward_null; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_null (id, name, reverse_id) FROM stdin;
34	testForward1	\N
35	Other	123
36	testForward1	\N
37	Other	124
38	testForward1	\N
39	Other	125
40	testForward1	126
41	Other	126
42	testForward1	\N
43	Other	127
44	testForward3	127
45	testForward1	128
46	Other	128
47	testForward1	129
48	Other	129
51	Other	132
52	testForward3	132
53	testForward1	133
54	Other	133
\.


--
-- TOC entry 3442 (class 0 OID 16602)
-- Dependencies: 225
-- Data for Name: forward_null_unique; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_null_unique (id, name, reverse_id) FROM stdin;
\.


--
-- TOC entry 3444 (class 0 OID 16607)
-- Dependencies: 227
-- Data for Name: forward_restrict; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_restrict (id, name, reverse_id) FROM stdin;
\.


--
-- TOC entry 3446 (class 0 OID 16612)
-- Dependencies: 229
-- Data for Name: forward_restrict_unique; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.forward_restrict_unique (id, name, reverse_id) FROM stdin;
\.


--
-- TOC entry 3448 (class 0 OID 16617)
-- Dependencies: 231
-- Data for Name: reverse; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.reverse (id, name) FROM stdin;
123	testReverse
124	testReverse
125	testReverse
126	testReverse
127	testReverse
128	testReverse
129	testReverse
132	testReverse
133	testReverse
\.


--
-- TOC entry 3449 (class 0 OID 16621)
-- Dependencies: 232
-- Data for Name: two_key; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.two_key (server, directory, file_name) FROM stdin;
cnn.com	us	news
google.com	drive	
google.com	mail	mail.html
google.com	news	news.php
mail.google.com	mail	inbox
yahoo.com		
\.


--
-- TOC entry 3451 (class 0 OID 16625)
-- Dependencies: 234
-- Data for Name: type_test; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.type_test (id, date, "time", date_time, ts, test_int, test_float, test_double, test_text, test_bit, test_varchar, test_blob) FROM stdin;
1	2019-01-02	06:17:28	2019-01-02 06:17:28+00	2002-07-02 14:04:03+00	5	1.2	3.33	Sample	t	Sample	\\x61626364
\.


--
-- TOC entry 3453 (class 0 OID 16635)
-- Dependencies: 236
-- Data for Name: unsupported_types; Type: TABLE DATA; Schema: goradd_unit; Owner: root
--

COPY goradd_unit.unsupported_types (type_set, type_enum, type_decimal, type_double, type_geo, type_tiny_blob, type_medium_blob, type_varbinary, type_longtext, type_binary, type_small, type_medium, type_big, type_polygon, type_serial, type_unsigned, type_multfk1, type_multifk2) FROM stdin;
\.


--
-- TOC entry 3468 (class 0 OID 0)
-- Dependencies: 218
-- Name: forward_cascade_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_cascade_id_seq', 1, true);


--
-- TOC entry 3469 (class 0 OID 0)
-- Dependencies: 220
-- Name: forward_cascade_unique_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_cascade_unique_id_seq', 55, true);


--
-- TOC entry 3470 (class 0 OID 0)
-- Dependencies: 222
-- Name: forward_null_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_null_id_seq', 126, true);


--
-- TOC entry 3471 (class 0 OID 0)
-- Dependencies: 224
-- Name: forward_null_unique_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_null_unique_id_seq', 55, true);


--
-- TOC entry 3472 (class 0 OID 0)
-- Dependencies: 226
-- Name: forward_restrict_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_restrict_id_seq', 73, true);


--
-- TOC entry 3473 (class 0 OID 0)
-- Dependencies: 228
-- Name: forward_restrict_unique_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.forward_restrict_unique_id_seq', 55, true);


--
-- TOC entry 3474 (class 0 OID 0)
-- Dependencies: 230
-- Name: reverse_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.reverse_id_seq', 403, true);


--
-- TOC entry 3475 (class 0 OID 0)
-- Dependencies: 233
-- Name: type_test_id_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.type_test_id_seq', 1, true);


--
-- TOC entry 3476 (class 0 OID 0)
-- Dependencies: 235
-- Name: unsupported_types_type_serial_seq; Type: SEQUENCE SET; Schema: goradd_unit; Owner: root
--

SELECT pg_catalog.setval('goradd_unit.unsupported_types_type_serial_seq', 1, true);


--
-- TOC entry 3258 (class 2606 OID 25081)
-- Name: double_index idx_16583_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.double_index
    ADD CONSTRAINT idx_16583_primary PRIMARY KEY (id);


--
-- TOC entry 3261 (class 2606 OID 25095)
-- Name: forward_cascade idx_16587_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade
    ADD CONSTRAINT idx_16587_primary PRIMARY KEY (id);


--
-- TOC entry 3264 (class 2606 OID 25114)
-- Name: forward_cascade_unique idx_16592_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade_unique
    ADD CONSTRAINT idx_16592_primary PRIMARY KEY (id);


--
-- TOC entry 3267 (class 2606 OID 25133)
-- Name: forward_null idx_16597_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null
    ADD CONSTRAINT idx_16597_primary PRIMARY KEY (id);


--
-- TOC entry 3270 (class 2606 OID 25152)
-- Name: forward_null_unique idx_16602_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null_unique
    ADD CONSTRAINT idx_16602_primary PRIMARY KEY (id);


--
-- TOC entry 3273 (class 2606 OID 25171)
-- Name: forward_restrict idx_16607_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict
    ADD CONSTRAINT idx_16607_primary PRIMARY KEY (id);


--
-- TOC entry 3276 (class 2606 OID 25179)
-- Name: forward_restrict_unique idx_16612_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict_unique
    ADD CONSTRAINT idx_16612_primary PRIMARY KEY (id);


--
-- TOC entry 3279 (class 2606 OID 25199)
-- Name: reverse idx_16617_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.reverse
    ADD CONSTRAINT idx_16617_primary PRIMARY KEY (id);


--
-- TOC entry 3281 (class 2606 OID 16668)
-- Name: two_key idx_16621_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.two_key
    ADD CONSTRAINT idx_16621_primary PRIMARY KEY (server, directory);


--
-- TOC entry 3283 (class 2606 OID 16673)
-- Name: type_test idx_16625_primary; Type: CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.type_test
    ADD CONSTRAINT idx_16625_primary PRIMARY KEY (id);


--
-- TOC entry 3262 (class 1259 OID 25101)
-- Name: idx_16587_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE INDEX idx_16587_reverse_id ON goradd_unit.forward_cascade USING btree (reverse_id);


--
-- TOC entry 3265 (class 1259 OID 25120)
-- Name: idx_16592_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE UNIQUE INDEX idx_16592_reverse_id ON goradd_unit.forward_cascade_unique USING btree (reverse_id);


--
-- TOC entry 3268 (class 1259 OID 25139)
-- Name: idx_16597_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE INDEX idx_16597_reverse_id ON goradd_unit.forward_null USING btree (reverse_id);


--
-- TOC entry 3271 (class 1259 OID 25158)
-- Name: idx_16602_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE UNIQUE INDEX idx_16602_reverse_id ON goradd_unit.forward_null_unique USING btree (reverse_id);


--
-- TOC entry 3274 (class 1259 OID 16652)
-- Name: idx_16607_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE INDEX idx_16607_reverse_id ON goradd_unit.forward_restrict USING btree (reverse_id);


--
-- TOC entry 3277 (class 1259 OID 25185)
-- Name: idx_16612_reverse_id; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE UNIQUE INDEX idx_16612_reverse_id ON goradd_unit.forward_restrict_unique USING btree (reverse_id);


--
-- TOC entry 3284 (class 1259 OID 16657)
-- Name: idx_16635_type_multfk1; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE INDEX idx_16635_type_multfk1 ON goradd_unit.unsupported_types USING btree (type_multfk1, type_multifk2);


--
-- TOC entry 3285 (class 1259 OID 16659)
-- Name: idx_16635_type_serial; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE UNIQUE INDEX idx_16635_type_serial ON goradd_unit.unsupported_types USING btree (type_serial);


--
-- TOC entry 3259 (class 1259 OID 25087)
-- Name: idx_dbl; Type: INDEX; Schema: goradd_unit; Owner: root
--

CREATE UNIQUE INDEX idx_dbl ON goradd_unit.double_index USING btree (field_int, field_string);


--
-- TOC entry 3286 (class 2606 OID 25205)
-- Name: forward_cascade forward_cascade_ibfk_1; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade
    ADD CONSTRAINT forward_cascade_ibfk_1 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3287 (class 2606 OID 25210)
-- Name: forward_cascade_unique forward_cascade_unique_ibfk_1; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_cascade_unique
    ADD CONSTRAINT forward_cascade_unique_ibfk_1 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3288 (class 2606 OID 25215)
-- Name: forward_null forward_null_ibfk_2; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null
    ADD CONSTRAINT forward_null_ibfk_2 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE SET NULL ON DELETE SET NULL;


--
-- TOC entry 3289 (class 2606 OID 25220)
-- Name: forward_null_unique forward_null_unique_ibfk_1; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_null_unique
    ADD CONSTRAINT forward_null_unique_ibfk_1 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE SET NULL ON DELETE SET NULL;


--
-- TOC entry 3290 (class 2606 OID 25200)
-- Name: forward_restrict forward_restrict_ibfk_1; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict
    ADD CONSTRAINT forward_restrict_ibfk_1 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- TOC entry 3291 (class 2606 OID 25225)
-- Name: forward_restrict_unique forward_restrict_unique_ibfk_1; Type: FK CONSTRAINT; Schema: goradd_unit; Owner: root
--

ALTER TABLE ONLY goradd_unit.forward_restrict_unique
    ADD CONSTRAINT forward_restrict_unique_ibfk_1 FOREIGN KEY (reverse_id) REFERENCES goradd_unit.reverse(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


-- Completed on 2022-11-19 03:04:11 UTC

--
-- PostgreSQL database dump complete
--

