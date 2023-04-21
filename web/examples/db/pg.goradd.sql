--
-- PostgreSQL database dump
--

-- Dumped from database version 15.0
-- Dumped by pg_dump version 15.2

-- Started on 2023-04-21 11:27:44 PDT

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
-- TOC entry 238 (class 1255 OID 16560)
-- Name: on_update_current_timestamp_person_with_lock(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.on_update_current_timestamp_person_with_lock() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
   NEW.sys_timestamp = now();
RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 16390)
-- Name: address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.address (
                                id integer NOT NULL,
                                person_id integer NOT NULL,
                                street character varying(100) NOT NULL,
                                city character varying(100) DEFAULT 'BOB'::character varying
);


--
-- TOC entry 216 (class 1259 OID 16389)
-- Name: address_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.address_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3483 (class 0 OID 0)
-- Dependencies: 216
-- Name: address_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.address_id_seq OWNED BY public.address.id;


--
-- TOC entry 219 (class 1259 OID 16396)
-- Name: employee_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.employee_info (
                                      id integer NOT NULL,
                                      person_id integer NOT NULL,
                                      employee_number integer NOT NULL
);


--
-- TOC entry 218 (class 1259 OID 16395)
-- Name: employee_info_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.employee_info_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3484 (class 0 OID 0)
-- Dependencies: 218
-- Name: employee_info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.employee_info_id_seq OWNED BY public.employee_info.id;


--
-- TOC entry 220 (class 1259 OID 16400)
-- Name: gift; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gift (
                             number integer NOT NULL,
                             name character varying(50) NOT NULL
);


--
-- TOC entry 3485 (class 0 OID 0)
-- Dependencies: 220
-- Name: TABLE gift; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.gift IS 'Table is keyed with an integer, but does not auto-increment';


--
-- TOC entry 222 (class 1259 OID 16404)
-- Name: login; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.login (
                              id integer NOT NULL,
                              person_id integer,
                              username character varying(20) NOT NULL,
                              password character varying(20) DEFAULT NULL::character varying,
                              is_enabled boolean DEFAULT true NOT NULL
);


--
-- TOC entry 221 (class 1259 OID 16403)
-- Name: login_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.login_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3486 (class 0 OID 0)
-- Dependencies: 221
-- Name: login_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.login_id_seq OWNED BY public.login.id;


--
-- TOC entry 224 (class 1259 OID 16411)
-- Name: milestone; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.milestone (
                                  id integer NOT NULL,
                                  project_id integer NOT NULL,
                                  name character varying(50) NOT NULL
);


--
-- TOC entry 223 (class 1259 OID 16410)
-- Name: milestone_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.milestone_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3487 (class 0 OID 0)
-- Dependencies: 223
-- Name: milestone_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.milestone_id_seq OWNED BY public.milestone.id;


--
-- TOC entry 226 (class 1259 OID 16416)
-- Name: person; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.person (
                               id integer NOT NULL,
                               first_name character varying(50) NOT NULL,
                               last_name character varying(50) NOT NULL
);


--
-- TOC entry 225 (class 1259 OID 16415)
-- Name: person_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.person_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3488 (class 0 OID 0)
-- Dependencies: 225
-- Name: person_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.person_id_seq OWNED BY public.person.id;


--
-- TOC entry 227 (class 1259 OID 16420)
-- Name: person_persontype_assn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.person_persontype_assn (
                                               person_id integer NOT NULL,
                                               person_type_id integer NOT NULL
);


--
-- TOC entry 229 (class 1259 OID 16424)
-- Name: person_type_enum; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.person_type_enum (
                                         id integer NOT NULL,
                                         name character varying(50) NOT NULL
);


--
-- TOC entry 228 (class 1259 OID 16423)
-- Name: person_type_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.person_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3489 (class 0 OID 0)
-- Dependencies: 228
-- Name: person_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.person_type_id_seq OWNED BY public.person_type_enum.id;


