INSERT INTO public.quiz_types ("name",deleted_at,created_at,updated_at) VALUES
	 ('MultipleChoice',NULL,now(),now()),
	 ('SingleChoice',NULL,now(),now());

INSERT INTO public.categories ("name","quiz_type_id",deleted_at,created_at,updated_at) VALUES
	 ('Geography',2,NULL,now(),now());

INSERT INTO public.questions (question,category_id,created_at,updated_at,deleted_at) VALUES
	 ('What is the capital of Germany?',2,now(),now(),NULL),
	 ('What is the capital of Portugal?',2,now(),now(),NULL),
	 ('What is the capital of Spain?',2,now(),now(),NULL);

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


	 INSERT INTO public.answers (answer,question_id,updated_at,deleted_at,created_at,is_correct) VALUES
	('Hamburg',1,now(),NULL,now(),false),
	 ('Berlin',1,now(),NULL,now(),true),
	 ('München',1,now(),NULL,now(),false);

	 INSERT INTO public.answers (answer,question_id,updated_at,deleted_at,created_at,is_correct) VALUES
	 ('Braga',now(),2,NULL,now(),false),
	 ('Porto',now(),2,NULL,now(),false),
	 ('Lisbon',now(),2,NULL,now(),true),

	 INSERT INTO public.answers (answer,question_id,updated_at,deleted_at,created_at,is_correct) VALUES
	 ('Valencia',3,now(),NULL,now(),false),
	 ('Madrid',3,now(),NULL,now(),true),
	 ('Barcelona',3,now(),NULL,now(),false),