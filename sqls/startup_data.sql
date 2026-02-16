INSERT INTO public.quiz_types ("name",deleted_at,created_at,updated_at) VALUES
	 ('MultipleChoice',NULL,now(),now()),
	 ('SingleChoice',NULL,now(),now());

INSERT INTO public.categories ("name",deleted_at,created_at,updated_at) VALUES
	 ('Geography',NULL,now(),now());

INSERT INTO public.questions (question,created_at,updated_at,deleted_at) VALUES
	 ('What is the capital of Germany?',now(),now(),NULL),
	 ('What is the capital of Portugal?',now(),now(),NULL),
	 ('What is the capital of Spain?',now(),now(),NULL);

INSERT INTO public.answers (answer,updated_at,deleted_at,created_at,is_correct) VALUES
	 ('Braga',now(),NULL,now(),false),
	 ('Porto',now(),NULL,now(),false),
	 ('Lisbon',now(),NULL,now(),true),
	 ('Valencia',now(),NULL,now(),false),
	 ('Madrid',now(),NULL,now(),true),
	 ('Barcelona',now(),NULL,now(),false),
	 ('Hamburg',now(),NULL,now(),false),
	 ('Berlin',now(),NULL,now(),true),
	 ('München',now(),NULL,now(),NULL);