--
-- TOC entry 231 (class 1259 OID 16429)
-- Name: person_with_lock; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.person_with_lock (
                                         id integer NOT NULL,
                                         first_name character varying(50) NOT NULL,
                                         last_name character varying(50) NOT NULL,
                                         sys_timestamp timestamp with time zone
);


--
-- TOC entry 230 (class 1259 OID 16428)
-- Name: person_with_lock_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.person_with_lock_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3490 (class 0 OID 0)
-- Dependencies: 230
-- Name: person_with_lock_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.person_with_lock_id_seq OWNED BY public.person_with_lock.id;


--
-- TOC entry 233 (class 1259 OID 16434)
-- Name: project; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project (
                                id integer NOT NULL,
                                num integer NOT NULL,
                                status_id integer NOT NULL,
                                manager_id integer,
                                name character varying(100) NOT NULL,
                                description text,
                                start_date date,
                                end_date date,
                                budget numeric(12,2) DEFAULT NULL::numeric,
                                spent numeric(12,2) DEFAULT NULL::numeric
);


--
-- TOC entry 3491 (class 0 OID 0)
-- Dependencies: 233
-- Name: COLUMN project.num; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.num IS 'To simplify checking test results and as a non pk id test';


--
-- TOC entry 232 (class 1259 OID 16433)
-- Name: project_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.project_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3492 (class 0 OID 0)
-- Dependencies: 232
-- Name: project_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.project_id_seq OWNED BY public.project.id;


--
-- TOC entry 235 (class 1259 OID 16443)
-- Name: project_status_enum; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project_status_enum (
                                            id integer NOT NULL,
                                            name character varying(50) NOT NULL,
                                            description text,
                                            guidelines text,
                                            is_active boolean NOT NULL
);


--
-- TOC entry 234 (class 1259 OID 16442)
-- Name: project_status_type_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.project_status_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3493 (class 0 OID 0)
-- Dependencies: 234
-- Name: project_status_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.project_status_type_id_seq OWNED BY public.project_status_enum.id;


--
-- TOC entry 236 (class 1259 OID 16449)
-- Name: related_project_assn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.related_project_assn (
                                             parent_id integer NOT NULL,
                                             child_id integer NOT NULL
);


--
-- TOC entry 237 (class 1259 OID 16452)
-- Name: team_member_project_assn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.team_member_project_assn (
                                                 team_member_id integer NOT NULL,
                                                 project_id integer NOT NULL
);


--
-- TOC entry 3247 (class 2604 OID 24780)
-- Name: address id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address ALTER COLUMN id SET DEFAULT nextval('public.address_id_seq'::regclass);


--
-- TOC entry 3249 (class 2604 OID 24799)
-- Name: employee_info id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_info ALTER COLUMN id SET DEFAULT nextval('public.employee_info_id_seq'::regclass);


--
-- TOC entry 3250 (class 2604 OID 24829)
-- Name: login id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login ALTER COLUMN id SET DEFAULT nextval('public.login_id_seq'::regclass);


--
-- TOC entry 3253 (class 2604 OID 24850)
-- Name: milestone id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.milestone ALTER COLUMN id SET DEFAULT nextval('public.milestone_id_seq'::regclass);


--
-- TOC entry 3254 (class 2604 OID 24869)
-- Name: person id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person ALTER COLUMN id SET DEFAULT nextval('public.person_id_seq'::regclass);


--
-- TOC entry 3255 (class 2604 OID 24932)
-- Name: person_type_enum id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_type_enum ALTER COLUMN id SET DEFAULT nextval('public.person_type_id_seq'::regclass);


--
-- TOC entry 3256 (class 2604 OID 24945)
-- Name: person_with_lock id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_with_lock ALTER COLUMN id SET DEFAULT nextval('public.person_with_lock_id_seq'::regclass);


--
-- TOC entry 3257 (class 2604 OID 24952)
-- Name: project id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project ALTER COLUMN id SET DEFAULT nextval('public.project_id_seq'::regclass);


