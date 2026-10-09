BEGIN;

INSERT INTO public.quiz_types (name)
SELECT seed.name
FROM (VALUES ('MultipleChoice'), ('SingleChoice')) AS seed(name)
WHERE NOT EXISTS (
    SELECT 1
    FROM public.quiz_types AS quiz_type
    WHERE quiz_type.name = seed.name
);

INSERT INTO public.categories (name)
SELECT 'Geography'
WHERE NOT EXISTS (
    SELECT 1
    FROM public.categories
    WHERE name = 'Geography'
);

INSERT INTO public.questions (question, category_id, quiz_type_id)
SELECT seed.question, category.id, quiz_type.id
FROM (
    VALUES
        ('What is the capital of Germany?', 'SingleChoice'),
        ('What is the capital of Portugal?', 'MultipleChoice'),
        ('What is the capital of Spain?', 'SingleChoice')
) AS seed(question, quiz_type_name)
JOIN public.categories AS category ON category.name = 'Geography'
JOIN LATERAL (
    SELECT id
    FROM public.quiz_types
    WHERE name = seed.quiz_type_name
       OR (seed.quiz_type_name = 'SingleChoice' AND name = 'SingleChoise')
    ORDER BY CASE WHEN name = seed.quiz_type_name THEN 0 ELSE 1 END, id
    LIMIT 1
) AS quiz_type ON true
WHERE NOT EXISTS (
    SELECT 1
    FROM public.questions AS existing
    WHERE existing.category_id = category.id
      AND existing.question = seed.question
);

INSERT INTO public.answers (answer, question_id, is_correct)
SELECT seed.answer, question.id, seed.is_correct
FROM (
    VALUES
        ('What is the capital of Germany?', 'Berlin', true),
        ('What is the capital of Germany?', 'Hamburg', false),
        ('What is the capital of Germany?', 'Munich', false),
        ('What is the capital of Portugal?', 'Lisbon', true),
        ('What is the capital of Portugal?', 'Porto', false),
        ('What is the capital of Spain?', 'Madrid', true),
        ('What is the capital of Spain?', 'Barcelona', false)
) AS seed(question, answer, is_correct)
JOIN public.questions AS question ON question.question = seed.question
JOIN public.categories AS category
    ON category.id = question.category_id
   AND category.name = 'Geography'
WHERE NOT EXISTS (
    SELECT 1
    FROM public.answers AS existing
    WHERE existing.question_id = question.id
      AND existing.answer = seed.answer
);

COMMIT;