--
-- TOC entry 3260 (class 2604 OID 25015)
-- Name: project_status_enum id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project_status_enum ALTER COLUMN id SET DEFAULT nextval('public.project_status_type_id_seq'::regclass);


--
-- TOC entry 3457 (class 0 OID 16390)
-- Dependencies: 217
-- Data for Name: address; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.address VALUES (1, 1, '1 Love Drive', NULL);
INSERT INTO public.address VALUES (2, 2, '2 Doves and a Pine Cone Dr.', 'Dallas');
INSERT INTO public.address VALUES (3, 3, '3 Gold Fish Pl.', 'New York');
INSERT INTO public.address VALUES (4, 3, '323 W QCubed', 'New York');
INSERT INTO public.address VALUES (5, 5, '22 Elm St', 'Palo Alto');
INSERT INTO public.address VALUES (6, 7, '1 Pine St', 'San Jose');
INSERT INTO public.address VALUES (7, 7, '421 Central Expw', 'Mountain View');


--
-- TOC entry 3459 (class 0 OID 16396)
-- Dependencies: 219
-- Data for Name: employee_info; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- TOC entry 3460 (class 0 OID 16400)
-- Dependencies: 220
-- Data for Name: gift; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.gift VALUES (1, 'Partridge in a pear tree');
INSERT INTO public.gift VALUES (2, 'Turtle doves');
INSERT INTO public.gift VALUES (3, 'French hens');


--
-- TOC entry 3462 (class 0 OID 16404)
-- Dependencies: 222
-- Data for Name: login; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.login VALUES (1, 1, 'jdoe', 'p@$$.w0rd', false);
INSERT INTO public.login VALUES (2, 3, 'brobinson', 'p@$$.w0rd', true);
INSERT INTO public.login VALUES (3, 4, 'mho', 'p@$$.w0rd', true);
INSERT INTO public.login VALUES (4, 7, 'kwolfe', 'p@$$.w0rd', false);
INSERT INTO public.login VALUES (5, NULL, 'system', 'p@$$.w0rd', true);


--
-- TOC entry 3464 (class 0 OID 16411)
-- Dependencies: 224
-- Data for Name: milestone; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.milestone VALUES (1, 1, 'Milestone A');
INSERT INTO public.milestone VALUES (2, 1, 'Milestone B');
INSERT INTO public.milestone VALUES (3, 1, 'Milestone C');
INSERT INTO public.milestone VALUES (4, 2, 'Milestone D');
INSERT INTO public.milestone VALUES (5, 2, 'Milestone E');
INSERT INTO public.milestone VALUES (6, 3, 'Milestone F');
INSERT INTO public.milestone VALUES (7, 4, 'Milestone G');
INSERT INTO public.milestone VALUES (8, 4, 'Milestone H');
INSERT INTO public.milestone VALUES (9, 4, 'Milestone I');
INSERT INTO public.milestone VALUES (10, 4, 'Milestone J');


--
-- TOC entry 3466 (class 0 OID 16416)
-- Dependencies: 226
-- Data for Name: person; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.person VALUES (1, 'John', 'Doe');
INSERT INTO public.person VALUES (2, 'Kendall', 'Public');
INSERT INTO public.person VALUES (3, 'Ben', 'Robinson');
INSERT INTO public.person VALUES (4, 'Mike', 'Ho');
INSERT INTO public.person VALUES (5, 'Alex', 'Smith');
INSERT INTO public.person VALUES (6, 'Wendy', 'Smith');
INSERT INTO public.person VALUES (7, 'Karen', 'Wolfe');
INSERT INTO public.person VALUES (8, 'Samantha', 'Jones');
INSERT INTO public.person VALUES (9, 'Linda', 'Brady');
INSERT INTO public.person VALUES (10, 'Jennifer', 'Smith');
INSERT INTO public.person VALUES (11, 'Brett', 'Carlisle');
INSERT INTO public.person VALUES (12, 'Jacob', 'Pratt');


--
-- TOC entry 3467 (class 0 OID 16420)
-- Dependencies: 227
-- Data for Name: person_persontype_assn; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.person_persontype_assn VALUES (1, 2);
INSERT INTO public.person_persontype_assn VALUES (1, 3);
INSERT INTO public.person_persontype_assn VALUES (2, 4);
INSERT INTO public.person_persontype_assn VALUES (2, 5);
INSERT INTO public.person_persontype_assn VALUES (3, 1);
INSERT INTO public.person_persontype_assn VALUES (3, 2);
INSERT INTO public.person_persontype_assn VALUES (3, 3);
INSERT INTO public.person_persontype_assn VALUES (5, 5);
INSERT INTO public.person_persontype_assn VALUES (7, 2);
INSERT INTO public.person_persontype_assn VALUES (7, 4);
INSERT INTO public.person_persontype_assn VALUES (9, 3);
INSERT INTO public.person_persontype_assn VALUES (10, 1);


--
-- TOC entry 3469 (class 0 OID 16424)
-- Dependencies: 229
-- Data for Name: person_type_enum; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.person_type_enum VALUES (4, 'Company Car');
INSERT INTO public.person_type_enum VALUES (1, 'Contractor');
INSERT INTO public.person_type_enum VALUES (3, 'Inactive');
INSERT INTO public.person_type_enum VALUES (2, 'Manager');
INSERT INTO public.person_type_enum VALUES (5, 'Works From Home');


--
-- TOC entry 3471 (class 0 OID 16429)
-- Dependencies: 231
-- Data for Name: person_with_lock; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.person_with_lock VALUES (1, 'John', 'Doe', NULL);
INSERT INTO public.person_with_lock VALUES (2, 'Kendall', 'Public', NULL);
INSERT INTO public.person_with_lock VALUES (3, 'Ben', 'Robinson', NULL);
INSERT INTO public.person_with_lock VALUES (4, 'Mike', 'Ho', NULL);
INSERT INTO public.person_with_lock VALUES (5, 'Alfred', 'Newman', NULL);
INSERT INTO public.person_with_lock VALUES (6, 'Wendy', 'Johnson', NULL);
INSERT INTO public.person_with_lock VALUES (7, 'Karen', 'Wolfe', NULL);
INSERT INTO public.person_with_lock VALUES (8, 'Samantha', 'Jones', NULL);
INSERT INTO public.person_with_lock VALUES (9, 'Linda', 'Brady', NULL);
INSERT INTO public.person_with_lock VALUES (10, 'Jennifer', 'Smith', NULL);
INSERT INTO public.person_with_lock VALUES (11, 'Brett', 'Carlisle', NULL);
INSERT INTO public.person_with_lock VALUES (12, 'Jacob', 'Pratt', NULL);


--
-- TOC entry 3473 (class 0 OID 16434)
-- Dependencies: 233
-- Data for Name: project; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.project VALUES (1, 1, 3, 7, 'ACME Website Redesign', 'The redesign of the main website for ACME Incorporated', '2004-03-01', '2004-07-01', 9560.25, 10250.75);
INSERT INTO public.project VALUES (2, 2, 1, 4, 'State College HR System', 'Implementation of a back-office Human Resources system for State College', '2006-02-15', NULL, 80500.00, 73200.00);
INSERT INTO public.project VALUES (3, 3, 1, 1, 'Blueman Industrial Site Architecture', 'Main website architecture for the Blueman Industrial Group', '2006-03-01', '2006-04-15', 2500.00, 4200.50);
INSERT INTO public.project VALUES (4, 4, 2, 7, 'ACME Payment System', 'Accounts Payable payment system for ACME Incorporated', '2005-08-15', '2005-10-20', 5124.67, 5175.30);


--
-- TOC entry 3475 (class 0 OID 16443)
-- Dependencies: 235
-- Data for Name: project_status_enum; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.project_status_enum VALUES (1, 'Open', 'The project is currently active', 'All projects that we are working on should be in this state', true);
INSERT INTO public.project_status_enum VALUES (2, 'Cancelled', 'The project has been canned', NULL, true);
INSERT INTO public.project_status_enum VALUES (3, 'Completed', 'The project has been completed successfully', 'Celebrate successes!', true);
INSERT INTO public.project_status_enum VALUES (4, 'Planned', 'Project is in the planning stages and has not been assigned a manager', 'Get ready', false);


--
-- TOC entry 3476 (class 0 OID 16449)
-- Dependencies: 236
-- Data for Name: related_project_assn; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.related_project_assn VALUES (1, 3);
INSERT INTO public.related_project_assn VALUES (1, 4);
INSERT INTO public.related_project_assn VALUES (4, 1);


--
-- TOC entry 3477 (class 0 OID 16452)
-- Dependencies: 237
-- Data for Name: team_member_project_assn; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.team_member_project_assn VALUES (1, 3);
INSERT INTO public.team_member_project_assn VALUES (1, 4);
INSERT INTO public.team_member_project_assn VALUES (2, 1);
INSERT INTO public.team_member_project_assn VALUES (2, 2);
INSERT INTO public.team_member_project_assn VALUES (2, 4);
INSERT INTO public.team_member_project_assn VALUES (3, 4);
INSERT INTO public.team_member_project_assn VALUES (4, 2);
INSERT INTO public.team_member_project_assn VALUES (4, 3);
INSERT INTO public.team_member_project_assn VALUES (5, 1);
INSERT INTO public.team_member_project_assn VALUES (5, 2);
INSERT INTO public.team_member_project_assn VALUES (5, 4);
INSERT INTO public.team_member_project_assn VALUES (6, 1);
INSERT INTO public.team_member_project_assn VALUES (6, 3);
INSERT INTO public.team_member_project_assn VALUES (7, 1);
INSERT INTO public.team_member_project_assn VALUES (7, 2);
INSERT INTO public.team_member_project_assn VALUES (8, 1);
INSERT INTO public.team_member_project_assn VALUES (8, 3);
INSERT INTO public.team_member_project_assn VALUES (8, 4);
INSERT INTO public.team_member_project_assn VALUES (9, 2);
INSERT INTO public.team_member_project_assn VALUES (10, 2);
INSERT INTO public.team_member_project_assn VALUES (10, 3);
INSERT INTO public.team_member_project_assn VALUES (11, 4);
INSERT INTO public.team_member_project_assn VALUES (12, 4);


--
-- TOC entry 3494 (class 0 OID 0)
-- Dependencies: 216
-- Name: address_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.address_id_seq', 161, true);


--
-- TOC entry 3495 (class 0 OID 0)
-- Dependencies: 218
-- Name: employee_info_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.employee_info_id_seq', 23, true);


--
-- TOC entry 3496 (class 0 OID 0)
-- Dependencies: 221
-- Name: login_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.login_id_seq', 6, true);


--
-- TOC entry 3497 (class 0 OID 0)
-- Dependencies: 223
-- Name: milestone_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.milestone_id_seq', 10, true);


--
-- TOC entry 3498 (class 0 OID 0)
-- Dependencies: 225
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.person_id_seq', 264, true);


--
-- TOC entry 3499 (class 0 OID 0)
-- Dependencies: 228
-- Name: person_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.person_type_id_seq', 5, true);


--
-- TOC entry 3500 (class 0 OID 0)
-- Dependencies: 230
-- Name: person_with_lock_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.person_with_lock_id_seq', 12, true);


--
-- TOC entry 3501 (class 0 OID 0)
-- Dependencies: 232
-- Name: project_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.project_id_seq', 31, true);


--
-- TOC entry 3502 (class 0 OID 0)
-- Dependencies: 234
-- Name: project_status_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.project_status_type_id_seq', 4, true);


--
-- TOC entry 3263 (class 2606 OID 24782)
-- Name: address idx_16390_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT idx_16390_primary PRIMARY KEY (id);


--
-- TOC entry 3266 (class 2606 OID 24801)
-- Name: employee_info idx_16396_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_info
    ADD CONSTRAINT idx_16396_primary PRIMARY KEY (id);


--
-- TOC entry 3268 (class 2606 OID 24824)
-- Name: gift idx_16400_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gift
    ADD CONSTRAINT idx_16400_primary PRIMARY KEY (number);


--
-- TOC entry 3272 (class 2606 OID 24831)
-- Name: login idx_16404_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login
    ADD CONSTRAINT idx_16404_primary PRIMARY KEY (id);


--
-- TOC entry 3275 (class 2606 OID 24852)
-- Name: milestone idx_16411_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.milestone
    ADD CONSTRAINT idx_16411_primary PRIMARY KEY (id);


--
-- TOC entry 3278 (class 2606 OID 24871)
-- Name: person idx_16416_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person
    ADD CONSTRAINT idx_16416_primary PRIMARY KEY (id);


--
-- TOC entry 3281 (class 2606 OID 24920)
-- Name: person_persontype_assn idx_16420_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_persontype_assn
    ADD CONSTRAINT idx_16420_primary PRIMARY KEY (person_id, person_type_id);


--
-- TOC entry 3284 (class 2606 OID 24934)
-- Name: person_type_enum idx_16424_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_type_enum
    ADD CONSTRAINT idx_16424_primary PRIMARY KEY (id);


--
-- TOC entry 3286 (class 2606 OID 24947)
-- Name: person_with_lock idx_16429_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_with_lock
    ADD CONSTRAINT idx_16429_primary PRIMARY KEY (id);


--
-- TOC entry 3290 (class 2606 OID 24954)
-- Name: project idx_16434_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT idx_16434_primary PRIMARY KEY (id);


--
-- TOC entry 3294 (class 2606 OID 25017)
-- Name: project_status_enum idx_16443_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project_status_enum
    ADD CONSTRAINT idx_16443_primary PRIMARY KEY (id);


--
-- TOC entry 3297 (class 2606 OID 25043)
-- Name: related_project_assn idx_16449_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.related_project_assn
    ADD CONSTRAINT idx_16449_primary PRIMARY KEY (parent_id, child_id);


--
-- TOC entry 3300 (class 2606 OID 25068)
-- Name: team_member_project_assn idx_16452_primary; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.team_member_project_assn
    ADD CONSTRAINT idx_16452_primary PRIMARY KEY (team_member_id, project_id);


--
-- TOC entry 3261 (class 1259 OID 24788)
-- Name: idx_16390_idx_address_1; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16390_idx_address_1 ON public.address USING btree (person_id);


--
-- TOC entry 3264 (class 1259 OID 24807)
-- Name: idx_16396_person_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16396_person_id ON public.employee_info USING btree (person_id);


--
-- TOC entry 3269 (class 1259 OID 24838)
-- Name: idx_16404_idx_login_1; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16404_idx_login_1 ON public.login USING btree (person_id);


--
-- TOC entry 3270 (class 1259 OID 16470)
-- Name: idx_16404_idx_login_2; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16404_idx_login_2 ON public.login USING btree (username);


--
-- TOC entry 3273 (class 1259 OID 24858)
-- Name: idx_16411_idx_milestoneproj_1; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16411_idx_milestoneproj_1 ON public.milestone USING btree (project_id);


--
-- TOC entry 3276 (class 1259 OID 16456)
-- Name: idx_16416_idx_person_1; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16416_idx_person_1 ON public.person USING btree (last_name);


--
-- TOC entry 3279 (class 1259 OID 24921)
-- Name: idx_16420_person_type_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16420_person_type_id ON public.person_persontype_assn USING btree (person_type_id);


--
-- TOC entry 3282 (class 1259 OID 16468)
-- Name: idx_16424_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16424_name ON public.person_type_enum USING btree (name);


--
-- TOC entry 3287 (class 1259 OID 24999)
-- Name: idx_16434_idx_project_2; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16434_idx_project_2 ON public.project USING btree (manager_id);


--
-- TOC entry 3288 (class 1259 OID 24759)
-- Name: idx_16434_num; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16434_num ON public.project USING btree (num);


--
-- TOC entry 3292 (class 1259 OID 16478)
-- Name: idx_16443_idx_projectstatustype_1; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_16443_idx_projectstatustype_1 ON public.project_status_enum USING btree (name);


--
-- TOC entry 3295 (class 1259 OID 25044)
-- Name: idx_16449_idx_relatedprojectassn_2; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16449_idx_relatedprojectassn_2 ON public.related_project_assn USING btree (child_id);


--
-- TOC entry 3298 (class 1259 OID 25069)
-- Name: idx_16452_idx_teammemberprojectassn_2; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16452_idx_teammemberprojectassn_2 ON public.team_member_project_assn USING btree (project_id);


--
-- TOC entry 3291 (class 1259 OID 49638)
-- Name: project_status_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX project_status_id_idx ON public.project USING btree (status_id);


--
-- TOC entry 3313 (class 2620 OID 16561)
-- Name: person_with_lock on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON public.person_with_lock FOR EACH ROW EXECUTE FUNCTION public.on_update_current_timestamp_person_with_lock();


--
-- TOC entry 3302 (class 2606 OID 24892)
-- Name: employee_info employee_info_ibfk_1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.employee_info
    ADD CONSTRAINT employee_info_ibfk_1 FOREIGN KEY (person_id) REFERENCES public.person(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3301 (class 2606 OID 24887)
-- Name: address person_address; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT person_address FOREIGN KEY (person_id) REFERENCES public.person(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3303 (class 2606 OID 24897)
-- Name: login person_login; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login
    ADD CONSTRAINT person_login FOREIGN KEY (person_id) REFERENCES public.person(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3305 (class 2606 OID 24935)
-- Name: person_persontype_assn person_persontype_assn_1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_persontype_assn
    ADD CONSTRAINT person_persontype_assn_1 FOREIGN KEY (person_type_id) REFERENCES public.person_type_enum(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- TOC entry 3306 (class 2606 OID 24909)
-- Name: person_persontype_assn person_persontype_assn_2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person_persontype_assn
    ADD CONSTRAINT person_persontype_assn_2 FOREIGN KEY (person_id) REFERENCES public.person(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- TOC entry 3307 (class 2606 OID 25000)
-- Name: project person_project; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT person_project FOREIGN KEY (manager_id) REFERENCES public.person(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- TOC entry 3311 (class 2606 OID 25057)
-- Name: team_member_project_assn person_team_member_project_assn; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.team_member_project_assn
    ADD CONSTRAINT person_team_member_project_assn FOREIGN KEY (team_member_id) REFERENCES public.person(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3304 (class 2606 OID 24970)
-- Name: milestone project_milestone; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.milestone
    ADD CONSTRAINT project_milestone FOREIGN KEY (project_id) REFERENCES public.project(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- TOC entry 3308 (class 2606 OID 49639)
-- Name: project project_status_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_status_fk FOREIGN KEY (status_id) REFERENCES public.project_status_enum(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3312 (class 2606 OID 25070)
-- Name: team_member_project_assn project_team_member_project_assn; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.team_member_project_assn
    ADD CONSTRAINT project_team_member_project_assn FOREIGN KEY (project_id) REFERENCES public.project(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3309 (class 2606 OID 25032)
-- Name: related_project_assn related_project_assn_1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.related_project_assn
    ADD CONSTRAINT related_project_assn_1 FOREIGN KEY (parent_id) REFERENCES public.project(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- TOC entry 3310 (class 2606 OID 25045)
-- Name: related_project_assn related_project_assn_2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.related_project_assn
    ADD CONSTRAINT related_project_assn_2 FOREIGN KEY (child_id) REFERENCES public.project(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


-- Completed on 2023-04-21 11:27:45 PDT

--
-- PostgreSQL database dump complete
--

